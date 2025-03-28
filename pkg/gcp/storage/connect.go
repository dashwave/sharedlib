package storage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/dashwave/sharedlib/pkg/vault"
	"google.golang.org/api/option"
)

const (
	US_VAULT    = "US-VAULT"
	INDIA_VAULT = "INDIA-VAULT"
)

func ConnectGCP(v *vault.VaultClient, accountLocation string) (*storage.Client, vault.VaultSecretMap, error) {
	secretPath := ""
	if accountLocation == US_VAULT {
		secretPath = "US-GCP-ACCOUNT"
	} else if accountLocation == INDIA_VAULT {
		secretPath = "INDIA-GCP-ACCOUNT"
	} else {
		return nil, nil, fmt.Errorf("invalid GCP account location provided : %s", accountLocation)
	}

	secrets, err := v.GetSecretMap(secretPath)
	if err != nil {
		return nil, nil, err
	}

	credentialsJSON := secrets["GCP_CREDENTIALS_JSON"].(string)
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(credentialsJSON)))
	if err != nil {
		return nil, nil, err
	}

	return client, secrets, nil
}

func ConnectStorage() (*storage.Client, vault.VaultSecretMap) {
	vaultClient, err := vault.NewVaultClient()
	if err != nil {
		panic(err)
	}

	storageClient, secrets, err := ConnectGCP(vaultClient, INDIA_VAULT)
	if err != nil {
		panic(err)
	}

	return storageClient, secrets
}

func GetGCPClient() (*storage.Client, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}
