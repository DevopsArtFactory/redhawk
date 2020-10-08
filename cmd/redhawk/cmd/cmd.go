package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
	"github.com/DevopsArtFactory/redhawk/pkg/version"
)

var (
	cfgFile string
	v       string
)

// Get root command
func NewRootCommand(out, stderr io.Writer) *cobra.Command {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:           "redhawk",
		Short:         "Opensource client resource audit and management tool",
		Long:          "Opensource client resource audit and management tool",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().SetOutput(out)

			// Setup logs
			if err := tools.SetUpLogs(stderr, v); err != nil {
				return err
			}

			version := version.Get()

			logrus.Infof("redhawk %+v", version)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	//Group by commands
	groups := templates.CommandGroups{
		{
			Message: "managing configuration of redhawk",
			Commands: []*cobra.Command{
				NewCmdScanResources(),
			},
		},
	}

	groups.Add(rootCmd)

	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(NewCmdVersion())

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", constants.DefaultLogLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	templates.ActsAsRootCommand(rootCmd, nil, groups...)

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
