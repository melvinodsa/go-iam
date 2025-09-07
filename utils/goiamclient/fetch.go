package goiamclient

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/client"
)

func GetGoIamClient(svc client.Service) (*sdk.Client, error) {
	prvs, err := svc.GetGoIamClients(context.Background(), sdk.ClientQueryParams{
		GoIamClient: true,
	})
	if err != nil {
		return nil, err
	}
	if len(prvs) == 0 {
		log.Warn("IAM running in insecure mode. Create a client for Go IAM to make the application secure")
		return nil, nil
	}
	return &prvs[0], nil
}
