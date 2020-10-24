package config

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

type viperStore struct{}

var vs *viperStore
var once sync.Once

type Interface interface {
	Database
}

func New() (Interface, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/packs/")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.Wrap(err, "config file not found")
		} else {
			return nil, errors.Wrap(err, "problem loading config file")
		}
	}

	newInstance := func() {
		vs = &viperStore{}
	}

	once.Do(newInstance)

	return vs, nil
}
