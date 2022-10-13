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
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/bot"
)

// User command plugin structure.
type UserPlugin struct {
	// Bot structure.
	bot *bot.Bot
	// User config variables.
	userCfg *config.UserConfig
	// Bot config variables.
	botCfg *config.BotConfig
	// User service.
	service service.User
}

// Creating a new user command plugin.
func NewUserPlugin(bot *bot.Bot, cfg *config.Config, service service.User) *UserPlugin {
	return &UserPlugin{bot: bot, userCfg: &cfg.User, botCfg: &cfg.Bot, service: service}
}

// Registering all user plugin commands.
func (p *UserPlugin) RegisterCommands() {
	// Register user plugin register bot command.
	p.RegisterCommand()
	// Register user plugin get user bot command.
	p.UserCommand()
	// Register user plugin create bot command.
	p.CreateCommand()
	// Register user plugin use bot command.
	p.UseCommand()
	// Register user plugin update balance command.
	p.UpdateBalanceCommand()
}
