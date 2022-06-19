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

package service

import (
	"context"

	"github.com/durudex/discord-promo-bot/internal/config"
	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/repository"
)

// Promo service interface.
type Promo interface {
	Update(ctx context.Context, discordId, promo string) error
	Use(ctx context.Context, discordId, promo string) error
}

// Promo service structure.
type PromoService struct {
	repos repository.User
	cfg   config.PromoConfig
}

// Creating a new promo service.
func NewPromoService(repos repository.User, cfg config.PromoConfig) *PromoService {
	return &PromoService{repos: repos, cfg: cfg}
}

// Updating a user promo.
func (s *PromoService) Update(ctx context.Context, discordId, promo string) error {
	// Validate promo code.
	if !domain.RxPromo.MatchString(promo) {
		return &domain.Error{Code: domain.CodeInvalidArgument, Message: "The promo code is invalid."}
	}

	// Updating a user promo.
	if err := s.repos.UpdatePromo(ctx, discordId, promo); err != nil {
		return err
	}

	return nil
}

// Using a user promo.
func (s *PromoService) Use(ctx context.Context, discordId, promo string) error {
	if err := s.repos.UsePromo(ctx, discordId, promo, s.cfg.Award); err != nil {
		return err
	}

	return nil
}
