package builder

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/schema"
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
	} else {
		if len(builder.Flags.Region) > 0 {
			builder.Config.Regions = []string{builder.Flags.Region}
			logrus.Debugf("region will be applied: %s", builder.Flags.Region)
		} else if len(builder.Config.Regions) == 0 {
			builder.Config.Regions = []string{constants.DefaultRegion}
			logrus.Debugf("default region will be applied: %s", constants.DefaultRegion)
		}
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

	return flags, nil
}
