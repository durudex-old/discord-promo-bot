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
	"time"

	"github.com/durudex/discord-promo-bot/internal/bot/response"
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/pkg/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Register bot command.
func (p *UserPlugin) RegisterCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.registerCommandApplication(),
		Handler:            p.registerCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// Register command application.
func (p *UserPlugin) registerCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "register",
		Description: "The command with the help of which you can register in the bot.",
	}
}

// Register command handler.
func (p *UserPlugin) registerCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	if createdAt.Add(p.userCfg.MinAge).Unix() > time.Now().Unix() {
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
	if err := p.service.Create(context.Background(), domain.User{Id: author.ID}); err != nil {
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
			Content: "You have successfully registered!",
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}
