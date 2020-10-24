package config

import "github.com/spf13/viper"

type Database interface {
	Username() string
	Password() string
	Address() string
	Port() int
	Database() string
}

func (v viperStore) Username() string {
	return viper.GetString("database.username")
}

func (v viperStore) Password() string {
	return viper.GetString("database.password")
}

func (v viperStore) Address() string {
	return viper.GetString("database.address")
}

func (v viperStore) Port() int {
	return viper.GetInt("database.port")
}

func (v viperStore) Database() string {
	return viper.GetString("database.database")
}
