package doconf

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Config Configuration

const CONFIG_PATH = "config.yaml"

type Configuration struct {
	Version string `yaml:"version"`

	Service struct {
		Web struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"web"`

		DB struct {
			Host    string `yaml:"host"`
			Port    int    `yaml:"port"`
			SSLMode string `yaml:"sslmode"`

			Environment struct {
				PostgresDB       string `yaml:"POSTGRES_DB"`
				PostgresUser     string `yaml:"POSTGRES_USER"`
				PostgresPassword string `yaml:"POSTGRES_PASSWORD"`
			} `yaml:"environment"`

			Settings struct {
				MaxIdleTime  int `yaml:"MAX_IDLE_TIME"`
				MaxConnLife  int `yaml:"MAX_CONN_LIFE"`
				MaxOpenConns int `yaml:"MAX_OPEN_CONNS"`
				MaxIdleConns int `yaml:"MAX_IDLE_CONNS"`
			} `yaml:"settings"`
		} `yaml:"db"`
	}
}

func init() {
	// Read YAML file
	file, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		msg := "Could not read file %v"
		log.Fatalf(msg, CONFIG_PATH)
	}

	// Parse YAML into struct
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		msg := "Could not parse %v file"
		log.Fatalf(msg, CONFIG_PATH)
	}
}
