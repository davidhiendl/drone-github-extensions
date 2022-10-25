package main

import (
	"dhswt.de/drone-github-extensions/plugin_converter"
	"dhswt.de/drone-github-extensions/plugin_env"
	"dhswt.de/drone-github-extensions/shared"
	"fmt"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/drone-go/plugin/environ"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg := new(shared.AppConfig)
	err := envconfig.Process("", cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if cfg.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if cfg.Bind == "" {
		cfg.Bind = ":3000"
	}

	environHandler := environ.Handler(
		cfg.Secret,
		plugin_env.New(cfg),
		logrus.StandardLogger(),
	)
	http.Handle("/env", environHandler)

	converterHandler := converter.Handler(
		plugin_converter.New(cfg),
		cfg.Secret,
		logrus.StandardLogger(),
	)
	http.Handle("/convert", converterHandler)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "OK")
	})

	logrus.Infof("server listening on address %s", cfg.Bind)
	logrus.Fatal(http.ListenAndServe(cfg.Bind, nil))
}
