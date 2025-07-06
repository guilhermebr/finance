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
	//AuthSecretKey  string `conf:"env:AUTH_SECRET_KEY,required"`
	Service struct {
		Address string `conf:"env:SERVICE_ADDRESS,default:0.0.0.0:3000"`
	}
	Web struct {
		Address    string `conf:"env:WEB_ADDRESS,default:0.0.0.0:8080"`
		ApiBaseURL string `conf:"env:API_BASE_URL,default:http://127.0.0.1:3000"`
	}
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
