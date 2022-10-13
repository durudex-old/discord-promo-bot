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

package monitor

import (
	"context"
	"fmt"

	"github.com/durudex/discord-promo-bot/internal/bot/response"
	"github.com/durudex/discord-promo-bot/pkg/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Epoch bot command.
func (p *MonitorPlugin) EpochCommand() {
	// Registering a new discord application command.
	if err := p.bot.RegisterCommand(&bot.Command{
		ApplicationCommand: p.epochCommandApplication(),
		Handler:            p.epochCommandHandler,
	}); err != nil {
		log.Error().Err(err).Msg("failed to register command")
	}
}

// Epoch command application.
func (p *MonitorPlugin) epochCommandApplication() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "epoch",
		Description: "The command outputs all public information about the monitor epoch",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "epoch",
				Description: "Get information about the specified epoch.",
				Required:    false,
			},
		},
	}
}

// Epoch command handler.
func (p *MonitorPlugin) epochCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		id      int
		current bool
	)

	// Setting the monitor query options.
	if i.ApplicationCommandData().Options != nil {
		id = int(i.ApplicationCommandData().Options[0].IntValue())
	} else {
		current = true
	}

	// Getting a promo monitor.
	monitor, err := p.service.Get(context.Background(), id, current, false)
	if err != nil {
		// Send a interaction respond error message.
		if err := response.InteractionError(s, i, err); err != nil {
			log.Warn().Err(err).Msg("failed to getting promo monitor")
		}

		return
	}

	// Send a interaction respond message.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Epoch %d", monitor.Id),
					Description: fmt.Sprintf("**Reward:** %d\n", monitor.Reward) +
						fmt.Sprintf("**Usage Limit:** %d\n", monitor.UsageLimit) +
						fmt.Sprintf("**Started In:** <t:%d:R>\n", monitor.StartedIn.Unix()) +
						fmt.Sprintf("**Updated At:** <t:%d:R>\n", monitor.UpdatedAt.Unix()),
					Color: p.botCfg.Color,
				},
			},
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}
