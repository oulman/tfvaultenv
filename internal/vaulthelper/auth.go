package vaulthelper

// import (
// 	"fmt"

// 	"github.com/hashicorp/go-hclog"
// 	"github.com/hashicorp/go-secure-stdlib/awsutil"
// )

// const (
// 	accessKey    = ""
// 	secretKey    = ""
// 	sessionToken = ""
// 	headerValue  = ""
// )

// // AuthAws() returns
// func AuthAws() (string, error) {

// 	// the Vault awsutil helpers require a logger but for our purposes
// 	// we can turn it off
// 	level := hclog.LevelFromString("Off")
// 	hlogger := hclog.Default()
// 	hlogger.SetLevel(level)

// 	creds, err := awsutil.RetrieveCreds(accessKey, secretKey, sessionToken, hlogger)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// default to us-east-1
// 	region := awsutil.DefaultRegion

// 	loginData, err := awsutil.GenerateLoginData(creds, headerValue, region, hlogger)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if loginData == nil {
// 		return nil, fmt.Errorf("got nil response from GenerateLoginData")
// 	}
// 	loginData["role"] = role
// 	path := fmt.Sprintf("auth/%s/login", mount)
// 	secret, err := c.Logical().Write(path, loginData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if secret == nil {
// 		return nil, fmt.Errorf("empty response from credential provider")
// 	}

// 	return secret, nil
// }
