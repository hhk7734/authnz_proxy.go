package logger

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	LOG_LEVEL_KEY  = "log_level"
	LOG_FORMAT_KEY = "log_format"
)

type Config struct {
	Level  string
	Format string
}

func PFlagSet() *pflag.FlagSet {
	f := pflag.NewFlagSet("log", pflag.ContinueOnError)
	f.String(LOG_LEVEL_KEY, "info", "log level")
	f.String(LOG_FORMAT_KEY, "json", "log format")
	return f
}

func ConfigFromViper() Config {
	return Config{
		Level:  viper.GetString(LOG_LEVEL_KEY),
		Format: viper.GetString(LOG_FORMAT_KEY),
	}
}

func SetGlobalZapLogger(cfg Config) {
	var l *zap.Logger
	var zapCfg zap.Config

	if cfg.Format != "json" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.TimeKey = "time"
	}

	err := zapCfg.Level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(err)
	}

	l, _ = zapCfg.Build()
	defer l.Sync()
	zap.ReplaceGlobals(l)

	zap.L().Info("logger config", zap.Dict("config",
		zap.String(LOG_LEVEL_KEY, cfg.Level),
		zap.String(LOG_FORMAT_KEY, cfg.Format),
	))
}
