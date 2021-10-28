package providers

import (
	"fmt"

	"github.com/pkg/errors"
)

func SetF5Env(username string, password string, extraEnvVars map[string]string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("F5 username or password are empty")
	}

	fmt.Printf("BIGIP_USER=%s\n", username)
	fmt.Printf("BIGIP_PASSWORD=%s\n", password)

	for k, v := range extraEnvVars {
		fmt.Printf("%s=%s\n", k, v)
	}
	return "", nil
}
