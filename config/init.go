package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Port           string
	HashSalt       string
	SigningKey     string
	TokenTtl       time.Duration
	RedisPort      string
	Sendler        string
	Password       string
	HostForSending string
	PortForSending string
}

func InitConfig() (Config, error) {
	var bytes, err = os.ReadFile("./config/config.json")
	if err != nil {
		return Config{}, err
	}
	var v Config
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return Config{}, err
	}
	return v, nil
}
