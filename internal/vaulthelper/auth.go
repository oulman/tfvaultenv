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
