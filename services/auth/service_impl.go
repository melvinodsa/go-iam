package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/authprovider"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/services/client"
	"github.com/melvinodsa/go-iam/services/encrypt"
	"github.com/melvinodsa/go-iam/services/jwt"
	"github.com/melvinodsa/go-iam/services/user"
)

type service struct {
	authP     authprovider.Service
	clientSvc client.Service
	cacheSvc  cache.Service
	jwtSvc    jwt.Service
	encSvc    encrypt.Service
	usrSvc    user.Service
}

func NewService(authP authprovider.Service, clientSvc client.Service, cacheSvc cache.Service, jwtSvc jwt.Service, encSvc encrypt.Service, usrSvc user.Service) Service {
	return &service{
		authP:     authP,
		clientSvc: clientSvc,
		cacheSvc:  cacheSvc,
		jwtSvc:    jwtSvc,
		encSvc:    encSvc,
		usrSvc:    usrSvc,
	}
}

func (s service) GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl string) (string, error) {
	/*
	 * We first get the client details from the client service if authproviderid is not provided
	 * Then we will get the auth provider details from the auth provider service for the default auth provider
	 * Then we will call the GetLoginUrl method on the auth provider
	 */
	if len(authProviderId) == 0 {
		client, err := s.clientSvc.Get(ctx, clientId, true)
		if err != nil {
			return "", fmt.Errorf("error fetching client details %w", err)
		}
		authProviderId = client.DefaultAuthProviderId
	}

	p, err := s.authP.Get(ctx, authProviderId, true)
	if err != nil {
		return "", fmt.Errorf("error fetching auth provider details %w", err)
	}
	sp, err := s.authP.GetProvider(*p)
	if err != nil {
		return "", fmt.Errorf("error getting service provider %w", err)
	}

	// it is important to note that we are combining the state with the client id
	newState, err := s.cacheState(ctx, state, clientId, p.Id, redirectUrl)
	if err != nil {
		return "", fmt.Errorf("error caching the state %w", err)
	}
	return sp.GetAuthCodeUrl(newState), nil
}
func (s service) Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error) {
	/*
	 * get the state, authprovider id and client id from the state
	 * generate the access token
	 * cache the token
	 * get the callback details from client service
	 * return the callback details
	 */
	clientId, oState, authProviderId, redirectUrl, err := s.getCacheState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("error getting the state from cache %w", err)
	}

	token, err := s.getToken(ctx, authProviderId, code)
	if err != nil {
		return nil, fmt.Errorf("error getting the token %w", err)
	}

	authCode, err := s.cacheAuthToken(ctx, *token)
	if err != nil {
		return nil, fmt.Errorf("error caching the token %w", err)
	}

	redirectUrl, err = s.getRedirectUrl(ctx, clientId, redirectUrl, authCode, oState)
	if err != nil {
		return nil, fmt.Errorf("error getting the callback url %w", err)
	}

	err = s.invalidateState(ctx, state)
	if err != nil {
		log.Errorf("error invalidating state %s", err)
	}

	return &sdk.AuthRedirectResponse{RedirectUrl: redirectUrl}, nil
}

func (s service) ClientCallback(ctx context.Context, code string) (*sdk.AuthVerifyCodeResponse, error) {
	/*
	 * get the code from the cache
	 * generate the access token and store the original token in cache
	 * invalidate the code from cache
	 * return the access token
	 */
	token, err := s.getAuthTokenFromCache(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error getting the token from cache %w", err)
	}

	accessTokenId, err := s.cacheAccessToken(ctx, *token, "")
	if err != nil {
		return nil, fmt.Errorf("error caching the access token %w", err)
	}

	// generate jwt access token
	accessToken, err := s.jwtSvc.GenerateToken(map[string]interface{}{"id": accessTokenId}, time.Now().AddDate(0, 0, 1).Unix())
	if err != nil {
		return nil, fmt.Errorf("error generating the access token %w", err)
	}

	err = s.invalidateAuthToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error invalidating the auth code %w", err)
	}

	return &sdk.AuthVerifyCodeResponse{AccessToken: accessToken}, nil

}

func (s service) GetIdentity(ctx context.Context, accessToken string) (*sdk.User, error) {
	/*
	 * get id from jwt access token
	 * get the access token from cache
	 * get the user details from the access token
	 */
	claims, err := s.jwtSvc.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("error validating the access token %w", err)
	}
	accessTokenId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("error getting the access token id from claims")
	}
	token, err := s.getAccessTokenFromCache(ctx, accessTokenId)
	if err != nil {
		return nil, fmt.Errorf("error getting the token from cache %w", err)
	}
	identity, err := s.getAuthProivderIdentity(ctx, *token, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error getting the identity from auth provider %w", err)
	}
	usr, err := s.getOrCreateUser(ctx, *identity)
	if err != nil {
		return nil, fmt.Errorf("error getting or creating the user %w", err)
	}

	return usr, nil
}

