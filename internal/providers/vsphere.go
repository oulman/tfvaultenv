package providers

import (
	"errors"
	"fmt"
)

func SetVsphereEnv(username string, password string, extraEnvVars map[string]string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("vsphere username or password are empty")
	}

	fmt.Printf("VSPHERE_USERNAME=%s\n", username)
	fmt.Printf("VSPHERE_PASSWORD=%s\n", password)

	for k, v := range extraEnvVars {
		fmt.Printf("%s=%s\n", k, v)
	}
	return "", nil
}
