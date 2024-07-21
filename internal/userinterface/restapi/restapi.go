package restapi

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

type RestAPI struct {
	engin  *gin.Engine
	server *http.Server
}

func NewRestAPI(cfg Config) *RestAPI {
	zap.L().Info("rest api config", zap.Dict("config",
		zap.String(PORT_KEY, cfg.Port),
	))

	lm := &middleware.Logger{}

	engin := gin.New()

	engin.RemoteIPHeaders = append([]string{"X-Envoy-External-Address"}, engin.RemoteIPHeaders...)
	engin.Use(lm.Logger([]string{}))
	engin.Use(lm.Recovery)
	engin.Use(middleware.RequestID(true))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engin,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &RestAPI{
		engin:  engin,
		server: server,
	}
}

func (g *RestAPI) Run() error {
	return g.server.ListenAndServe()
}

func (g *RestAPI) Shutdown() error {
	return g.server.Shutdown(context.Background())
}
