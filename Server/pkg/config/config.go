package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

type Config struct {
	HTTPServer `yaml:"server"`
	Enviroment string `yaml:"env"`
	Database   `yaml:"database"`
	Auth       `yaml:"auth"`
}
type Auth struct {
	JwtToken string `yaml:"jwttoken"`
}

type Database struct {
	Url string `yaml:"url"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Path to config yml file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("ERROR :: PLEASE PROVIDE CONFIG PATH")
		}
		_, err := os.Stat(configPath)

		if err != nil {
			log.Fatalf("Config file doesn't exist on this path %s", configPath)
		}

	}
	var config Config

	err := cleanenv.ReadConfig(configPath, &config)

	fmt.Println(config)

	if err != nil {
		log.Fatalf("ERROR :: Can't read config file :: %s", err.Error())
	}
	return &config

}
