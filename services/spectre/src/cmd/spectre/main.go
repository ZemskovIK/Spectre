package main

/*
	Spectre is a REST-api based service for working with letters and users.
	Service includes jwt-based auth and working with enryption (look at crypto service).
*/

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	server "spectre/internal/srv"
	"spectre/internal/srv/proxy"
	"spectre/internal/storage/sqlite"
	"spectre/pkg/config"
	"spectre/pkg/logger"
	"syscall"
	"time"
)

const (
	logLevel = "DEBUG"
	dbPath   = "../db/spec.db"
)

func main() {
	// init logger
	log := logger.New(logLevel)

	// init config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init storage
	st, err := sqlite.NewStorage(dbPath, log)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Info("init storage")

	// init crypto client
	cryptoADDR := fmt.Sprintf("%s:%s", cfg.Server.ProxyHost, cfg.Server.ProxyPort)
	crypto := proxy.NewCryptoClient( // ! TODO ip!!!
		cryptoADDR+cfg.Routes.ProxyEncryptPoing,
		cryptoADDR+cfg.Routes.ProxyDecryptPoing,
		cryptoADDR+cfg.Routes.ProxyECDHPoing,
	// "http://127.0.0.1:7654/encrypt",
	// "http://127.0.0.1:7654/decrypt",
	// "http://127.0.0.1:7654/ecdh",
	)
	log.Infof("ready proxy at %s", cryptoADDR)

	// init router
	r := server.NewRouter(st, log, crypto)
	r.Use(server.AuthMiddleware(log))
	r.Use(server.CORSMiddleware)
	r.Use(server.JSONRespMiddleware)
	log.Info("init router")

	// init server
	srvADDR := fmt.Sprintf("%s:%s", cfg.Server.SpectreHost, cfg.Server.SpectrePort)
	srv := &http.Server{
		Addr: srvADDR,
		// Addr:    ":5000",
		Handler: r,
	}
	log.Info("init server")

	// run
	go func() {
		log.Infof("running server on %s", srvADDR)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cannot run server on :%s : %v", err, srvADDR)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Info("server exited gracefully :)")
}
