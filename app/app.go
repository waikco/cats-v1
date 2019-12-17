package app

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/waikco/cats-v1/conf"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/waikco/cats-v1/model"
)

// App ...
type App struct {
	Server  *http.Server
	Storage model.Storage
	Router  http.Handler
	Config  conf.Config
}

// Bootstrap prepares app for run by setting things up based on provided config.
func (a *App) Bootstrap() {
	a.BootstrapLogger()
	storage, err := model.BootstrapPostgres(a.Config.Database)
	if err != nil {
		log.Fatal().Err(err)
	}
	a.Storage = storage
	a.BootstrapServer()
}

func (a *App) BootstrapServer() {
	router := httprouter.New()

	// add actual api routes
	router.GET("/cats/v1/health", a.Health)
	router.GET("/cats/v1/cats/:id", a.GetCat)
	router.GET("/cats/v1/cats", a.GetCats)
	router.POST("/cats/v1/", a.CreateCat)
	//router.POST("/cats/v1/bulkcatadd", a.MassCreateCat)
	router.PUT("/cats/v1/:id", a.UpdateCat)
	router.DELETE("/cats/v1/cats/:id", a.DeleteCat)

	a.Router = router

	cfg := &tls.Config{}
	if a.Config.Server.TLS {
		cert, err := tls.LoadX509KeyPair(
			a.Config.Server.Cert,
			a.Config.Server.Key)

		if err != nil {
			log.Fatal().Msgf("Unable to load cert/key: %s", err)
		}

		cfg = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       false,
			Certificates:             []tls.Certificate{cert},
		}
		cfg.BuildNameToCertificate()
	}

	addr := fmt.Sprintf(":%s", a.Config.Server.Port)
	a.Server = &http.Server{
		Addr:      addr,
		Handler:   a.Router,
		TLSConfig: cfg,
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Info().Msgf("initialized server to listen on: %+v", a.Server.Addr)
}

func (a *App) BootstrapLogger() {
	level := zerolog.InfoLevel
	if l, err := zerolog.ParseLevel(a.Config.Logging.Level); err == nil {
		level = l
	}
	zerolog.SetGlobalLevel(level)

}

//RunApp
func (a *App) Run() {
	log.Fatal().Err(a.Server.ListenAndServe())
}
