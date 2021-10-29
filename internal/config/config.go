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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"

	providers "github.com/oulman/tfvaultenv/internal/providers"
	vaulthelper "github.com/oulman/tfvaultenv/internal/vaulthelper"
)

type Config struct {
	Auth     []*Auth     `hcl:"auth,block"`
	Ad       []*Ad       `hcl:"ad,block"`
	Aws      *Aws        `hcl:"aws,block"`
	KvSecret []*KvSecret `hcl:"kv_secret,block"`
}

type Auth struct {
	Name   string `hcl:",label"`
	Method string `hcl:"method"`
	Path   string `hcl:"path"`
	When   *When  `hcl:"when,block"`
}

type Aws struct {
	Name         string            `hcl:",label"`
	Method       string            `hcl:"method"`
	Role         string            `hcl:"role"`
	RoleArn      string            `hcl:"role_arn,optional"`
	Ttl          string            `hcl:"ttl"`
	ExtraEnvVars map[string]string `hcl:"extra_env_vars,optional"`
	Mount        string            `hcl:"mount,optional"`
}

func (a *Aws) IsEmpty() bool {
	return reflect.DeepEqual(a, Aws{})
}

type Ad struct {
	Name           string            `hcl:",label"`
	Role           string            `hcl:"role"`
	Mount          string            `hcl:"mount,optional"`
	TargetProvider string            `hcl:"target_provider"`
	ExtraEnvVars   map[string]string `hcl:"extra_env_vars,optional"`
}

type KvSecret struct {
	Name           string            `hcl:",label"`
	Path           string            `hcl:"path"`
	Mount          string            `hcl:"mount,optional"`
	TargetProvider string            `hcl:"target_provider"`
	AttributeMap   map[string]string `hcl:"attribute_map,optional"`
	ExtraEnvVars   map[string]string `hcl:"extra_env_vars,optional"`
}

type When struct {
	EnvPresent string `hcl:"env_present"`
}

// Parse will parse file content into valid config.
func ParseConfig(filename string) (c *Config, err error) {
	// First, open a file handle to the input filename.
	input, err := os.Open(filename)
	if err != nil {
		return &Config{}, fmt.Errorf(
			"error in ReadConfig opening config file: %w", err,
		)
	}
	defer input.Close()

	// Next, read that file into a byte slice for use as a buffer. Because HCL
	// decoding must happen in the context of a whole file, it does not take an
	// io.Reader as an input, instead relying on byte slices.
	src, err := ioutil.ReadAll(input)
	if err != nil {
		return &Config{}, fmt.Errorf(
			"error in ReadConfig reading input `%s`: %w", filename, err,
		)
	}
	var diags hcl.Diagnostics

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("config parse: %w", diags)
	}

	c = &Config{}

	diags = gohcl.DecodeBody(file.Body, nil, c)
	if diags.HasErrors() {
		return nil, fmt.Errorf("config parse: %w", diags)
	}

	return c, nil
}

func ProcessConfig(c *Config) error {
	client, err := vaulthelper.GetVaultApiClient()
	if err != nil {
		return errors.Wrap(err, "failed to setup vault client")
	}

	// loop over Ad secrets engine entries and process them
	for _, v := range c.Ad {
		switch v.TargetProvider {
		case "vsphere":
			secret, err := vaulthelper.ReadAdSecretsEngine(client, v.Role, v.Mount)
			if err != nil {
				return errors.Wrap(err, "reading Vault Ad secrets engine")
			}

			_, err = providers.SetVsphereEnv(secret.Username, secret.CurrentPassword, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set vsphere environment variables")
			}
		case "infoblox":
			secret, err := vaulthelper.ReadAdSecretsEngine(client, v.Role, v.Mount)
			if err != nil {
				return errors.Wrap(err, "reading Vault Ad secrets engine")
			}

			_, err = providers.SetInfobloxEnv(secret.Username, secret.CurrentPassword, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set infoblox environment variables")
			}
		case "f5":
			secret, err := vaulthelper.ReadAdSecretsEngine(client, v.Role, v.Mount)
			if err != nil {
				return errors.Wrap(err, "reading Vault Ad secrets engine")
			}

			_, err = providers.SetF5Env(secret.Username, secret.CurrentPassword, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set F5 environment variables")
			}
		default:
			return fmt.Errorf("invalid target_provider for engine ad: %s", v.TargetProvider)
		}
	}

	// loop over Kv secrets engine entries and process them
	for _, v := range c.KvSecret {
		switch v.TargetProvider {
		case "vsphere":
			secret, err := vaulthelper.ReadKv2SecretsEngine(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetVsphereEnv(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set vsphere environment variables")
			}
		case "infoblox":
			secret, err := vaulthelper.ReadKv2SecretsEngine(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetInfobloxEnv(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set Infoblox environment variables")
			}
		case "f5":
			secret, err := vaulthelper.ReadKv2SecretsEngine(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetF5Env(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set vsphere environment variables")
			}
		default:
			return fmt.Errorf("invalid target_provider for engine kv_secret: %s", v.TargetProvider)
		}
	}

	if c.Aws != nil && !c.Aws.IsEmpty() {
		switch c.Aws.Method {
		case "assumed_role":
			secret, err := vaulthelper.ReadAwsStsSecretsEngine(client, c.Aws.Mount, c.Aws.Role, c.Aws.RoleArn, c.Aws.Ttl)
			if err != nil {
				return errors.Wrap(err, "reading Vault AWS secrets engine")
			}

			_, err = providers.SetAwsStsEnv(*secret, c.Aws.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set AWS STS environment variables")
			}
		}
	}
	return nil
}
