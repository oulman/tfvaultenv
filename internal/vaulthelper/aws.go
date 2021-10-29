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
	defaultAwsSecretsEnginePath   = "aws"
	defaultAwsStsSecretsEngineTtl = "900s"
)

type AwsStsSecretsEngineResponse struct {
	AccessKey     string
	SecretKey     string
	SecurityToken string
}

func (r *AwsStsSecretsEngineResponse) IsPopulated() bool {
	return r.AccessKey != "" && r.SecretKey != "" && r.SecurityToken != ""
}

func ReadAwsStsSecretsEngine(client *api.Client, mount string, role string, rolearn string, ttl string) (*AwsStsSecretsEngineResponse, error) {
	resp := &AwsStsSecretsEngineResponse{}

	if mount == "" {
		mount = defaultAwsSecretsEnginePath
	}

	if role == "" {
		return resp, errors.New("empty role provided to Aws Sts secrets engine reader")
	}

	if rolearn == "" {
		return resp, errors.New("empty role_arn provided to Aws Sts secrets engine reader")
	}

	if ttl == "" {
		ttl = defaultAwsStsSecretsEngineTtl
	}

	path := fmt.Sprintf("%s/sts/%s", mount, role)

	secret, err := client.Logical().Write(path, map[string]interface{}{
		"name":     role,
		"role_arn": rolearn,
		"ttl":      ttl,
	})
	if err != nil {
		return resp, errors.Wrap(err, fmt.Sprintf("failed to read secret from Vault at %s\n", path))
	}

	if secret == nil {
		return resp, errors.Wrap(err, fmt.Sprintf("no value from Vault at %s\n", path))
	}

	resp.AccessKey = secret.Data["access_key"].(string)
	resp.SecretKey = secret.Data["secret_key"].(string)
	resp.SecurityToken = secret.Data["security_token"].(string)

	return resp, nil
}
