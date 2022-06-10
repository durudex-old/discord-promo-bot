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

package plugin

import (
	"github.com/durudex/discord-promo-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Bot commands plugin structure.
type BotPlugin struct{}

// Creating a new bot commands plugin.
func NewBotPlugin() *BotPlugin {
	return &BotPlugin{}
}

// Registering bot plugin commands.
func (p *BotPlugin) RegisterCommands(handler *command.Handler) error {
	// Register github command.
	err := handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "github",
			Description: "The command sends a link to the bot's source code.",
		},
		Handler: p.github,
	})

	return err
}

// Github command handler.
func (p *BotPlugin) github(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send a interaction respond message.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "https://github.com/durudex/discord-promo-bot",
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}
