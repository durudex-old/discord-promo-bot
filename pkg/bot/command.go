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

// Discord application command structure.
type Command struct {
	// Discord application command.
	discordgo.ApplicationCommand
	// Discord bot application command handler.
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Registering a new discord application command.
func (b *Bot) RegisterCommand(c *Command) error {
	// Creating a new global discord application command.
	cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", &c.ApplicationCommand)
	if err != nil {
		return err
	}

	// Save the discord application command.
	b.commands[c.Name] = &Command{ApplicationCommand: *cmd, Handler: c.Handler}

	return nil
}

// Unregister all discord application commands.
func (b *Bot) UnregisterCommands() error {
	// Getting all discord application commands.
	commands, err := b.session.ApplicationCommands(b.session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, command := range commands {
		// Delete discord application commands.
		if err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", command.ID); err != nil {
			return err
		}
	}

	return nil
}
