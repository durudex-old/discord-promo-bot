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
	"os"
	"os/signal"
	"syscall"

	"github.com/durudex/discord-promo-bot/internal/bot/event"
	"github.com/durudex/discord-promo-bot/internal/bot/plugin"
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/repository"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"
	"github.com/durudex/discord-promo-bot/pkg/database/mongodb"

	"github.com/bwmarrin/discordgo"
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

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create discord session")
	}

	// Open a websocket connection to Discord and begin listening.
	if err := session.Open(); err != nil {
		log.Fatal().Err(err).Msg("failed to open discord websocket connection")
	}

	// Creating a new discord command handler.
	commandHandler := command.NewHandler(session)

	// Initializing the discord event handlers.
	event.NewEvent(commandHandler).InitEvents(session)

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
	service := service.NewService(repos, cfg)

	// Registering all discord commands.
	plugin.NewPlugin(service, cfg).RegisterPlugins(commandHandler)

	// Quit in application.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Close the websocket connection to Discord.
	if err := session.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close discord websocket connection")
	}

	log.Info().Msg("Discord Promo Bot stopping!")
}
