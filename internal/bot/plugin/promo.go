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

package plugin

import (
	"context"
	"fmt"

	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Promo commands plugin structure.
type PromoPlugin struct {
	service service.User
	handler *command.Handler
	bot     *config.BotConfig
}

// Creating a new promo commands plugin.
func NewPromoPlugin(service service.User, handler *command.Handler, cfg *config.Config) *PromoPlugin {
	return &PromoPlugin{service: service, handler: handler, bot: &cfg.Bot}
}

// Registering promo plugin commands.
func (p *PromoPlugin) RegisterCommands() {
	// Register create promo command.
	p.createPromoCommand()
	// Register use promo command.
	p.usePromoCommand()
}

// The command creating a new user promo code.
func (p *PromoPlugin) createPromoCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:         "create",
			Description:  "The command creating a new user promo code.",
			DMPermission: &DMPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "promo",
					Description: "Promo code.",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Checking where the command was use.
			if i.Interaction.User == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.Interaction.User
			}

			// Updating a user.
			if err := p.service.Update(
				context.Background(),
				domain.User{
					Id:    author.ID,
					Promo: i.ApplicationCommandData().Options[0].StringValue(),
				},
			); err != nil {
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
					Content: fmt.Sprintf("You created promo code `%s`", i.ApplicationCommandData().Options[0].StringValue()),
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}

			// Send bot log message.
			if _, err := s.ChannelMessageSendEmbed(
				p.bot.LogChannel,
				&discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						URL:     "https://discord.com/users/" + author.ID,
						Name:    author.Username,
						IconURL: author.AvatarURL("128x128"),
					},
					Description: fmt.Sprintf(
						"User created a new promo code `%s`.",
						i.ApplicationCommandData().Options[0].StringValue(),
					),
					Color: p.bot.Color,
				},
			); err != nil {
				log.Warn().Err(err).Msg("failed to send channel message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// The command use a user promo code.
func (p *PromoPlugin) usePromoCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:         "use",
			Description:  "The command use a user promo code.",
			DMPermission: &DMPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "promo",
					Description: "Promo code.",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Checking where the command was use.
			if i.Interaction.User == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.Interaction.User
			}

			// Use a promo code.
			reward, err := p.service.UsePromo(
				context.Background(),
				author.ID,
				i.ApplicationCommandData().Options[0].StringValue(),
			)
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
					Content: fmt.Sprintf("You used promo code `%s`", i.ApplicationCommandData().Options[0].StringValue()),
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}

			// Send bot log message.
			if _, err := s.ChannelMessageSendEmbed(
				p.bot.LogChannel,
				&discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						URL:     "https://discord.com/users/" + author.ID,
						Name:    author.Username,
						IconURL: author.AvatarURL("128x128"),
					},
					Description: fmt.Sprintf(
						"User used the promo code `%s` and received %d DUR.",
						i.ApplicationCommandData().Options[0].StringValue(),
						reward,
					),
					Color: p.bot.Color,
				},
			); err != nil {
				log.Warn().Err(err).Msg("failed to send channel message")
			}
		},
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}
