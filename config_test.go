package viper_x

import (
	"fmt"
	"testing"
)

type Config struct {
	Server Server `viper:"server"`
	DB     DB     `viper:"db"`
}

type Server struct {
	Port        int    `viper:"port"`
	PassportKey string `viper:"passport_key"`
}

type DB struct {
	Host     string `viper:"host"`
	Port     int    `viper:"port"`
	User     string `viper:"user"`
	Password string `viper:"password"`
	Name     string `viper:"name"`
}

func TestReadConf(t *testing.T) {
	var conf Config

	ReadConf("demo", &conf)

	fmt.Println(conf)
}
