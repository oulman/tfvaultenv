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

package providers

import (
	"fmt"

	"github.com/pkg/errors"
)

func SetInfobloxEnv(username string, password string, extraEnvVars map[string]string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("Infoblox username or password are empty")
	}

	fmt.Printf("INFOBLOX_USERNAME=%s\n", username)
	fmt.Printf("INFOBLOX_PASSWORD=%s\n", password)

	for k, v := range extraEnvVars {
		fmt.Printf("%s=%s\n", k, v)
	}
	return "", nil
}
