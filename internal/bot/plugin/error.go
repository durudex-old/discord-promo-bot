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
	"errors"

	"github.com/durudex/discord-promo-bot/internal/domain"

	"github.com/bwmarrin/discordgo"
)

// Sending a interaction discord error message.
func discordInteractionError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errorHandler(err),
		},
	})
}

// Bot plugin error handler.
func errorHandler(err error) string {
	var e *domain.Error

	// Check if error is a domain.Error.
	if errors.As(err, &e) {
		switch e.Code {
		case domain.CodeNotFound:
			return e.Message
		case domain.CodeAlreadyExists:
			return e.Message
		case domain.CodeInvalidArgument:
			return e.Message
		case domain.CodeInternal:
			return "Internal bot error"
		}
	}

	return "Internal bot error"
}
