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
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Promo commands plugin structure.
type PromoPlugin struct{ service service.Promo }

// Creating a new promo commands plugin.
func NewPromoPlugin(service service.Promo) *PromoPlugin {
	return &PromoPlugin{service: service}
}

// Registering promo plugin commands.
func (p *PromoPlugin) RegisterCommands(handler *command.Handler) error {
	// Register create promo command.
	if err := handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "create",
			Description: "The command creating a new user promo code.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "promo",
					Description: "Promo code.",
					Required:    true,
				},
			},
		},
		Handler: p.create,
	}); err != nil {
		return err
	}
	// Register use promo command.
	if err := handler.RegisterCommand(&command.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
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
		},
		Handler: p.use,
	}); err != nil {
		return err
	}

	return nil
}

// Creating a new promo code handler.
func (p *PromoPlugin) create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send a interaction respond message.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "test",
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}

// Use a promo code handler.
func (p *PromoPlugin) use(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send a interaction respond message.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "test",
		},
	}); err != nil {
		log.Warn().Err(err).Msg("failed to send interaction respond message")
	}
}
