package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config - структура конфигурации проекта
type Config struct {
	Address        string `mapstructure:"ADDRESS"`
	GrpcAddress    string `mapstructure:"GRPC_ADDRESS"`
	DatabaseDSN    string `mapstructure:"DATABASE_DSN"`
	SecretKey      string `mapstructure:"SECRET_KEY"`
	PasswordSecret string `mapstructure:"PASSWORD_SECRET"`
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`
}

// LoadConfig - загрузка конфига проекта
func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("ADDRESS", "localhost:8080")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Infof("filed to find conf file, set default value: %v.", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
