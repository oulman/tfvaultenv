package vaulthelper

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/pkg/errors"
)

func GetVaultApiClient() (*api.Client, error) {

	config := api.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, errors.Wrap(err, "failed to read environment")
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}

	return client, nil
}

func SetVaultTokenFromEnvOrHandler(client *api.Client) error {
	// Get the token if it came in from the environment
	clientToken := client.Token()

	// If we don't have a token, check the token helper
	if clientToken == "" {
		helper, err := token.NewInternalTokenHelper()
		if err != nil {
			return errors.Wrap(err, "failed to get token helper")
		}
		clientToken, err = helper.Get()
		if err != nil {
			return errors.Wrap(err, "failed to get token from token helper")
		}
	}

	// Set the token
	if clientToken != "" {
		client.SetToken(clientToken)
		return nil
	} else {
		return fmt.Errorf("failed to get token from environment or credential helper")
	}
}

func PrintVaultToken(token string) {
	fmt.Printf("VAULT_TOKEN=%s\n", token)
}
