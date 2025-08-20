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
	Driver string `mapstructure:"driver"`
	URL    string `mapstructure:"url"`
}

type Web struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type App struct {
	PasswordMinLength int `mapstructure:"password_min_length"`
}

type Config struct {
	Database       Database `mapstructure:"database"`
	Web            Web      `mapstructure:"web"`
	App            App      `mapstructure:"app"`
	Logfile        string   `mapstructure:"logfile"`
	SecretKey      string   `mapstructure:"secretkey"`
	TrustedOrigins []string `mapstructure:"trustedorigins"`
	Debug          bool     `mapstructure:"debug"`
}

var Cfg = &Config{
	Database: Database{
		Driver: "sqlite",
		URL:    "someblocks.sqlite",
	},
	Web: Web{
		Host: "127.0.0.1",
		Port: 8080,
	},
	App: App{
		PasswordMinLength: 8,
	},
	Logfile:        "someblocks.log",
	SecretKey:      makeSecretKey(32),
	TrustedOrigins: []string{"localhost:8000", "127.0.0.1:8000"},
	Debug:          false,
}

func loadDefaultConfig() error {
	cfgMap := make(map[string]interface{})
	err := mapstructure.Decode(Cfg, &cfgMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode default config.")
		return err
	}

	cfgYamlBytes, err := yaml.Marshal(&cfgMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal default config.")
		return err
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(cfgYamlBytes))
	if err != nil {
		log.Error().Err(err).Msg("Failed to read default config.")
		return err
	}
	return nil
}

func Load(cfgFile string) error {
	err := loadDefaultConfig()
	if err != nil {
		log.Error().Msgf("Error loading default config.")
		return err
	}

	// Setup viper to load our config from the filesystem and ENV
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for the default config in the src config directory

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	err = viper.MergeInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Error().Msgf("Error using config file: %s", err)
		return err
	}
	viper.AutomaticEnv()

	if err := viper.MergeInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(Cfg); err != nil {
		log.Error().Err(err).Msgf("Couldn't unmarshal viper config into Cfg")
		return err
	}
	return err
}

func ToYAML() string {
	bs, err := yaml.Marshal(Cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

func makeSecretKey(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't generate a secret key.")
	}
	//return string(b[:])
	return base64.URLEncoding.EncodeToString(b)
}
