package get

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetOptions struct {
	Settings *settings.Settings

	// Args
	ParameterName string
}

func NewCmdGet(settings *settings.Settings) *cobra.Command {
	opts := &GetOptions{
		Settings: settings,
	}

	var getCmd = &cobra.Command{
		Use:     "get",
		Version: settings.Version,
		Short:   "Get any field from the buddy config file",
		Long:    `Get any field from the buddy config file`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ParameterName = args[0]
			opts.Settings.Logger = opts.Settings.Logger.With("parameter", opts.ParameterName)

			return RunGet(opts)
		},
	}

	return getCmd
}

func RunGet(opts *GetOptions) error {
	opts.Settings.Logger.Debug("Getting field from buddy config")

	projectConfig := opts.Settings.ProjectConfig

	caser := cases.Title(language.English)

	configKey := caser.String(strings.ToLower(opts.ParameterName))
	opts.Settings.Logger.Debug("Titlefy field name", "field", configKey)

	r := reflect.ValueOf(projectConfig)
	value := reflect.Indirect(r).FieldByName(configKey)

	if !value.IsValid() {
		opts.Settings.Logger.Error("Field not found", "field", configKey)
		return fmt.Errorf("Field %s not found", configKey)
	}

	if value.Kind() != reflect.String {
		opts.Settings.Logger.Error("Field is not a printable field", "field", configKey, "kind", value.Kind().String())
		return fmt.Errorf("Field %s is not a printable field", configKey)
	}

	fmt.Println(value.String())

	return nil
}
