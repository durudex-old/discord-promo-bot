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
	"time"

	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// User commands plugin structure.
type UserPlugin struct {
	service service.User
	handler *command.Handler
	user    *config.UserConfig
	bot     *config.BotConfig
}

// Creating a new user commands plugin.
func NewUserPlugin(service service.User, handler *command.Handler, cfg *config.Config) *UserPlugin {
	return &UserPlugin{service: service, handler: handler, user: &cfg.User, bot: &cfg.Bot}
}

// Registering user plugin commands.
func (p *UserPlugin) RegisterCommands() {
	// Register user register command.
	p.registerUserCommand()
	// Register user command.
	p.userCommand()
	// Register user update balance command.
	p.updateBalanceCommand()
}

// The command registers a new user.
func (p *UserPlugin) registerUserCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:         "register",
			Description:  "The command registers a new user.",
			DMPermission: &DMPermission,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Checking where the command was use.
			if i.Interaction.User == nil {
				author = i.Interaction.Member.User
			} else {
				author = i.Interaction.User
			}

			// Getting creating user timestamp.
			createdAt, err := discordgo.SnowflakeTimestamp(author.ID)
			if err != nil {
				return
			}

			// Checking min user account age.
			if createdAt.Add(p.user.MinAge).Unix() > time.Now().Unix() {
				// Send a interaction respond error message.
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Your account is well new!",
					},
				}); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond error message")
				}

				return
			}

			// Creating a new user.
			if err := p.service.Create(context.Background(), domain.User{DiscordId: author.ID}); err != nil {
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
			Name:         "user",
			Description:  "The command getting public information about the user.",
			DMPermission: &DMPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User",
					Required:    false,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var author *discordgo.User

			// Setting the author.
			if i.ApplicationCommandData().Options == nil {
				// Checking where the command was use.
				if i.Interaction.User == nil {
					author = i.Interaction.Member.User
				} else {
					author = i.Interaction.User
				}
			} else {
				author = i.ApplicationCommandData().Options[0].UserValue(s)
			}

			// Getting a user.
			user, err := p.service.Get(context.Background(), author.ID)
			if err != nil {
				// Send a interaction respond error message.
				if err := discordInteractionError(s, i, err); err != nil {
					log.Warn().Err(err).Msg("failed to getting user")
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
							Color: p.bot.Color,
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

// The command updating the user balance.
func (p *UserPlugin) updateBalanceCommand() {
	if err := p.handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "update-balance",
			Description: "The command updating the user balance.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User who needs to update the balance.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Quantity to be added or removed.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "Reason for the change.",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Check is command user in dm.
			if i.Interaction.User != nil {
				// Send a interaction respond message.
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This command cannot be used in dm!",
					},
				}); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond message")
				}

				return
			}

			// Checking if the user has the review role.
			if !hasRole(i.Interaction.Member.Roles, p.user.ReviewRole) {
				// Send a interaction respond message.
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have access to this command!",
					},
				}); err != nil {
					log.Warn().Err(err).Msg("failed to send interaction respond message")
				}

				return
			}

			// Updating the user balance.
			if err := p.service.UpdateBalance(
				context.Background(),
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
				int(i.ApplicationCommandData().Options[1].IntValue()),
			); err != nil {
				log.Error().Err(err).Msg("failed to updating user balance")
			}

			// Send a interaction respond message.
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						"You have updated the balance of user <@%s> on `%d`",
						i.ApplicationCommandData().Options[0].UserValue(s).ID,
						i.ApplicationCommandData().Options[1].IntValue(),
					),
				},
			}); err != nil {
				log.Warn().Err(err).Msg("failed to send interaction respond message")
			}

			// Send bot log message.
			if _, err := s.ChannelMessageSendEmbed(
				p.bot.LogChannel,
				&discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						URL:     "https://discord.com/users/" + i.Interaction.Member.User.ID,
						Name:    i.Interaction.Member.User.Username,
						IconURL: i.Interaction.Member.User.AvatarURL("128x128"),
					},
					Description: fmt.Sprintf(
						"User balance <@%s> has been updated on `%d`.",
						i.ApplicationCommandData().Options[0].UserValue(s).ID,
						i.ApplicationCommandData().Options[1].IntValue(),
					),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Reason",
							Value: "> " + i.ApplicationCommandData().Options[2].StringValue(),
						},
					},
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

// Check is target role in the list of roles.
func hasRole(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}

	return false
}
