package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const root = "/tmp/wallpaper/"
const favorite = "/home/alex/Pictures/Wallpapers/"

var Conf Configuration

type Configuration struct {
	Root     string
	Favorite string
	ApiKey   string
}

func Config() error {
	// NOTE: Right now we are setting env variable manually
	pwd := os.Getenv("PWD")
	p := fmt.Sprint(pwd, "/config")
	os.Setenv("WSWITCH_CONFIG", p)

	path := os.Getenv("WSWITCH_CONFIG")

	viper.SetConfigName(".wswitch")
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Config: %w", err)
	}

	viper.SetDefault("ApiKey", nil)
	viper.SetDefault("Root", root)
	viper.SetDefault("Favorite", favorite)

	viper.Unmarshal(&Conf)

	return nil
}
