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

package user

import (
	"context"
	"fmt"

	"github.com/durudex/discord-promo-bot/internal/bot/response"
	"github.com/durudex/discord-promo-bot/pkg/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// User bot command.
func (p *UserPlugin) UserCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.userCommandApplication(),
		Handler:            p.userCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// User command application.
func (p *UserPlugin) userCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "user",
		Description: "The command getting public information about the user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Get information about the specified user.",
				Required:    false,
			},
		},
	}
}

// User command handler.
func (p *UserPlugin) userCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		if err := response.InteractionError(s, i, err); err != nil {
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
					Color: p.botCfg.Color,
				},
			},
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}
