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

	// Get the token if it came in from the environment
	clientToken := client.Token()

	// If we don't have a token, check the token helper
	if clientToken == "" {
		helper, err := token.NewInternalTokenHelper()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get token helper")
		}
		clientToken, err = helper.Get()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get token from token helper")
		}
	}

	// Set the token
	if clientToken != "" {
		client.SetToken(clientToken)
	} else {
		return nil, errors.Wrap(err, "failed to get token from environment or credential helper")
	}

	return client, nil
}
