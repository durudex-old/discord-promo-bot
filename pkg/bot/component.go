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

package bot

import "github.com/bwmarrin/discordgo"

// Discord message component structure.
type Component struct {
	// Custom component id.
	ComponentID string
	// Discord message component handler.
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Registering a new discord message component.
func (b *Bot) RegisterComponent(c *Component) {
	b.components[c.ComponentID] = c
}
