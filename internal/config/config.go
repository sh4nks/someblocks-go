package config

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
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
	Database  Database `mapstructure:"database"`
	Web       Web      `mapstructure:"web"`
	Logfile   string   `mapstructure:"logfile"`
	SecretKey string   `mapstructure:"secretkey"`
	Debug     bool     `mapstructure:"debug"`
}

var Cfg = &Config{
	Database: Database{
		Driver:   "sqlite3",
		Dbname:   "someblocks.sqlite",
		Host:     "",
		Username: "",
		Password: "",
	},
	Web: Web{
		Host: "127.0.0.1",
		Port: 8080,
	},
	Logfile:   "someblocks.log",
	SecretKey: makeSecretKey(32),
	Debug:     false,
}

func loadDefaultConfig(cfg interface{}) {
	cfgMap := make(map[string]interface{})
	err := mapstructure.Decode(cfg, &cfgMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal default config")
	}

	cfgYamlBytes, err := yaml.Marshal(&cfgMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal default config")
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(cfgYamlBytes))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read default config")
	}
}

func makeSecretKey(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't generate secret key.")
	}
	//return string(b[:])
	return base64.URLEncoding.EncodeToString(b)
}

func Load(cfgFile string) {
	loadDefaultConfig(Cfg)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for the default config in the src config directory

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	err := viper.MergeInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		log.Fatal().Msgf("Error using config file: %s", err)
	}
	viper.AutomaticEnv()

	if err := viper.MergeInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatal().Msgf("Couldn't unmarshal viper config into Cfg", err)
	}
}

func ToYAML() string {
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}
