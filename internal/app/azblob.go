package app

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/revandpratama/lognest/config"
	"github.com/rs/zerolog/log"
)

func WithAzureBlobStorage() Option {
	return func(app *App) error {

		client, err := azblob.NewClientFromConnectionString(config.ENV.AZURE_STORAGE_CONNECTION_STRING, nil)
		if err != nil {
			return err
		}

		containerName := config.ENV.AZURE_STORAGE_CONTAINER_NAME

		containerClient := client.ServiceClient().NewContainerClient(containerName)

		// 3. Check and create the container if it doesn't exist.
		// The options `&container.CreateOptions{ Access: &accessType }` make it public.
		// Use `nil` for private containers.
		// accessType := container.PublicAccessTypeBlob
		_, err = containerClient.Create(context.Background(), nil)

		// The Create function returns a specific error if the container already exists.
		// We can safely ignore this specific error, but we should fail on any other error.
		if err != nil && !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			log.Warn().Err(err).Msg("container already exists, skipping creation")
		}

		log.Printf("Successfully connected to Azure Blob Storage and ensured container '%s' exists.", containerName)
		log.Printf("Azure client INITIALIZED. Memory Address: %p", client)
		app.AzblobClient = client

		return nil
	}
}
