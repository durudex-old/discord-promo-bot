/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Default config path.
const defaultConfigPath string = "configs/main"

type (
	// Config variables.
	Config struct {
		Bot      BotConfig
		Database DatabaseConfig
		User     UserConfig
		Promo    PromoConfig
	}

	// Discord bot config variables.
	BotConfig struct {
		Color      int    `mapstructure:"color"`
		LogChannel string `mapstructure:"log-channel"`
		Token      string
	}

	// Database config variables.
	DatabaseConfig struct {
		Mongodb MongodbConfig `mapstructure:"mongodb"`
	}

	// Mongodb config variables.
	MongodbConfig struct {
		URI      string
		Username string
		Password string
		Database string        `mapstructure:"database"`
		Timeout  time.Duration `mapstructure:"timeout"`
	}

	// User config variables.
	UserConfig struct {
		ReviewRole string        `mapstructure:"review-role"`
		MinAge     time.Duration `mapstructure:"min-age"`
	}

	// Promo config variables.
	PromoConfig struct {
		Award int `mapstructure:"award"`
	}
)

// Initialize config.
func Init() (*Config, error) {
	log.Debug().Msg("Initialize config...")

	// Parsing specified when starting the config file.
	if err := parseConfigFile(); err != nil {
		return nil, err
	}

	var cfg Config

	// Unmarshal config keys.
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set env configurations.
	setFromEnv(&cfg)

	return &cfg, nil
}

// Parsing specified when starting the config file.
func parseConfigFile() error {
	// Get config path variable.
	configPath := os.Getenv("CONFIG_PATH")

	// Check is config path variable empty.
	if configPath == "" {
		configPath = defaultConfigPath
	}

	log.Debug().Msgf("Parsing config file: %s", configPath)

	// Split path to folder and file.
	dir, file := filepath.Split(configPath)

	viper.AddConfigPath(dir)
	viper.SetConfigName(file)

	// Read config file.
	return viper.ReadInConfig()
}

// Unmarshal config keys.
func unmarshal(cfg *Config) error {
	log.Debug().Msg("Unmarshal config keys...")

	// Unmarshal bot keys.
	if err := viper.UnmarshalKey("bot", &cfg.Bot); err != nil {
		return err
	}
	// Unmarshal user keys.
	if err := viper.UnmarshalKey("user", &cfg.User); err != nil {
		return err
	}
	// Unmarshal promo keys.
	if err := viper.UnmarshalKey("promo", &cfg.Promo); err != nil {
		return err
	}
	// Unmarshal database keys.
	return viper.UnmarshalKey("database", &cfg.Database)
}

// Setting environment variables from .env file.
func setFromEnv(cfg *Config) {
	log.Debug().Msg("Set from environment configurations...")

	// Discord bot variables.
	cfg.Bot.Token = os.Getenv("BOT_TOKEN")

	// Mongo database variables.
	cfg.Database.Mongodb.URI = os.Getenv("MONGO_URI")
	cfg.Database.Mongodb.Username = os.Getenv("MONGO_USERNAME")
	cfg.Database.Mongodb.Password = os.Getenv("MONGO_PASSWORD")
}
