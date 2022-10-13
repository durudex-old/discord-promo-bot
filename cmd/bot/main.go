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

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/durudex/discord-promo-bot/internal/bot/command"
	"github.com/durudex/discord-promo-bot/internal/bot/event"
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/repository"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/bot"
	"github.com/durudex/discord-promo-bot/pkg/database/mongodb"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initialize application.
func init() {
	// Set logger mode.
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// A function that running the bot.
func main() {
	// Initialize config.
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize config.")
	}

	// Creating a new discord bot.
	b, err := bot.New(cfg.Bot.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a discord session")
	}

	// Running the discord bot.
	if err := b.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to running discord bot")
	}

	// Initializing the discord event handlers.
	event.NewEvent(b).InitEvents()

	// Creating a new mongodb client.
	client, err := mongodb.NewClient(&mongodb.MongoConfig{
		URI:      cfg.Database.Mongodb.URI,
		Username: cfg.Database.Mongodb.Username,
		Password: cfg.Database.Mongodb.Password,
		Timeout:  cfg.Database.Mongodb.Timeout,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create mongodb client")
	}

	// Creating a new repository.
	repos := repository.NewRepository(client.Database(cfg.Database.Mongodb.Database))
	// Creating a new service.
	service := service.NewService(repos)

	// Starting promo monitoring.
	startMonitor(service.Monitor, cfg.Promo.AutoSaveTTL)

	// Registering all discord commands.
	command.NewCommandPlugin(b, cfg, service).Register()

	// Quit in application.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Closing a bot connections.
	if err := b.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close discord connection")
	}

	// Saving promo monitor.
	if err := service.Monitor.Save(context.Background()); err != nil {
		log.Error().Err(err).Msg("error saving monitor")
	}

	log.Info().Msg("Discord Promo Bot stopping!")
}

// Starting promo monitoring.
func startMonitor(mon service.Monitor, ttl time.Duration) {
	// Sync promo monitor with database.
	if err := mon.Sync(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("error sync monitor")
	}

	go func() {
		for {
			time.Sleep(ttl)

			// Saving promo monitor.
			if err := mon.Save(context.Background()); err != nil {
				log.Error().Err(err).Msg("error saving monitor")
			}
		}
	}()
}
