/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
		Short:         "Opensource cloud resources audit and management tool",
		Long:          "Opensource cloud resources audit and management tool",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().SetOutput(out)

			// Setup logs
			if err := tools.SetUpLogs(stderr, v); err != nil {
				return err
			}

			version := version.Get()

			logrus.Debugf("redhawk %+v", version)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	//Group by commands
	groups := templates.CommandGroups{
		{
			Message: "checking all resources in cloud provider",
			Commands: []*cobra.Command{
				NewCmdList(),
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
