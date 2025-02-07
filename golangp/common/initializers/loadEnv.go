package initializers

import (
	"fmt"
	"hajime/golangp/common/logging"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const packagePath = "golangp/common/initializers"

func FindDevConfigDir(targets []string) string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return ""
	}

	for {
		for _, target := range targets {
			configPath := filepath.Join(wd, packagePath, target)
			if _, err := os.Stat(configPath); err == nil {
				return filepath.Join(wd, packagePath)
			}
		}

		parentDir := filepath.Dir(wd)
		if parentDir == wd || wd == "/" {
			logging.Info("Dev Config file not found")
			return ""
		}

		wd = parentDir
	}
}

func LoadEnv(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath(FindDevConfigDir([]string{"app.env", "app.dev.env"})) // for bazel run
	//Config file names search in order: app.env, app.dev.env
	configNames := []string{"app", "app.dev"}

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	for _, configName := range configNames {
		viper.SetConfigName(configName)
		err = viper.ReadInConfig()
		if err == nil {
			logging.Info("Using config file: %s", viper.ConfigFileUsed())
			break
		}
	}

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	if config.StorageType == "s3" {
		err = viper.Unmarshal(&config.Minio)
	} else {
		err = viper.Unmarshal(&config.LocalStorage)
	}
	if err != nil {
		return
	}

	return
}
