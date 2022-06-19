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

	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/repository"
)

// User service interface.
type User interface {
	Create(ctx context.Context, user domain.User) error
}

// User service structure.
type UserService struct{ repos repository.User }

// Creating a new user service.
func NewUserService(repos repository.User) *UserService {
	return &UserService{repos: repos}
}

// Creating a new user.
func (s *UserService) Create(ctx context.Context, user domain.User) error {
	return s.repos.CreateUser(ctx, user)
}
