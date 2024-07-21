package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hhk7734/authnz_proxy.go/internal/pkg/config"
	"github.com/hhk7734/authnz_proxy.go/internal/pkg/logger"
	"github.com/hhk7734/authnz_proxy.go/internal/userinterface/restapi"
	"go.uber.org/zap"
)

func main() {
	config.Load(
		logger.PFlagSet(),
		restapi.PFlagSet(),
	)

	logger.SetGlobalZapLogger(logger.ConfigFromViper())

	server := restapi.NewRestAPI(restapi.ConfigFromViper())

	listenErr := make(chan error, 1)
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		zap.L().Error("failed to listen and serve", zap.Error(err))
	case <-shutdown:
	}

	zap.L().Info("shutting down server...")

	wg := &sync.WaitGroup{}

	go func() {
		defer wg.Done()
		// blocked until all connections are closed or timeout
		if err := server.Shutdown(); err != nil {
			zap.L().Error("failed to shutdown server", zap.Error(err))
		}
	}()
	wg.Add(1)

	wg.Wait()
}
