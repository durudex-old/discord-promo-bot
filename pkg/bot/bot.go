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

// Bot structure.
type Bot struct {
	// Discord bot session.
	session *discordgo.Session
	// Discord bot application commands.
	commands map[string]*Command
	// Discord bot application commands components.
	components map[string]*Component
}

// Creating a new discord bot.
func New(token string) (*Bot, error) {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:    session,
		commands:   make(map[string]*Command),
		components: make(map[string]*Component),
	}, nil
}

// Running the discord bot.
func (b *Bot) Run() error {
	return b.session.Open()
}

// Closing a bot connections.
func (b *Bot) Close() error {
	return b.session.Close()
}

// Registering a discord bot handler.
func (b *Bot) RegisterHandler(handler any) func() {
	return b.session.AddHandler(handler)
}

// Handle discord application command.
func (b *Bot) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// Handle discord bot application command.
		if c, ok := b.commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		// Handle discord message component.
		if c, ok := b.components[i.MessageComponentData().CustomID]; ok {
			c.Handler(s, i)
		}
	}
}
