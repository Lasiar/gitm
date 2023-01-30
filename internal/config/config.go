package config

import (
	"runtime"

	"github.com/spf13/viper"
)

type Commit struct {
	Branch        string
	CommitMessage string
}

type Config struct {

	// Repositories URLs of repositories.
	Repositories  []string
	Script        []string
	ParallelLimit int
	Commit        *Commit
}

func LoadConfig(configPath string) (Config, error) {
	c := Config{
		ParallelLimit: runtime.NumCPU(),
	}

	v := viper.New()
	v.SetConfigFile(configPath)

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return c, v.Unmarshal(&c)
}
