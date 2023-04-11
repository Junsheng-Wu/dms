package config

// Package config saves configuration for running EMLA components.
//
// Config can be configured from command line flags and configuration file.
// Command line flags hold higher priority than configuration file. But if
// component Endpoint/Host/APIServer was left empty, all of that component
// command line flags will be ignored, use configuration file instead.

import (
	"dms/client/rule"
	"reflect"
	"strings"
)

type Config struct {
	RuleOptions *rule.Options `json:"rule" yaml:"rule" mapstructure:"rule"`
}

// New creates a default non-empty Config
func New() *Config {
	return &Config{
		RuleOptions: rule.NewOptions(),
	}
}

//
//// TryLoadFromDisk loads configuration from default location after server startup
//// return nil error if configuration file not exists
//func TryLoadFromDisk() (*Config, error) {
//	viper.SetConfigName(defaultConfigurationName)
//	viper.AddConfigPath(defaultConfigurationPath)
//
//	// Load from current working directory, only used for debugging
//	viper.AddConfigPath(".")
//
//	if err := viper.ReadInConfig(); err != nil {
//		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
//			return nil, err
//		} else {
//			return nil, fmt.Errorf("error parsing configuration file %s", err)
//		}
//	}
//
//	conf := New()
//
//	if err := viper.Unmarshal(conf); err != nil {
//		return nil, err
//	}
//
//	return conf, nil
//}

// ToMap simply converts config to map[string]bool
// to hide sensitive information
func (conf *Config) ToMap() map[string]bool {
	conf.stripEmptyOptions()
	result := make(map[string]bool)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}

// Remove invalid options before serializing to json or yaml
func (conf *Config) stripEmptyOptions() {
	if conf.RuleOptions != nil && conf.RuleOptions.Path == "" {
		conf.RuleOptions = nil
	}

}
