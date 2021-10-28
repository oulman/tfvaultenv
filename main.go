package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	config "github.com/oulman/tfvaultenv/internal/config"
	"github.com/pkg/errors"
)

const (
	configFileName = ".tfvaultenv.config.hcl"
	defaultDepth   = 0
)

func main() {
	if err := inner(); err != nil {
		fmt.Printf("tfvaultenv error: %s\n", err.Error())
		os.Exit(1)
	}
}

func inner() error {

	// There is a random function for the HCL configuration.
	rand.Seed(time.Now().Unix())

	depth := defaultDepth
	d := os.Getenv("TFVAULENV_CONFIG_DEPTH")
	if d != "" {
		v, err := strconv.Atoi(d)
		if err != nil {
			depth = v
		} else {
			return errors.Wrap(err, "invalid TFVAULENV_CONFIG_DEPTH")
		}
	}

	configFilePath, err := config.FindConfigFile(depth, configFileName)
	if err != nil {
		return err
	}

	configParsed, err := config.ParseConfig(configFilePath)
	if err != nil {
		return err
	}

	err = config.ProcessConfig(configParsed)
	if err != nil {
		return err
	}

	return nil
}
