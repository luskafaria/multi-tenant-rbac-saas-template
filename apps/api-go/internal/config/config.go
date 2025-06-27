package config

import "github.com/spf13/viper"

type Config struct {
	DB_SOURCE                    string `mapstructure:"DB_SOURCE"`
	API_PORT                     string `mapstructure:"API_PORT"`
	JWT_SECRET                   string `mapstructure:"JWT_SECRET"`
	GITHUB_OAUTH_CLIENT_ID       string `mapstructure:"GITHUB_OAUTH_CLIENT_ID"`
	GITHUB_OAUTH_CLIENT_SECRET   string `mapstructure:"GITHUB_OAUTH_CLIENT_SECRET"`
	GITHUB_OAUTH_CLIENT_REDIRECT_URI string `mapstructure:"GITHUB_OAUTH_CLIENT_REDIRECT_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
