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
