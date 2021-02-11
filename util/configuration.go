/*
Package util util consists of general utility functions and structures.

Copyright 2021 Michael Bungenstock

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// GlobalConfig contains the global configuration
type GlobalConfig struct {
	Name   string
	Config string
	Debug  bool
	Human  bool
}

// NewGlobalConfig returns an initilized configuration.
func NewGlobalConfig(name string, cmd *cobra.Command) *GlobalConfig {
	o := &GlobalConfig{}
	o.Name = name
	flags := cmd.PersistentFlags()

	flags.StringVar(&o.Config, "config", "", "path to configuration file")
	flags.BoolVar(&o.Debug, "debug", false, "sets log level to debug")
	flags.BoolVar(&o.Human, "human", false, "human readable logging to console")
	return o
}

// GetPreRunE returns the function for PreRunE function of cobra.Comamand
func (o *GlobalConfig) GetPreRunE() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return o.Configure(cmd)
	}
}

// Configure loads the configuration file encoded in json or yaml.
func (o *GlobalConfig) Configure(cmd *cobra.Command) error {
	// set up logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if o.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if o.Human {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if o.Config != "" {
		viper.SetConfigFile(o.Config)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/" + o.Name + "/")
		if home, err := os.UserHomeDir(); err != nil {
			viper.AddConfigPath(path.Join(home, "."+o.Name+"/"))
		}
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error
			} else {
				// Config file was found but another error was produced
				return err
			}
		} else {
			log.Debug().Str("Service", o.Name).Str("Config file", viper.ConfigFileUsed()).Msg("Successfully loaded configuration file")
		}
	}

	for act := cmd; act != nil; act = act.Parent() {
		act.Flags().VisitAll(func(f *pflag.Flag) {
			if !f.Changed && viper.IsSet(f.Name) {
				_ = f.Value.Set(viper.GetString(f.Name))
			}
		})
	}
	return nil
}

// ConfigFile represents a configuration file in YAML format
type ConfigFile struct {
	Node *yaml.Node
	Path string
}

// Query returns this configuraion as YAMLQuery
func (f *ConfigFile) Query() *YAMLQuery {
	return NewYAMLQuery(f.Node)
}

// Organization returns the current organization value
func (f *ConfigFile) Organization() string {
	return f.Query().Str("organization")
}

// SetOrganization sets the new organization value
func (f *ConfigFile) SetOrganization(org string) {
	f.Query().SetMap("organization", org)
}

// TechAccount returns the current technical account value
func (f *ConfigFile) TechAccount() string {
	return f.Query().Str("tech-account")
}

// SetTechAccount sets the new technical account value
func (f *ConfigFile) SetTechAccount(account string) {
	f.Query().SetMap("tech-account", account)
}

// ClientID returns the current client id value
func (f *ConfigFile) ClientID() string {
	return f.Query().Str("client-id")
}

// SetClientID sets the new client id value
func (f *ConfigFile) SetClientID(id string) {
	f.Query().SetMap("client-id", id)
}

// ClientSecret returns the current client secret value
func (f *ConfigFile) ClientSecret() string {
	return f.Query().Str("client-secret")
}

// SetClientSecret sets the new client secret value
func (f *ConfigFile) SetClientSecret(secret string) {
	f.Query().SetMap("client-secret", secret)
}

// Key returns the current private key path value
func (f *ConfigFile) Key() string {
	return f.Query().Str("key")
}

// SetKey sets the new private key path value
func (f *ConfigFile) SetKey(key string) {
	f.Query().SetMap("key", key)
}

// Sandbox returns the current sandbox value
func (f *ConfigFile) Sandbox() string {
	return f.Query().Str("sandbox")
}

// SetSandbox sets the new sandbox value
func (f *ConfigFile) SetSandbox(sandbox string) {
	f.Query().SetMap("sandbox", sandbox)
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	return path.Join(home, ".aepctl", "config.yaml"), err
}

// LoadConfigFile loads the configuration file in YAML format for the passed
// path
func LoadConfigFile(path string) (*ConfigFile, error) {
	if path == "" {
		cp, err := configPath()
		if err != nil {
			return nil, err
		}
		path = cp
	}
	_, err := os.Stat(path)
	if err == nil {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		node := &yaml.Node{}
		err = yaml.Unmarshal(data, node)
		if err != nil {
			return nil, err
		}
		return &ConfigFile{
			Path: path,
			Node: node,
		}, nil
	}
	if os.IsNotExist(err) {
		// return an empty config
		return &ConfigFile{
			Path: path,
			Node: &yaml.Node{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:    yaml.MappingNode,
						Content: make([]*yaml.Node, 0, 16),
					},
				},
				HeadComment: "# Generated by aepctl configure",
			},
		}, nil
	}
	return nil, err
}

// Save saves the configuration in YAML format
func (f *ConfigFile) Save() error {
	if _, err := os.Stat(f.Path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(path.Dir(f.Path), 0700); err != nil {
			return err
		}
	}
	data, err := yaml.Marshal(f.Node)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.Path, data, 0600)
}
