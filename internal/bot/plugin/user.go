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
	"fmt"

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
	// Register user register command.
	p.registerUserCommand()
	// Register user command.
	p.userCommand()
}

// The command registers a new user.
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

// The command getting public information about the user.
func (p *UserPlugin) userCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "user",
			Description: "The command getting public information about the user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Member",
					Required:    false,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Setting the author.
			if i.ApplicationCommandData().Options == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.ApplicationCommandData().Options[0].UserValue(s)
			}

			// Getting a user.
			user, err := p.service.Get(context.Background(), author.ID)
			if err != nil {
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
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: author.Username,
							Description: fmt.Sprintf("**Token Balance:** %d\n", user.Balance) +
								fmt.Sprintf("**Used Promo:** %s\n", user.Used) +
								fmt.Sprintf("**Own Promo:** %s\n", user.Promo),
							Color: 0xa735ed,
						},
					},
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}
