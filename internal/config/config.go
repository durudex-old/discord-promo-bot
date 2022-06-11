/*
 * Copyright © 2022 Durudex

 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.

 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

type (
	// Config variables.
	Config struct{ Bot BotConfig }

	// Discord bot config variables.
	BotConfig struct{ Token string }
)

// Initialize config.
func Init() (*Config, error) {
	log.Debug().Msg("Initialize config...")

	var cfg Config

	// Set env configurations.
	setFromEnv(&cfg)

	return &cfg, nil
}

// Setting environment variables from .env file.
func setFromEnv(cfg *Config) {
	log.Debug().Msg("Set from environment configurations...")

	// Discord bot variables.
	cfg.Bot.Token = os.Getenv("BOT_TOKEN")
}
