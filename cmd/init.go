package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

func init() {
	buddyCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:                   "init [options] [directory]",
	DisableFlagsInUseLine: true,
	Short:                 "Initialize a new buddy file",
	Long: `Initialize a new buddy file.
	If a directory is provided, buddy.json will be created in the directory.
	If no directory is provided, buddy.json will be created in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			err := os.Mkdir(args[0], 0755)
			if err != nil {
				return
			}

			os.Chdir(args[0])
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("buddy.json"); err == nil {
			return errors.New("buddy.json already exists")
		}

		workingDir, err := os.Getwd()
		if err != nil {
			return err
		}

		projectName := path.Base(workingDir)

		if len(args) > 0 {
			projectName = args[0]
		}

		buddyConfig := models.NewBuddyConfig(projectName, "0.0.1", "A new buddy project", "Anonymous", map[string]string{})

		file, err := os.Create("buddy.json")
		if err != nil {
			return err
		}
		defer file.Close()

		json, err := buddyConfig.ToJson()
		if err != nil {
			return err
		}

		_, err = file.WriteString(string(json))
		if err != nil {
			return err
		}

		fmt.Println("buddy.json created")

		return nil
	},
}
