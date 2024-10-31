package config

import "github.com/spf13/viper"

type Config struct {
	MaxRequestsPerSecond      int
	BlockDuration             int
	MaxRequestsPerSecondToken int
	BlockDurationToken        int
	RedisAddress              string
	RedisPassword             string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		MaxRequestsPerSecond:      viper.GetInt("MAX_REQUESTS_PER_SECOND"),
		BlockDuration:             viper.GetInt("BLOCK_DURATION"),
		MaxRequestsPerSecondToken: viper.GetInt("MAX_REQUESTS_PER_SECOND_TOKEN"),
		BlockDurationToken:        viper.GetInt("BLOCK_DURATION_TOKEN"),
		RedisAddress:              viper.GetString("REDIS_ADDRESS"),
		RedisPassword:             viper.GetString("REDIS_PASSWORD"),
	}

	return config, nil
}
