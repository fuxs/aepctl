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
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
