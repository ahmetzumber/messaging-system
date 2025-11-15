package config

import (
	"time"

	"github.com/spf13/viper"
	"github.com/yudai/pp"
)

type Config struct {
	AppName   string
	Server    *Server
	Client    *Client
	Mongo     *Mongo
	Redis     *Redis
	Processor *Processor
}

type Server struct {
	Port int
}

type Client struct {
	URL    string
	ApiKey string
}

type Mongo struct {
	URI               string
	Database          string
	MessageCollection string
}

type Redis struct {
	URI      string
	Password string
	DB       int
	TTL      time.Duration
}

type Processor struct {
	BatchSize int
}

func NewConfig(configPath, configName string) (Config, error) {
	config := Config{}

	viperConfig, err := readConfig(configPath, configName)
	if err != nil {
		return config, err
	}

	if err := viperConfig.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func readConfig(configPath, configName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Config) Print() {
	_, err := pp.Println(c)
	if err != nil {
		return
	}
}
