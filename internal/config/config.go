package config

import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig() (*Config, error) {

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "test")
	viper.SetDefault("database.password", "test")
	viper.SetDefault("database.dbname", "testdb")
	viper.SetDefault("database.sslmode", "disable")

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}
