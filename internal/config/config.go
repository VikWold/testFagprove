package config

import "github.com/spf13/viper"

type Config struct {
	DB *DatabaseConfig `json:db`
}

type DatabaseConfig struct {
	DSN     string `json:"-"`
	Timeout int    `json:"timeout"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/")
	viper.AddConfigPath("$HOME/.config/testFagprove")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.BindEnv("db.dsn", "TESTFAGPROVE_DB_DSN")
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
