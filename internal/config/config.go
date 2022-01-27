package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty/function"

	providers "github.com/oulman/tfvaultenv/internal/providers"
	vaulthelper "github.com/oulman/tfvaultenv/internal/vaulthelper"
)

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

	// add hcl functions
	ectx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			"env": EnvFunc,
		},
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("config parse: %w", diags)
	}

	c = &Config{}

	diags = gohcl.DecodeBody(file.Body, ectx, c)
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

	// add an implicit empty auth which will default to env/credential helpers
	implicitAuth := &Auth{
		Priority: 0,
	}
	authMethods := []*Auth{}
	authMethods = append(authMethods, implicitAuth)

	// loop over Auth blocks and build a slice of valid entries based on conditions
	// TODO: refactor to support more conditions
	for _, v := range c.Auth {
		if !v.IsEmpty() {
			if v.When != nil && v.When.EnvPresent != "" {
				env := os.Getenv(v.When.EnvPresent)
				if env != "" {
					authMethods = append(authMethods, v)
				}
			} else {
				authMethods = append(authMethods, v)
			}
		}
	}

	// sort auth methods by priority
	sort.Slice(authMethods, func(i, j int) bool {
		return authMethods[i].Priority > authMethods[j].Priority
	})

	auth := authMethods[0]

	switch auth.Method {
	case "jwt":
		if auth.Jwt == nil {
			return fmt.Errorf("config error: jwt{} block not set for auth.name = %s", auth.Name)
		}
		secret, err := vaulthelper.AuthJwt(client, auth.Path, auth.Jwt.Role, auth.Jwt.Token)
		if err != nil {
			return errors.Wrap(err, "failed to authenticate against JWT endpoint")
		}

		if secret.Auth.ClientToken != "" {
			client.SetToken(secret.Auth.ClientToken)
		} else {
			return fmt.Errorf("no ClientToken in JWT response")
		}
	// fall back to credentials helper / env
	default:
		err = vaulthelper.SetVaultTokenFromEnvOrHandler(client)
		if err != nil {
			return errors.Wrap(err, "failed to set client token")
		}
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
		case "generic":
			secret, err := vaulthelper.ReadAdSecretsEngine(client, v.Role, v.Mount)
			if err != nil {
				return errors.Wrap(err, "reading Vault Ad secrets engine")
			}

			if v.UsernameEnvVar == "" || v.PasswordEnvVar == "" {
				return fmt.Errorf("username_env_var and password_env_var are both required for AD generic provider")
			}

			// build a secret map between the desired environment variables
			// and the AD-provided values
			secretMap := map[string]string {
				v.UsernameEnvVar: secret.Username,
				v.PasswordEnvVar: secret.CurrentPassword,
			}

			_, err = providers.SetGenericEnv(secretMap, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set generic environment variables")
			}
		default:
			return fmt.Errorf("invalid target_provider for engine ad: %s", v.TargetProvider)
		}
	}

	// loop over Kv secrets engine entries and process them
	for _, v := range c.KvSecret {
		switch v.TargetProvider {
		case "vsphere":
			secret, err := vaulthelper.ReadKv2SecretsEngineUserPass(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetVsphereEnv(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set vsphere environment variables")
			}
		case "infoblox":
			secret, err := vaulthelper.ReadKv2SecretsEngineUserPass(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetInfobloxEnv(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set Infoblox environment variables")
			}
		case "f5":
			secret, err := vaulthelper.ReadKv2SecretsEngineUserPass(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetF5Env(secret.Username, secret.Password, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set f5 environment variables")
			}
		case "generic":
			secretMap, err := vaulthelper.ReadKv2SecretsEngineGeneric(client, v.Mount, v.Path, v.AttributeMap)
			if err != nil {
				return errors.Wrap(err, "reading Vault kv2 secrets engine")
			}

			_, err = providers.SetGenericEnv(secretMap, v.ExtraEnvVars)
			if err != nil {
				return errors.Wrap(err, "failed to set generic environment variables")
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
		default:
			return fmt.Errorf("invalid method %s for AWS secrets engine", c.Aws.Method)
		}
	}

	if c.Azure != nil {
		secret, err := vaulthelper.ReadAzureSecretsEngine(client, c.Azure.Mount, c.Azure.Role)
		if err != nil {
			return errors.Wrap(err, "reading Vault Azure secrets engine")
		}

		_, err = providers.SetAzureEnv(*secret, c.Azure.ExtraEnvVars)
		if err != nil {
			return errors.Wrap(err, "failed to set Azure environment variables")
		}
	}

	return nil
}
