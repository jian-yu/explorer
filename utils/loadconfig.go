package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func LoadConfigPath(relativePath string) string {
	configPath, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("laod config file err:%s.\n", err.Error())
		return ""
	}

	if configPath == "" {
		panic(`must had config file`)
	}

	fmt.Printf("config file path:%s.\n", configPath)
	return configPath
}

func LoadViper(envPrefix, configPath string) (*viper.Viper, error) {
	if len(envPrefix) > 0 {
		viper.SetEnvPrefix(envPrefix)
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(*os.PathError); !ok {
			return nil, err
		}
	}
	return viper.GetViper(), nil
}
