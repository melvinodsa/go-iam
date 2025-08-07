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
	authP      authprovider.Service
	clientSvc  client.Service
	cacheSvc   cache.Service
	jwtSvc     jwt.Service
	encSvc     encrypt.Service
	usrSvc     user.Service
	tokenTTL   int64
	refetchTTL int64
}

func NewService(authP authprovider.Service, clientSvc client.Service, cacheSvc cache.Service, jwtSvc jwt.Service, encSvc encrypt.Service, usrSvc user.Service, tokenTTL int64, refetchTTL int64) *service {
	return &service{
		authP:      authP,
		clientSvc:  clientSvc,
		cacheSvc:   cacheSvc,
		jwtSvc:     jwtSvc,
		encSvc:     encSvc,
		usrSvc:     usrSvc,
		tokenTTL:   tokenTTL,
		refetchTTL: refetchTTL,
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
	sp, err := s.authP.GetProvider(ctx, *p)
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

	// get the auth token from cache

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
	usr, err := s.getUserFromCache(ctx, accessToken)
	if err == nil && usr != nil {
		log.Debugf("fetched user records from cache - %s", usr.Id)
		return usr, nil
	}
	token, err := s.getAccessTokenFromCache(ctx, accessTokenId)
	if err != nil {
		return nil, fmt.Errorf("error getting the token from cache %w", err)
	}
	identity, err := s.getAuthProivderIdentity(ctx, token, accessTokenId)
	if err != nil {
		return nil, fmt.Errorf("error getting the identity from auth provider %w", err)
	}
	usr, err = s.getOrCreateUser(ctx, *identity)
	if err != nil {
		return nil, fmt.Errorf("error getting or creating the user %w", err)
	}
	log.Debugf("fetched user records from auth provider - %s", usr.Id)
	err = s.cacheUserDetails(ctx, accessToken, *usr)
	if err != nil {
		return nil, fmt.Errorf("error caching the user details %w", err)
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

func (s service) getAuthProivderIdentity(ctx context.Context, token *sdk.AuthToken, accessTokenId string) (*sdk.User, error) {
	/*
	 * get the service provider
	 * call the get identity method on the service provider
	 */
	p, err := s.authP.Get(ctx, token.AuthProviderID, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching auth provider details %w", err)
	}
	sp, err := s.authP.GetProvider(ctx, *p)
	if err != nil {
		return nil, fmt.Errorf("error getting service provider %w", err)
	}

	// if the token is expired, we need to refresh the token
	if token.ExpiresAt.Before(time.Now()) {
		newToken, err := s.refreshAuthToken(ctx, accessTokenId, *token, sp)
		if err != nil {
			return nil, fmt.Errorf("error refreshing the token %w", err)
		}
		*token = *newToken
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

func (s service) cacheUserDetails(ctx context.Context, accessToken string, user sdk.User) error {
	/*
	 * encode the user to json
	 * generate a new user id
	 * save the user in cache
	 */
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(user)
	if err != nil {
		return fmt.Errorf("error encoding the user %w", err)
	}
	userEnc, err := s.encSvc.Encrypt(b.String())
	if err != nil {
		return fmt.Errorf("error encrypting the user %w", err)
	}
	err = s.cacheSvc.Set(ctx, fmt.Sprintf("token-%s", accessToken), userEnc, time.Minute*time.Duration(s.refetchTTL))
	if err != nil {
		return fmt.Errorf("error saving the user %w", err)
	}
	return nil
}

func (s service) getUserFromCache(ctx context.Context, accessToken string) (*sdk.User, error) {
	/*
	 * get the value from cache
	 */
	val, err := s.cacheSvc.Get(ctx, fmt.Sprintf("token-%s", accessToken))
	if err != nil {
		return nil, fmt.Errorf("error fetching the value from cache %w", err)
	}

	userDec, err := s.encSvc.Decrypt(val)
	if err != nil {
		return nil, fmt.Errorf("error decrypting the user %w", err)
	}
	result := sdk.User{}
	err = json.NewDecoder(strings.NewReader(userDec)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding the user %w", err)
	}
	return &result, nil
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
	sp, err := s.authP.GetProvider(ctx, *p)
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

	err = s.cacheSvc.Set(ctx, fmt.Sprintf("access-token-%s", accessToken), auEnc, time.Minute*time.Duration(s.tokenTTL))
	if err != nil {
		return "", fmt.Errorf("error saving the access token %w", err)
	}
	return accessToken, nil
}

func (s service) getAccessTokenFromCache(ctx context.Context, accessToken string) (*sdk.AuthToken, error) {
	/*
	 * get the value from cache
	 */

	val, err := s.cacheSvc.Get(ctx, fmt.Sprintf("access-token-%s", accessToken))
	if err != nil {
		return nil, fmt.Errorf("error fetching the value from cache %w", err)
	}

	auDec, err := s.encSvc.Decrypt(val)
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
	err = s.cacheSvc.Set(ctx, fmt.Sprintf("auth-code-%s", authCode), auEnc, time.Minute)
	if err != nil {
		return "", fmt.Errorf("error saving the auth code %w", err)
	}
	return authCode, nil
}

func (s service) getAuthTokenFromCache(ctx context.Context, authCode string) (*sdk.AuthToken, error) {
	/*
	 * get the value from cache
	 */

	val, err := s.cacheSvc.Get(ctx, fmt.Sprintf("auth-code-%s", authCode))
	if err != nil {
		return nil, fmt.Errorf("error fetching the value from cache %w", err)
	}

	auDec, err := s.encSvc.Decrypt(val)
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
	err := s.cacheSvc.Delete(ctx, fmt.Sprintf("auth-code-%s", authCode))
	if err != nil {
		return fmt.Errorf("error deleting the value from cache %w", err)
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
	err = s.cacheSvc.Set(ctx, fmt.Sprintf("state-%s", stateId), st, time.Minute*5)
	if err != nil {
		return "", fmt.Errorf("error saving the state %w", err)
	}
	return stateId, nil
}

func (s service) getCacheState(ctx context.Context, stateId string) (string, string, string, string, error) {
	/*
	 * get the state from cache
	 */
	val, err := s.cacheSvc.Get(ctx, fmt.Sprintf("state-%s", stateId))
	if err != nil {
		return "", "", "", "", fmt.Errorf("error fetching the state from cache %w", err)
	}

	state, err := s.encSvc.Decrypt(val)
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
	err := s.cacheSvc.Delete(ctx, fmt.Sprintf("state-%s", stateId))
	if err != nil {
		return fmt.Errorf("error deleting the value from cache %w", err)
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

func (s service) ClientCredentials(ctx context.Context, clientId, clientSecret string) (*sdk.ClientCredentialsDataResponse, error) {

	// Get client details
	cl, err := s.clientSvc.Get(ctx, clientId, true)
	if err != nil {
		return nil, fmt.Errorf("invalid client_id: %w", err)
	}

	// Verify client is enabled
	if !cl.Enabled {
		return nil, errors.New("client is disabled")
	}

	// Verify client secret
	if !s.clientSvc.VerifySecret(clientSecret, cl.Secret) {
		return nil, errors.New("invalid client_secret")
	}
	// Check if client has a linked user
	if cl.LinkedUserId == "" {
		return nil, errors.New("client does not support client credentials flow")
	}

	// Get the linked user
	user, err := s.usrSvc.GetById(ctx, cl.LinkedUserId)
	if err != nil {
		return nil, fmt.Errorf("linked user not found: %w", err)
	}

	// Verify user is enabled
	if !user.Enabled {
		return nil, errors.New("linked user is disabled")
	}

	// Check user expiry
	if user.Expiry != nil && user.Expiry.Before(time.Now()) {
		return nil, errors.New("linked user has expired")
	}

	// Generate access token ID and cache user details
	accessTokenId := uuid.New().String()

	expiryTime := time.Now().AddDate(0, 0, 1)
	claims := map[string]interface{}{
		"id":         accessTokenId,
		"grant_type": "client_credentials",
		"client_id":  clientId,
	}

	accessToken, err := s.jwtSvc.GenerateToken(claims, expiryTime.Unix())
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}
	err = s.cacheUserDetails(ctx, accessToken, *user)
	if err != nil {
		return nil, fmt.Errorf("error caching user details: %w", err)
	}

	log.Debugf("client credentials authentication successful for client %s as user %s", clientId, user.Id)

	return &sdk.ClientCredentialsDataResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours in seconds
	}, nil
}
