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

package builder

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"

	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/color"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/schema"
	"github.com/DevopsArtFactory/redhawk/pkg/templates"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

type Builder struct {
	Config *schema.Config
	Flags  Flags
}

type Flags struct {
	Detail    bool   `json:"detail"`
	All       bool   `json:"all"`
	Config    string `json:"config"`
	Resources string `json:"resources"`
	Output    string `json:"output"`
	Region    string `json:"region"`
}

// ValidateFlags checks validation of flags
func ValidateFlags(flags Flags) error {
	if len(flags.Resources) == 0 {
		return printHelp(flags.Region)
	}

	if !tools.IsStringInArray(flags.Output, constants.ValidFormats) {
		return fmt.Errorf("output format is not supported: %s", flags.Output)
	}

	if flags.All && len(flags.Region) > 0 {
		return fmt.Errorf("--all and --region cannot be used at the same time")
	}

	if len(flags.Resources) > 0 {
		split := strings.Split(flags.Resources, ",")
		for _, resource := range split {
			if _, ok := constants.ResourceGlobal[resource]; !ok {
				return fmt.Errorf("resource is not supported: %s", resource)
			}
		}
	}

	return nil
}

// CreateNewBuilder creates new builder
func CreateNewBuilder(flags Flags) (*Builder, error) {
	var config schema.Config

	if len(flags.Config) == 0 {
		logrus.Debug("you have no config file")
		return New(nil, flags), nil
	}

	if !tools.FileExists(flags.Config) {
		return nil, fmt.Errorf("configuration file does not exist: %s", flags.Config)
	}

	file, err := ioutil.ReadFile(flags.Config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return New(&config, flags), nil
}

// Create new builder
func New(config *schema.Config, flags Flags) *Builder {
	return SetDefault(&Builder{
		Config: config,
		Flags:  flags,
	})
}

// SetDefault will set the default values for builder
func SetDefault(builder *Builder) *Builder {
	if builder.Config == nil {
		builder.Config = &schema.Config{
			Provider: constants.DefaultProvider,
		}
	}

	if builder.Flags.All {
		builder.Config.Regions = constants.AllAWSRegions
		logrus.Debugf("all regions will be applied: [%s]", strings.Join(constants.AllAWSRegions, ","))
	} else if len(builder.Flags.Region) > 0 {
		builder.Config.Regions = []string{builder.Flags.Region}
		logrus.Debugf("region will be applied: %s", builder.Flags.Region)
	}

	if len(builder.Flags.Resources) > 0 {
		var selectedResources []schema.Resource
		resources := strings.Split(builder.Flags.Resources, ",")
		for _, resource := range resources {
			selectedResources = append(selectedResources, schema.Resource{
				Name:   resource,
				Global: constants.ResourceGlobal[resource],
			})
		}

		builder.Config.Resources = selectedResources
	} else {
		if len(builder.Config.Resources) > 0 {
			for i, resource := range builder.Config.Resources {
				builder.Config.Resources[i].Global = constants.ResourceGlobal[resource.Name]
			}
		} else {
			builder.Config.Resources = defaultResources()
			logrus.Debug("all resources will be applied")
		}
	}

	return builder
}

// defaultResources returns a list of all target resources
func defaultResources() []schema.Resource {
	var ret []schema.Resource
	for _, rc := range constants.ResourceConfigs {
		if rc.Default {
			ret = append(ret, schema.Resource{
				Name:   rc.Name,
				Global: constants.ResourceGlobal[rc.Name],
			})
		}
	}
	return ret
}

// GetFlags makes flags from command
func GetFlags() (Flags, error) {
	keys := viper.AllKeys()
	flags := Flags{}

	val := reflect.ValueOf(&flags).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		key := strings.ReplaceAll(typeField.Tag.Get("json"), "_", "-")
		if tools.IsStringInArray(key, keys) {
			t := val.FieldByName(typeField.Name)
			if t.CanSet() {
				switch t.Kind() {
				case reflect.String:
					t.SetString(viper.GetString(key))
				case reflect.Int:
					t.SetInt(viper.GetInt64(key))
				case reflect.Bool:
					t.SetBool(viper.GetBool(key))
				}
			}
		}
	}

	// if region is not specified, then apply default region
	if len(flags.Region) == 0 && !flags.All {
		defaultRegion, err := getDefaultRegion("default")
		if err != nil {
			flags.Region = constants.DefaultRegion
		}

		logrus.Debugf("default region will be applied: %s", defaultRegion)
		flags.Region = defaultRegion
	}

	return flags, nil
}

// getDefaultRegion gets default region with env or configuration file
func getDefaultRegion(profile string) (string, error) {
	if len(os.Getenv(constants.DefaultRegionVariable)) > 0 {
		return os.Getenv(constants.DefaultRegionVariable), nil
	}

	functions := []func() (*ini.File, error){
		ReadAWSCredentials,
		ReadAWSConfig,
	}

	for _, f := range functions {
		cfg, err := f()
		if err != nil {
			return constants.EmptyString, err
		}

		section, err := cfg.GetSection(profile)
		if err != nil {
			return constants.EmptyString, err
		}

		if _, err := section.GetKey("region"); err == nil && len(section.Key("region").String()) > 0 {
			return section.Key("region").String(), nil
		}
	}
	return constants.EmptyString, errors.New("no aws region configuration exists")
}

// ReadAWSCredentials parse an aws credentials
func ReadAWSCredentials() (*ini.File, error) {
	if !tools.FileExists(constants.AWSCredentialsPath) {
		return ReadAWSConfig()
	}

	cfg, err := ini.Load(constants.AWSCredentialsPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ReadAWSConfig parse an aws configuration
func ReadAWSConfig() (*ini.File, error) {
	if !tools.FileExists(constants.AWSConfigPath) {
		return nil, fmt.Errorf("no aws configuration file exists in $HOME/%s", constants.AWSConfigPath)
	}

	cfg, err := ini.Load(constants.AWSConfigPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// printHelp shows the basic usage of redhawk
func printHelp(region string) error {
	sts, err := client.NewSTSClient(region)
	if err != nil {
		return err
	}

	account, err := sts.CheckWhoIam()
	if err != nil {
		return err
	}

	var data = struct {
		Account string
	}{
		Account: *account,
	}

	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
	}

	// Template for scan result
	w := tabwriter.NewWriter(os.Stdout, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("How to use").Funcs(funcMap).Parse(templates.HelperTemplates))

	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	w.Flush()

	return errors.New("please check command")
}
