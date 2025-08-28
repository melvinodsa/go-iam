package auth

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

func (s *service) HandleEvent(e utils.Event[sdk.Client]) {
	switch e.Name() {
	case goiamuniverse.EventClientCreated:
		s.handleClientUpdates(e)
	case goiamuniverse.EventClientUpdated:
		s.handleClientUpdates(e)
	default:
		return
	}

}

func (s *service) handleClientUpdates(e utils.Event[sdk.Client]) {
	s.cacheClientSecret(context.Background(), e.Payload().Id, e.Payload().Secret)
}
