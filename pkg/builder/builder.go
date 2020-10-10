package builder

import (
	"fmt"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"

	"github.com/DevopsArtFactory/redhawk/pkg/schema"
)

type Builder struct {
	Config *schema.Config
}

// CreateNewBuilder creates new builder
func CreateNewBuilder() (*Builder, error) {
	var config schema.Config
	configPath := viper.GetString("config")

	if len(configPath) == 0 {
		logrus.Debug("you have no config file")
		return nil, nil
	}

	if !tools.FileExists(configPath) {
		return nil, fmt.Errorf("configuration file does not exist: %s", configPath)
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return New(&config), nil
}

// Create new builder
func New(config *schema.Config) *Builder {
	return &Builder{Config: config}
}

// SetDefault will set the default values for builder
func SetDefault(builder *Builder) *Builder {
	if builder == nil {
		builder = &Builder{
			Config: &schema.Config{
				Provider: constants.DefaultProvider,
			},
		}
	}

	if len(builder.Config.Regions) == 0 {
		builder.Config.Regions = constants.AllAWSRegions
		logrus.Debugf("all regions will be applied: [%s]", strings.Join(constants.AllAWSRegions, ","))
	}

	if len(builder.Config.Resources) == 0 {
		builder.Config.Resources = defaultResources()
		logrus.Debug("all resources will be applied")
	}
	return builder
}

// defaultResources returns a list of all target resources
func defaultResources() []schema.Resource {
	return []schema.Resource{
		{
			Name: constants.EC2ResourceName,
		},
		{
			Name: constants.SGResourceName,
		},
		{
			Name: constants.Route53ResourceName,
		},
	}
}
