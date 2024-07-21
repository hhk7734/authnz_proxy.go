package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CONFIG_KEY = "config"
)

func Load(pFlagSets ...*pflag.FlagSet) {
	// pflag
	pflag.CommandLine.String(CONFIG_KEY, "config.yaml", "config file path")
	for _, f := range pFlagSets {
		pflag.CommandLine.AddFlagSet(f)
	}

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env
	for _, env := range []string{".env", "../.env", "../../.env"} {
		// If file exists, load it
		if _, err := os.Stat(env); err == nil {
			if err := godotenv.Load(env); err != nil {
				panic(fmt.Errorf("failed to load env file: %w", err))
			}
			break
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// config file
	configPath := viper.GetString(CONFIG_KEY)
	dir := filepath.Dir(configPath)
	base := filepath.Base(configPath)
	ext := filepath.Ext(base)
	viper.AddConfigPath(dir)
	viper.SetConfigName(base)
	if len(ext) > 1 {
		viper.SetConfigType(ext[1:])
	}

	workDir, _ := os.Getwd()
	for {
		if workDir == "/" {
			break
		}
		viper.AddConfigPath(workDir)
		workDir = filepath.Dir(workDir)
	}

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			fmt.Printf(`{"level":"warn","msg":"config file not found", "path":"%s"}`+"\n", configPath)
		} else {
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
	}
}
