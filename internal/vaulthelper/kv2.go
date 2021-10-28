package vaulthelper

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

const (
	defaultKv2SecretEnginePath      = "secret"
	defaultKv2UsernameAttributeName = "username"
	defaultKv2PasswordAttributeName = "password"
)

type Kv2SecretsEngineResponse struct {
	Username string
	Password string
}

func ReadKv2SecretsEngine(client *api.Client, mountPath string, secretPath string, attributeMap map[string]string) (*Kv2SecretsEngineResponse, error) {
	resp := &Kv2SecretsEngineResponse{}

	if mountPath == "" {
		mountPath = defaultKv2SecretEnginePath
	}

	if secretPath == "" {
		return resp, errors.New("empty rolename provided to secrets engine reader")
	}

	am := parseAttributeMap(attributeMap)
	if len(am) == 0 {
		return resp, fmt.Errorf("unable to parse attribute map")
	}

	vaultSecretPath := fmt.Sprintf("%s/%s", mountPath, secretPath)

	secret, err := client.Logical().Read(vaultSecretPath)
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", vaultSecretPath))
	}

	if secret == nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", vaultSecretPath))
	}

	resp.Username = secret.Data[am["username"]].(string)
	resp.Password = secret.Data[am["password"]].(string)

	return resp, nil
}

func parseAttributeMap(attributeMap map[string]string) map[string]string {
	if _, ok := attributeMap["username"]; !ok {
		attributeMap["username"] = defaultKv2UsernameAttributeName
	}
	if _, ok := attributeMap["password"]; !ok {
		attributeMap["password"] = defaultKv2PasswordAttributeName
	}
	return attributeMap
}
