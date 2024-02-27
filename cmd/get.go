package cmd

import (
	"reflect"
	"strings"

	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

func init() {
	buddyCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get any field from the buddy config file",
	Long:  `Get any field from the buddy config file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configKey := args[0]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		configKey = strings.Title(strings.ToLower(configKey))

		r := reflect.ValueOf(buddyConfig)
		value := reflect.Indirect(r).FieldByName(configKey)

		if !value.IsValid() {
			cmd.PrintErrf("Field %s not found", configKey)
			return
		}

		cmd.Println(value.String())
	},
}
