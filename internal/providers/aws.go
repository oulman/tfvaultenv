package providers

import (
	"fmt"

	"github.com/oulman/tfvaultenv/internal/vaulthelper"
	"github.com/pkg/errors"
)

func SetAwsStsEnv(secret vaulthelper.AwsStsSecretsEngineResponse, extraEnvVars map[string]string) (string, error) {

	if !secret.IsPopulated() {
		return "", errors.New("AWS sts secret is not valid")
	}

	fmt.Printf("AWS_ACCESS_KEY_ID=%s\n", secret.AccessKey)
	fmt.Printf("AWS_ACCESS_SECRET_KEY=%s\n", secret.SecretKey)
	fmt.Printf("AWS_ACCESS_SESSION_TOKEN=%s\n", secret.SecurityToken)

	for k, v := range extraEnvVars {
		fmt.Printf("%s=%s\n", k, v)
	}
	return "", nil
}
