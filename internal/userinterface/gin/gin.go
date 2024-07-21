package gin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhk7734/authnz_proxy.go/internal/userinterface/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	PORT_KEY = "port"
)

type Config struct {
	Port string
}

func PFlagSet() *pflag.FlagSet {
	f := pflag.NewFlagSet("gin", pflag.ContinueOnError)
	f.String(PORT_KEY, "8080", "port")
	return f
}

func ConfigFromViper() Config {
	return Config{
		Port: viper.GetString(PORT_KEY),
	}
}

type GinRestAPI struct {
	engin  *gin.Engine
	server *http.Server
}

func NewGinRestAPI(cfg Config) *GinRestAPI {
	zap.L().Info("gin rest api config", zap.Dict("config",
		zap.String(PORT_KEY, cfg.Port),
	))

	lm := &middleware.GinLoggerMiddleware{}

	engin := gin.New()

	engin.RemoteIPHeaders = append([]string{"X-Envoy-External-Address"}, engin.RemoteIPHeaders...)
	engin.Use(lm.Logger([]string{}))
	engin.Use(lm.Recovery)
	engin.Use(middleware.GinRequestIDMiddleware(true))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engin,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &GinRestAPI{
		engin:  engin,
		server: server,
	}
}

func (g *GinRestAPI) Run() error {
	return g.server.ListenAndServe()
}

func (g *GinRestAPI) Shutdown() error {
	return g.server.Shutdown(context.Background())
}
