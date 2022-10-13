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
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/bot"
)

// Monitor command plugin structure.
type MonitorPlugin struct {
	// Bot structure.
	bot *bot.Bot
	// Bot config variables.
	botCfg *config.BotConfig
	// Monitor service.
	service service.Monitor
}

// Creating a new monitor service.
func NewMonitorPlugin(bot *bot.Bot, cfg *config.BotConfig, service service.Monitor) *MonitorPlugin {
	return &MonitorPlugin{bot: bot, botCfg: cfg, service: service}
}

// Registering all monitor plugin commands.
func (p *MonitorPlugin) RegisterCommands() {
	// Register monitor plugin epoch bot command.
	p.EpochCommand()
}
