package get

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetOptions struct {
	GlobalConfig *config.GlobalConfig

	// Args
	ParameterName string
}

func NewCmdGet(globalConfig *config.GlobalConfig) *cobra.Command {
	opts := &GetOptions{
		GlobalConfig: globalConfig,
	}

	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get any field from the buddy config file",
		Long:  `Get any field from the buddy config file`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ParameterName = args[0]

			return runGet(opts)
		},
	}

	return getCmd
}

func runGet(opts *GetOptions) error {
	buddyConfig, err := config.ParseMergeProjectConfigFile(opts.GlobalConfig)
	if err != nil {
		return err
	}

	caser := cases.Title(language.English)

	configKey := caser.String(strings.ToLower(opts.ParameterName))

	r := reflect.ValueOf(buddyConfig)
	value := reflect.Indirect(r).FieldByName(configKey)

	if !value.IsValid() {
		return fmt.Errorf("Field %s not found", configKey)
	}

	if value.Kind() != reflect.String {
		return fmt.Errorf("Field %s is not a printable field", configKey)
	}

	fmt.Println(value.String())

	return nil
}
