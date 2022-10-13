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

	"github.com/durudex/discord-promo-bot/pkg/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var UpdateCommandMemberPermission int64 = discordgo.PermissionManageMessages

// Update balance bot command
func (p *UserPlugin) UpdateBalanceCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.updateBalanceCommandApplication(),
		Handler:            p.updateBalanceCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// Update balance command application.
func (p *UserPlugin) updateBalanceCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:                     "update-balance",
		Description:              "The command updating the user balance.",
		DefaultMemberPermissions: &UpdateCommandMemberPermission,
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
	}
}

// Register command handler.
func (p *UserPlugin) updateBalanceCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	if !hasRole(i.Interaction.Member.Roles, p.userCfg.ReviewRole) {
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
		p.botCfg.LogChannel,
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
			Color: p.botCfg.Color,
		},
	); err != nil {
		log.Warn().Err(err).Msg("failed to send channel message")
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
