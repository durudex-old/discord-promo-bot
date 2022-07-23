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

package plugin

import (
	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/service"
	"github.com/durudex/discord-promo-bot/pkg/command"
)

// Discord command plugin structure.
type Plugin struct {
	service *service.Service
	cfg     *config.Config
}

// Discord dm commands permission.
var DMPermission bool = true

// Creating a new discord command plugin.
func NewPlugin(service *service.Service, cfg *config.Config) *Plugin {
	return &Plugin{service: service, cfg: cfg}
}

// Registering all discord commands.
func (p *Plugin) RegisterPlugins(handler *command.Handler) {
	// Registering user plugin commands.
	NewUserPlugin(p.service.User, handler, &p.cfg.User).RegisterCommands()
	// Register promo commands.
	NewPromoPlugin(p.service.Promo, handler).RegisterCommands()
	// Register bot commands.
	NewBotPlugin(handler).RegisterCommands()
}
