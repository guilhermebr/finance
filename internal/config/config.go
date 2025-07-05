package config

import (
	"errors"
	"fmt"

	"github.com/ardanlabs/conf/v3"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment    string `conf:"env:ENVIRONMENT,default:development"`
	DatabaseEngine string `conf:"env:DATABASE_ENGINE,default:postgres"`
	ApiAddress     string `conf:"env:API_ADDRESS,default:0.0.0.0:8000"`
	//AuthSecretKey  string `conf:"env:AUTH_SECRET_KEY,required"`
}

func (c *Config) Load(prefix string) error {
	if help, err := conf.Parse(prefix, c); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return err
		}
		return err
	}
	return nil
}
