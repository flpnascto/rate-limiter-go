package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	MaxIpRequests    int
	MaxTokenRequests int
	CleanupInterval  string
	WebServerPort    string
	IpBlockTime      int
	TokenBlockTime   int
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
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
