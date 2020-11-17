package config

import (
	"bytes"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

/*
map[string]interface{} =
map[string]interface{}{
    "port":     9090,
    "hostname": "localhost",
    "auth": map[string]string{
      "username": "titpetric",
      "password": "12fa",
    },
})
*/

func LoadConfig() {
	loadConfigFromStruct(Cfg)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

type Database struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Web struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Web      Web      `mapstructure:"web"`
	Logfile  string   `mapstructure:"logfile"`
}

var Cfg = &Config{
	Database: Database{
		Driver:   "sqlite3",
		Dbname:   "someblocks.sqlite",
		Username: "",
		Password: "",
	},
	Web: Web{
		Host: "127.0.0.1",
		Port: 8080,
	},
	Logfile: "someblocks.log",
}

func loadConfigFromStruct(cfg interface{}) {
	cfgMap := make(map[string]interface{})
	err := mapstructure.Decode(cfg, &cfgMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal default config")
	}

	cfgJsonBytes, err := json.Marshal(&cfgMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal default config")
	}

	viper.SetConfigType("json")
	err = viper.ReadConfig(bytes.NewReader(cfgJsonBytes))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read default config")
	}
}

func Load(cfgFile string) {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")               // look for the default config in the src config directory
	viper.AddConfigPath("/etc/someblocks") // look for config in the src config directory

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal().Msgf("Error using config file: %s", err)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatal().Msgf("Couldn't unmarshal viper config into Cfg", err)
	}
}
