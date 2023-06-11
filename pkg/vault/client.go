package vault

import (
	"context"
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Cli *vault.Client
}

func NewVaultClient() (*VaultClient, error) {
	config := vault.DefaultConfig()

	config.Address = "https://vault.dashwave.io"

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return &VaultClient{Cli: client}, nil
}

func (vc *VaultClient) GetSecret(secretPath, secretKey string) (string, error) {
	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := vc.Cli.KVv2("kv-v2").Get(context.Background(), secretPath)
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %v", err)
	}

	value, ok := secret.Data[secretKey].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data[secretKey], secret.Data[secretKey])
	}

	return value, nil
}

func (vc *VaultClient) PutSecret(secretpath, secretKey, secretValue string) error {
	secretData := map[string]interface{}{
		secretKey: secretValue,
	}

	// Write a secret
	_, err := vc.Cli.KVv2("kv-v2").Put(context.Background(), secretpath, secretData)
	if err != nil {
		return fmt.Errorf("unable to write secret: %v", err)
	}

	return nil
}