func (s service) getOrCreateUser(ctx context.Context, usr sdk.User) (*sdk.User, error) {
	/*
	 * get the user from the user service
	 * if user not found, create the user
	 */
	var u *sdk.User
	var err error
	if len(usr.Email) > 0 {
		u, err = s.usrSvc.GetByEmail(ctx, usr.Email, usr.ProjectId)
	} else if len(usr.Phone) > 0 {
		u, err = s.usrSvc.GetByPhone(ctx, usr.Phone, usr.ProjectId)
	} else {
		return nil, fmt.Errorf("email or phone is required")
	}
	if err != nil && errors.Is(err, user.ErrorUserNotFound) {
		// we need to create the user
		err = s.usrSvc.Create(ctx, &usr)
		if err != nil {
			return nil, fmt.Errorf("error creating the user %w", err)
		}
		u = &usr
	}
	if !u.Enabled {
		return nil, errors.New("user is disabled")
	}

	if u.Expiry != nil && u.Expiry.Before(time.Now()) {
		return nil, errors.New("user expired")
	}

	return u, nil
}

func (s service) getAuthProivderIdentity(ctx context.Context, token sdk.AuthToken, accessToken string) (*sdk.User, error) {
	/*
	 * get the service provider
	 * call the get identity method on the service provider
	 */
	p, err := s.authP.Get(ctx, token.AuthProviderID, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching auth provider details %w", err)
	}
	sp, err := s.authP.GetProvider(*p)
	if err != nil {
		return nil, fmt.Errorf("error getting service provider %w", err)
	}

	// if the token is expired, we need to refresh the token
	if token.ExpiresAt.Before(time.Now()) {
		newToken, err := s.refreshAuthToken(ctx, accessToken, token, sp)
		if err != nil {
			return nil, fmt.Errorf("error refreshing the token %w", err)
		}
		token = *newToken
	}

	identity, err := sp.GetIdentity(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error getting the identity from service provider %w", err)
	}
	user := sdk.User{ProjectId: p.ProjectId}
	for _, id := range identity {
		id.UpdateUserDetails(&user)
	}
	return &user, nil

}

func (s service) refreshAuthToken(ctx context.Context, accessToken string, token sdk.AuthToken, sp sdk.ServiceProvider) (*sdk.AuthToken, error) {
	tk, err := sp.RefreshToken(token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error refreshing the token with servivce provider%w", err)
	}
	// cache the service provider token
	_, err = s.cacheAccessToken(ctx, token, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error caching the access token %w", err)
	}
	return tk, nil
}

func (s service) getToken(ctx context.Context, authProviderId, code string) (*sdk.AuthToken, error) {
	/*
	 * get the client details
	 * get the auth provider details
	 * get the service provider
	 * call the verify code on the service provider
	 */
	p, err := s.authP.Get(ctx, authProviderId, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching auth provider details %w", err)
	}
	sp, err := s.authP.GetProvider(*p)
	if err != nil {
		return nil, fmt.Errorf("error getting service provider %w", err)
	}

	token, err := sp.VerifyCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error verifying the code %w", err)
	}
	token.AuthProviderID = authProviderId
	return token, nil
}

func (s service) cacheAccessToken(ctx context.Context, token sdk.AuthToken, accessToken string) (string, error) {
	/*
	 * generate a new access token id
	 * save the token in cache
	 */
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(token)
	if err != nil {
		return "", fmt.Errorf("error encoding the token %w", err)
	}
	auEnc, err := s.encSvc.Encrypt(b.String())
	if err != nil {
		return "", fmt.Errorf("error encrypting the access token %w", err)
	}
	if len(accessToken) == 0 {
		accessToken = uuid.New().String()
	}
	res := s.cacheSvc.Redis.Set(ctx, fmt.Sprintf("access-token-%s", accessToken), auEnc, time.Hour*24)
	if res.Err() != nil {
		return "", fmt.Errorf("error saving the access token %w", res.Err())
	}
	return accessToken, nil
}

func (s service) getAccessTokenFromCache(ctx context.Context, accessToken string) (*sdk.AuthToken, error) {
	/*
	 * get the value from cache
	 */
	res := s.cacheSvc.Redis.Get(ctx, fmt.Sprintf("access-token-%s", accessToken))
	if res.Err() != nil {
		return nil, fmt.Errorf("error fetching the value from cache %w", res.Err())
	}
	auDec, err := s.encSvc.Decrypt(res.Val())
	if err != nil {
		return nil, fmt.Errorf("error decrypting the access token %w", err)
	}
	result := sdk.AuthToken{}
	err = json.NewDecoder(strings.NewReader(auDec)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding the token %w", err)
	}
	return &result, nil
}

