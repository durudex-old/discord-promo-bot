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

package service

import (
	"context"

	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/repository"
)

// User service interface.
type User interface {
	// Creating a new user.
	Create(ctx context.Context, user domain.User) error
	// Getting a user.
	Get(ctx context.Context, id string) (domain.User, error)
	// Updating a user.
	Update(ctx context.Context, user domain.User) error
	// Using a user promo.
	UsePromo(ctx context.Context, discordId, promo string) (int, error)
	// Updating a user balance.
	UpdateBalance(ctx context.Context, id string, amount int) error
}

// User service structure.
type UserService struct {
	// User repository.
	repos repository.User
	// Monitor service.
	monitor Monitor
}

// Creating a new user service.
func NewUserService(repos repository.User, monitor Monitor) *UserService {
	return &UserService{repos: repos, monitor: monitor}
}

// Creating a new user.
func (s *UserService) Create(ctx context.Context, user domain.User) error {
	return s.repos.Create(ctx, user)
}

// Getting a user.
func (s *UserService) Get(ctx context.Context, id string) (domain.User, error) {
	return s.repos.Get(ctx, id)
}

// Updating a user.
func (s *UserService) Update(ctx context.Context, user domain.User) error {
	// Validating a user.
	if err := user.Validate(); err != nil {
		return err
	}

	return s.repos.UpdatePromo(ctx, user.Id, user.Promo)
}

// Using a user promo.
func (s *UserService) UsePromo(ctx context.Context, discordId, promo string) (int, error) {
	// Using a promo code with monitor.
	reward, err := s.monitor.Use()
	if err != nil {
		return 0, err
	}

	// Using a promo code.
	if err := s.repos.UsePromo(ctx, discordId, promo, reward); err != nil {
		s.monitor.DeUse()
		return 0, err
	}

	return reward, nil
}

// Updating a user balance.
func (s *UserService) UpdateBalance(ctx context.Context, id string, amount int) error {
	return s.repos.UpdateBalance(ctx, id, amount)
}
