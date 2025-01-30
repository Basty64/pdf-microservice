package options

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Api Api
	S3  S3
}

type Api struct {
	Name      string `mapstructure:"name"`
	Port      string `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	LocalSave bool   `mapstructure:"local_save"`
}

type S3 struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	UseSSL          bool   `mapstructure:"useSSL"`
	BucketName      string `mapstructure:"bucketName"`
	Region          string `mapstructure:"region"`
	FilePath        string `mapstructure:"filePath"`
	ObjectKey       string `mapstructure:"objectKey"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath) // Указываем путь к config.toml
	viper.SetConfigType("toml")

	viper.AutomaticEnv() // Чтение переменных окружения

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
