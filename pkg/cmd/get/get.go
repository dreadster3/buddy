package get

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get any field from the buddy config file",
		Long:  `Get any field from the buddy config file`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[0]

			buddyConfig, err := config.ParseProjectConfigFile("buddy.json")
			if err != nil {
				return err
			}

			// TODO: Use the non-deprecated version of this function
			configKey = strings.Title(strings.ToLower(configKey))

			r := reflect.ValueOf(buddyConfig)
			value := reflect.Indirect(r).FieldByName(configKey)

			if !value.IsValid() {
				return fmt.Errorf("Field %s not found", configKey)
			}

			fmt.Print(value.String())

			return nil
		},
	}

	return getCmd
}
