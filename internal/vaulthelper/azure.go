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
	defaultAzureSecretsEnginePath = "azure"
)

type AzureSecretsEngineResponse struct {
	ClientId     string
	ClientSecret string
}

func (r *AzureSecretsEngineResponse) IsPopulated() bool {
	return r.ClientId != "" && r.ClientSecret != ""
}

func ReadAzureSecretsEngine(client *api.Client, mount string, role string) (*AzureSecretsEngineResponse, error) {
	resp := &AzureSecretsEngineResponse{}

	if mount == "" {
		mount = defaultAzureSecretsEnginePath
	}

	if role == "" {
		return resp, errors.New("empty role provided to Azure secrets engine reader")
	}

	path := fmt.Sprintf("%s/creds/%s", mount, role)

	secret, err := client.Logical().Read(path)
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", path))
	}

	if secret == nil {
		return resp, errors.Wrap(err, fmt.Sprintf("no value from Vault at %s\n", path))
	}

	resp.ClientId = secret.Data["client_id"].(string)
	resp.ClientSecret = secret.Data["client_secret"].(string)

	return resp, nil
}
