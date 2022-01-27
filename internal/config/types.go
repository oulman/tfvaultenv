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

import "reflect"

type Config struct {
	Auth     []*Auth     `hcl:"auth,block"`
	Ad       []*Ad       `hcl:"ad,block"`
	Aws      *Aws        `hcl:"aws,block"`
	Azure    *Azure      `hcl:"azure,block"`
	KvSecret []*KvSecret `hcl:"kv_secret,block"`
}

type Auth struct {
	Name     string   `hcl:",label"`
	Method   string   `hcl:"method"`
	Path     string   `hcl:"path"`
	When     *When    `hcl:"when,block"`
	Priority int      `hcl:"priority,optional"`
	Jwt      *AuthJwt `hcl:"jwt,block"`
}

func (a *Auth) IsEmpty() bool {
	return reflect.DeepEqual(a, Auth{})
}

type AuthJwt struct {
	Role  string `hcl:"role"`
	Token string `hcl:"token"`
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

type Azure struct {
	Name         string            `hcl:",label"`
	Role         string            `hcl:"role"`
	ExtraEnvVars map[string]string `hcl:"extra_env_vars,optional"`
	Mount        string            `hcl:"mount,optional"`
}

func (a *Azure) IsEmpty() bool {
	return reflect.DeepEqual(a, Azure{})
}

type Ad struct {
	Name           string            `hcl:",label"`
	Role           string            `hcl:"role"`
	Mount          string            `hcl:"mount,optional"`
	TargetProvider string            `hcl:"target_provider"`
	UsernameEnvVar string	 		 `hcl:"username_env_var,optional"`
	PasswordEnvVar string	 		 `hcl:"password_env_var,optional"`
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
	EnvPresent string `hcl:"env_present,optional"`
}
