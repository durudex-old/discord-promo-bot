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

import "github.com/durudex/discord-promo-bot/pkg/command"

// Discord command plugin structure.
type Plugin struct{}

// Creating a new discord command plugin.
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Registering all discord commands.
func (p *Plugin) RegisterPlugins(handler *command.Handler) error {
	// Register promo commands.
	if err := NewPromoPlugin().RegisterCommands(handler); err != nil {
		return err
	}
	// Register bot commands.
	return NewBotPlugin().RegisterCommands(handler)
}
