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
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"
	"github.com/rs/zerolog/log"
)

// User commands plugin structure.
type UserPlugin struct {
	service service.User
	handler *command.Handler
}

// Creating a new user commands plugin.
func NewUserPlugin(service service.User, handler *command.Handler) *UserPlugin {
	return &UserPlugin{service: service, handler: handler}
}

// Registering user plugin commands.
func (p *UserPlugin) RegisterCommands() {
	// Register user command.
	p.registerUserCommand()
}

// Registering a new user command.
func (p *UserPlugin) registerUserCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "register",
			Description: "The command registers a new user.",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Creating a new user.
			if err := p.service.Create(context.Background(), domain.User{
				DiscordId: i.Interaction.Member.User.ID,
			}); err != nil {
				// Send a interaction respond error message.
				if err := discordInteractionError(s, i, err); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond error message")
				}

				return
			}

			// Send a interaction respond message.
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You have successfully registered!",
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}
