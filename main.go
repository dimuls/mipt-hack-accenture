package main

import (
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type config struct {
	PostgresURI string `yaml:"postgres_uri"`
}

func main() {
	if len(os.Args) != 2 {
		logrus.Fatal("expected exact one argument: path to config file")
	}

	cYAML, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	var c config

	err = yaml.Unmarshal(cYAML, &c)
	if err != nil {
		logrus.WithError(err).Fatal("failed to YAML unmarshal config")
	}

	db, err := sqlx.Connect("postgres", c.PostgresURI)
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to postgres")
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// TODO
	//e.GET("/resource-group", func(c echo.Context) error {
	//
	//})

}
