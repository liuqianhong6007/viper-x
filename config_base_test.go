package viperx

import (
	"fmt"
	"testing"
)

type Config struct {
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Server struct {
	Port        int    `yaml:"port"`
	PassportKey string `yaml:"passport_key"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func TestLoad(t *testing.T) {
	var conf Config
	Load(&conf)
	fmt.Println(conf)
}
