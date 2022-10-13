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

// Use bot command.
func (p *UserPlugin) UseCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.useCommandApplication(),
		Handler:            p.useCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// Use command application.
func (p *UserPlugin) useCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "use",
		Description: "The command use a user promo code.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "promo",
				Description: "Promo code.",
				Required:    true,
			},
		},
	}
}

// Use command handler.
func (p *UserPlugin) useCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		if err := response.InteractionError(s, i, err); err != nil {
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
		p.botCfg.LogChannel,
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
			Color: p.botCfg.Color,
		},
	); err != nil {
		log.Warn().Err(err).Msg("failed to send channel message")
	}
}
