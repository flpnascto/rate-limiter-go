package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	MaxIpRequests    string
	MaxTokenRequests string
	CleanupInterval  string
	WebServerPort    string
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
