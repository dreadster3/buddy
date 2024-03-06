package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dreadster3/buddy/pkg/cmd/root"
	"github.com/dreadster3/buddy/pkg/config"
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
	viper.AddConfigPath("$HOME/.config/buddy")
	viper.AddConfigPath("$XDG_CONFIG_HOME/buddy")
	viper.SetDefault("author", defaultAuthor)
	viper.SetDefault("filename", "buddy.json")
	viper.SetDefault("scripts", map[string]string{})
	viper.ReadInConfig()

	if viper.ConfigFileUsed() == "" {
		fmt.Println("No config file found, creating one at $HOME/.config/buddy")

		err := os.Mkdir(os.ExpandEnv("$HOME/.config/buddy"), 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		err = viper.SafeWriteConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
				fmt.Println("Error writing config file:", err)
				os.Exit(1)
			}
		}
	}

	globalConfig, err := config.GlobalConfigFromViper(viper.GetViper())
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	cmd := root.NewRootCmd(globalConfig)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
