/*
Copyright © 2021 James Oulman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vaulthelper

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

const (
	defaultAdProviderPath = "ad"
)

type AdSecretsEngineResponse struct {
	Username        string
	CurrentPassword string
	LastPassword    string
}

func ReadAdSecretsEngine(client *api.Client, rolename string, vaultpath string) (*AdSecretsEngineResponse, error) {
	resp := &AdSecretsEngineResponse{}

	if rolename == "" {
		return resp, errors.New("empty rolename provided to secrets engine reader")
	}

	if vaultpath == "" {
		vaultpath = defaultAdProviderPath
	}

	secretPath := fmt.Sprintf("%s/creds/%s", vaultpath, rolename)

	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s", secretPath))
	}

	if secret == nil {
		return resp, fmt.Errorf("no value found at %s", secretPath)
	}

	resp.Username = secret.Data["username"].(string)
	resp.CurrentPassword = secret.Data["current_password"].(string)

	return resp, nil
}
