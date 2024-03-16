package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dreadster3/buddy/pkg/cmd/root"
	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/log"
	"github.com/spf13/viper"
)

func main() {
	defaultAuthor := "Anonymous"
	currentUser, _ := user.Current()
	if currentUser != nil {
		defaultAuthor = currentUser.Username
	}

	// Initialize configurations
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$XDG_CONFIG_HOME/buddy")
	viper.AddConfigPath("$HOME/.config/buddy")
	viper.SetDefault("author", defaultAuthor)
	viper.SetDefault("filename", "buddy.json")
	viper.SetDefault("scripts", map[string]string{})
	viper.SetDefault("templates_path", "templates")
	viper.ReadInConfig()

	if viper.ConfigFileUsed() == "" {
		configDir := os.ExpandEnv("$HOME/.config/buddy")

		log.Logger.Info("No config file found, creating one", "path", configDir)

		err := os.Mkdir(configDir, 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		err = viper.SafeWriteConfig()
		if err != nil {
			fmt.Println("Error writing config file:", err)
			os.Exit(1)
		}
	}

	globalConfig, err := config.GlobalConfigFromViper(viper.GetViper())
	if err != nil {
		log.Logger.Error("Error parsing config file", "error", err)
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	settings := settings.New(Version, globalConfig)
	cmd := root.NewRootCmd(settings)

	if err := cmd.Execute(); err != nil {
		log.Logger.Error("Error executing command", "error", err)
		os.Exit(1)
	}
}
