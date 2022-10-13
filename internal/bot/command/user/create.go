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
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/pkg/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Create bot command.
func (p *UserPlugin) CreateCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.createCommandApplication(),
		Handler:            p.createCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// Create command application.
func (p *UserPlugin) createCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "create",
		Description: "The command creating a new user promo code.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "promo",
				Description: "Unique promo code.",
				Required:    true,
			},
		},
	}
}

// Create command handler.
func (p *UserPlugin) createCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		if err := response.InteractionError(s, i, err); err != nil {
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
		p.botCfg.LogChannel,
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
			Color: p.botCfg.Color,
		},
	); err != nil {
		log.Warn().Err(err).Msg("failed to send channel message")
	}
}
