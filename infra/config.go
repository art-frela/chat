package infra

import (
	"fmt"
	"os"
	"path"
	"strings"

	"flag"

	"github.com/spf13/viper"
)

const (
	configFormat = "yaml"
	fPort        = "p"
	fHost        = "host"
	fLevel       = "level"
	fPathConfig  = "c"
	envPrefix    = "CHAT"
)

var (
	configName = "config"
	configPath = "."
)

func (s *Server) setConfig() {
	flag.String(fPathConfig, path.Join(configPath, configName+"."+configFormat), "path to config file for application")
	flag.Int(fPort, 8080, "http port for application")
	flag.String(fHost, "0.0.0.0", "ip address for bind application")
	flag.String(fLevel, "debug", "loglevel fog logging information")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case fPort:
			os.Setenv(envPrefix+"_HTTPD_PORT", f.Value.String())
		case fHost:
			os.Setenv(envPrefix+"_HTTPD_HOST", f.Value.String())
		case fLevel:
			os.Setenv(envPrefix+"_LOG_LEVEL", f.Value.String())
		case fPathConfig:
			configPath, configName = path.Split(f.Value.String())
			configName = strings.Replace(configName, path.Ext(configName), "", -1)
		}

	})

	config := viper.New()
	config.SetConfigType(configFormat)
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetEnvPrefix(envPrefix)
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	config.AutomaticEnv()
	s.config = config
}