func (s service) cacheAuthToken(ctx context.Context, token sdk.AuthToken) (string, error) {
	/*
	 * encode the token to json
	 * generate a new auth code id
	 * save the token in cache
	 */
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(token)
	if err != nil {
		return "", fmt.Errorf("error encoding the token %w", err)
	}
	auEnc, err := s.encSvc.Encrypt(b.String())
	if err != nil {
		return "", fmt.Errorf("error encrypting the access token %w", err)
	}
	authCode := uuid.New().String()
	res := s.cacheSvc.Redis.Set(ctx, fmt.Sprintf("auth-code-%s", authCode), auEnc, time.Minute)
	if res.Err() != nil {
		return "", fmt.Errorf("error saving the auth code %w", res.Err())
	}
	return authCode, nil
}

func (s service) getAuthTokenFromCache(ctx context.Context, authCode string) (*sdk.AuthToken, error) {
	/*
	 * get the value from cache
	 */
	res := s.cacheSvc.Redis.Get(ctx, fmt.Sprintf("auth-code-%s", authCode))
	if res.Err() != nil {
		return nil, fmt.Errorf("error fetching the value from cache %w", res.Err())
	}
	auDec, err := s.encSvc.Decrypt(res.Val())
	if err != nil {
		return nil, fmt.Errorf("error decrypting the access token %w", err)
	}
	result := sdk.AuthToken{}
	err = json.NewDecoder(strings.NewReader(auDec)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding the token %w", err)
	}
	return &result, nil
}

func (s service) invalidateAuthToken(ctx context.Context, authCode string) error {
	/*
	 * delete the value from cache
	 */
	res := s.cacheSvc.Redis.Del(ctx, fmt.Sprintf("auth-code-%s", authCode))
	if res.Err() != nil {
		return fmt.Errorf("error deleting the value from cache %w", res.Err())
	}
	return nil
}

func (s service) cacheState(ctx context.Context, state, clientId, providerId, redirectUrl string) (string, error) {
	/*
	 * add the extra info required in cache
	 * generate a new state id
	 * save the state in cache
	 */
	newState := fmt.Sprintf("%s:%s:%s:%s", state, clientId, providerId, url.QueryEscape(redirectUrl))
	st, err := s.encSvc.Encrypt(newState)
	if err != nil {
		return "", fmt.Errorf("error encrypting the state %w", err)
	}
	stateId := uuid.New().String()
	res := s.cacheSvc.Redis.Set(ctx, fmt.Sprintf("state-%s", stateId), st, time.Minute*5)
	if res.Err() != nil {
		return "", fmt.Errorf("error saving the state %w", res.Err())
	}
	return stateId, nil
}

func (s service) getCacheState(ctx context.Context, stateId string) (string, string, string, string, error) {
	/*
	 * get the state from cache
	 */
	res := s.cacheSvc.Redis.Get(ctx, fmt.Sprintf("state-%s", stateId))
	if res.Err() != nil {
		return "", "", "", "", fmt.Errorf("error fetching the state from cache %w", res.Err())
	}
	state, err := s.encSvc.Decrypt(res.Val())
	if err != nil {
		return "", "", "", "", fmt.Errorf("error decrypting the state %w", err)
	}
	stateParts := strings.Split(state, ":")
	if len(stateParts) != 4 {
		return "", "", "", "", fmt.Errorf("invalid state. expected to have 4 parts but got %d", len(stateParts))
	}
	clientId := stateParts[1]
	oState := stateParts[0]
	authProviderId := stateParts[2]
	redirectUrl := stateParts[3]
	urlDecoded, err := url.QueryUnescape(redirectUrl)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error decoding the redirect url %w", err)
	}
	return clientId, oState, authProviderId, urlDecoded, nil
}

func (s service) invalidateState(ctx context.Context, stateId string) error {
	/*
	 * delete the value from cache
	 */
	res := s.cacheSvc.Redis.Del(ctx, fmt.Sprintf("state-%s", stateId))
	if res.Err() != nil {
		return fmt.Errorf("error deleting the value from cache %w", res.Err())
	}
	return nil
}

func (s service) getRedirectUrl(ctx context.Context, clientId, redirectUrl, authCode, state string) (string, error) {
	/*
	 * get the client details
	 * get the auth provider details
	 * get the service provider
	 * get the callback url
	 */
	cl, err := s.clientSvc.Get(ctx, clientId, true)
	if err != nil {
		return "", fmt.Errorf("error fetching client details %w", err)
	}
	found := false
	for _, cb := range cl.RedirectURLs {
		if strings.EqualFold(cb, redirectUrl) {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("callback url not found in the client details - %s", redirectUrl)
	}
	if strings.Contains(redirectUrl, "?") {
		redirectUrl = fmt.Sprintf("%s&code=%s&state=%s", redirectUrl, authCode, state)
	} else {
		redirectUrl = fmt.Sprintf("%s?code=%s&state=%s", redirectUrl, authCode, state)
	}
	return redirectUrl, nil
}
