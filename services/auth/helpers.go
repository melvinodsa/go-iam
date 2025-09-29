package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
)

func (s *service) cacheClientSecret(ctx context.Context, clientId string, secret string) {
	err := s.cacheSvc.Set(ctx, fmt.Sprintf("client-%s", clientId), secret, time.Hour*24*365)
	if err != nil {
		log.Errorf("failed to cache client secret: %w", err)
	}
}

func (s *service) cacheClientProjectID(ctx context.Context, clientId string, projectId string) {
	err := s.cacheSvc.Set(ctx, fmt.Sprintf("client-%s-project", clientId), projectId, time.Hour*24*365)
	if err != nil {
		log.Errorf("failed to cache client project ID: %w", err)
	}
}

func (s *service) getClientProjectID(ctx context.Context, clientId string) (string, error) {
	projectId, err := s.cacheSvc.Get(ctx, fmt.Sprintf("client-%s-project", clientId))
	if err == nil {
		return projectId, nil
	}
	cl, err := s.clientSvc.Get(ctx, clientId, true)
	if err != nil {
		return "", fmt.Errorf("couldn't get the client even from db: %w", err)
	}
	s.cacheClientSecret(ctx, cl.Id, cl.Secret)
	s.cacheClientProjectID(ctx, cl.Id, cl.ProjectId)
	return cl.ProjectId, nil
}

func (s *service) getClientSecret(ctx context.Context, clientId string) (string, error) {
	secret, err := s.cacheSvc.Get(ctx, fmt.Sprintf("client-%s", clientId))
	if err == nil {
		return secret, nil
	}
	cl, err := s.clientSvc.Get(ctx, clientId, true)
	if err != nil {
		return "", fmt.Errorf("couldn't get the client even from db: %w", err)
	}
	s.cacheClientSecret(ctx, cl.Id, cl.Secret)
	s.cacheClientProjectID(ctx, cl.Id, cl.ProjectId)
	return cl.Secret, nil
}

func (s *service) handlePrivateClient(ctx context.Context, clientId, clientSecret string) error {
	secret, err := s.getClientSecret(ctx, clientId)
	if err != nil {
		return fmt.Errorf("error getting client secret: %w", err)
	}
	err = s.clientSvc.VerifySecret(clientSecret, secret)
	if err != nil {
		return fmt.Errorf("invalid client secret: %w", err)
	}
	return nil
}

func (s *service) handlePublicClient(clientId, codeChallenge string, token sdk.AuthToken) error {
	// Implement public client handling logic here
	if token.CodeChallengeMethod != "S256" {
		return fmt.Errorf("invalid code challenge")
	}
	calculatedVerifier := generateCodeChallengeS256(token.CodeChallenge)
	// Verify the code verifier
	if strings.Compare(calculatedVerifier, codeChallenge) != 0 {
		log.Debugw("invalid code verifier", "calculated_verifier", calculatedVerifier, "code_challenge", codeChallenge)
		return fmt.Errorf("invalid code verifier")
	}
	if strings.Compare(token.ClientId, clientId) != 0 {
		return fmt.Errorf("invalid client id")
	}
	return nil
}

func generateCodeChallengeS256(codeChallenge string) string {
	hash := sha256.Sum256([]byte(codeChallenge))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
