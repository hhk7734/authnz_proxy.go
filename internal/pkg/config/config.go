package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CONFIG_KEY = "config"
)

func Load(pFlagSets ...*pflag.FlagSet) {
	// pflag
	for _, f := range pFlagSets {
		pflag.CommandLine.AddFlagSet(f)
	}

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env
	viper.AutomaticEnv()

	// config file
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	workDir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(workDir, ".env")); err == nil {
			viper.AddConfigPath(workDir)
			break
		}
		if workDir == "/" {
			break
		}
		workDir = filepath.Dir(workDir)
	}

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
}
