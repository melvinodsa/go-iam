package client

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromModelListToSdk(clients []models.Client) []sdk.Client {
	sdkClients := []sdk.Client{}
	for _, client := range clients {
		sdkClients = append(sdkClients, *fromModelToSdk(&client))
	}
	return sdkClients
}

func fromModelToSdk(client *models.Client) *sdk.Client {
	return &sdk.Client{
		Id:                    client.Id,
		Name:                  client.Name,
		Description:           client.Description,
		Secret:                client.Secret,
		Tags:                  client.Tags,
		RedirectURLs:          client.RedirectURLs,
		DefaultAuthProviderId: client.DefaultAuthProviderId,
		GoIamClient:           client.GoIamClient,
		ProjectId:             client.ProjectId,
		Scopes:                client.Scopes,
		LinkedUserId:          client.LinkedUserId,
		ServiceAccountEmail:   client.ServiceAccountEmail,
		Enabled:               client.Enabled,
		CreatedAt:             client.CreatedAt,
		CreatedBy:             client.CreatedBy,
		UpdatedAt:             client.UpdatedAt,
		UpdatedBy:             client.UpdatedBy,
	}
}

func fromSdkToModel(client sdk.Client) models.Client {
	return models.Client{
		Id:                    client.Id,
		Name:                  client.Name,
		Description:           client.Description,
		Secret:                client.Secret,
		Tags:                  client.Tags,
		RedirectURLs:          client.RedirectURLs,
		ProjectId:             client.ProjectId,
		DefaultAuthProviderId: client.DefaultAuthProviderId,
		GoIamClient:           client.GoIamClient,
		ServiceAccountEmail:   client.ServiceAccountEmail,
		Scopes:                client.Scopes,
		LinkedUserId:          client.LinkedUserId,
		Enabled:               client.Enabled,
		CreatedAt:             client.CreatedAt,
		CreatedBy:             client.CreatedBy,
		UpdatedAt:             client.UpdatedAt,
		UpdatedBy:             client.UpdatedBy,
	}
}

func generateRandomSecret(length int) (string, error) {
	// Calculate the number of bytes needed
	byteLength := (length*6 + 7) / 8 // Convert bit length to byte length
	bytes := make([]byte, byteLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode bytes to a URL-safe base64 string and truncate to the desired length
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
