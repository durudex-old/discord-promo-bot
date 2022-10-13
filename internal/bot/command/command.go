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

package command

import (
	"github.com/durudex/discord-promo-bot/internal/bot/command/basic"
	"github.com/durudex/discord-promo-bot/internal/bot/command/user"
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/bot"
)

// Command plugin structure.
type CommandPlugin struct {
	// Bot structure.
	bot *bot.Bot
	// Config variables.
	cfg *config.Config
	// Service structure.
	service *service.Service
}

// Creating a new command plugin.
func NewCommandPlugin(bot *bot.Bot, cfg *config.Config, service *service.Service) *CommandPlugin {
	return &CommandPlugin{bot: bot, cfg: cfg, service: service}
}

// Registering command plugins.
func (p *CommandPlugin) Register() {
	// Registering all basic commands.
	basic.NewBasicPlugin(p.bot).RegisterCommands()
	// Registering all user commands.
	user.NewUserPlugin(p.bot, p.cfg, p.service.User).RegisterCommands()
}
