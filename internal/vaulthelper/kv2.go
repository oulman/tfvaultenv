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
	defaultKv2SecretEnginePath      = "secret"
	defaultKv2UsernameAttributeName = "username"
	defaultKv2PasswordAttributeName = "password"
)

type Kv2SecretsEngineUserPassResponse struct {
	Username string
	Password string
}

func ReadKv2SecretsEngineUserPass(client *api.Client, mountPath string, secretPath string, attributeMap map[string]string) (*Kv2SecretsEngineUserPassResponse, error) {
	resp := &Kv2SecretsEngineUserPassResponse{}

	if mountPath == "" {
		mountPath = defaultKv2SecretEnginePath
	}

	if secretPath == "" {
		return resp, errors.New("empty rolename provided to secrets engine reader")
	}

	// initialize the attributeMap if it wasnt defined in the config
	if attributeMap == nil {
		attributeMap = make(map[string]string)
	}

	am := parseAttributeMap(attributeMap)
	if len(am) == 0 {
		return resp, fmt.Errorf("unable to parse attribute map")
	}

	vaultSecretPath := fmt.Sprintf("%s/data/%s", mountPath, secretPath)

	secret, err := client.Logical().Read(vaultSecretPath)
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", vaultSecretPath))
	}

	if secret == nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", vaultSecretPath))
	}

	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to extract secrets from kv2 response")
	}

	resp.Username, ok = m[am["username"]].(string)
	if !ok {
		return nil, fmt.Errorf("unable to extract username from kv2 response")
	}

	resp.Password, ok = m[am["password"]].(string)
	if !ok {
		return nil, fmt.Errorf("unable to extract password from kv2 response")
	}

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

func ReadKv2SecretsEngineGeneric(client *api.Client, mountPath string, secretPath string, attributeMap map[string]string) (map[string]string, error) {
	resp := make(map[string]string)

	if mountPath == "" {
		mountPath = defaultKv2SecretEnginePath
	}

	if secretPath == "" {
		return resp, errors.New("empty rolename provided to secrets engine reader")
	}

	if attributeMap == nil {
		return resp, errors.New("empty attributeMap provided to secrets engine reader")
	}

	vaultSecretPath := fmt.Sprintf("%s/data/%s", mountPath, secretPath)

	secret, err := client.Logical().Read(vaultSecretPath)
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", vaultSecretPath))
	}

	if secret == nil {
		return resp, fmt.Errorf("failed to read secret from Vault at %s", vaultSecretPath)
	}

	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to extract secrets from kv2 response")
	}

	for k, v := range attributeMap {
		resp[k], ok = m[v].(string)
		if !ok {
			return nil, fmt.Errorf("unable to extract %s=%v from kv2 response", k, v)
		}
		if resp[k] == "" {
			return nil, fmt.Errorf("unable to extract %s=%v from kv2 response", k, v)
		}
	}

	return resp, nil
}
