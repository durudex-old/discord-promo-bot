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

package repository

import (
	"context"

	"github.com/durudex/discord-promo-bot/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo db database collection.
const userCollection string = "user"

// User repository interface.
type User interface {
	Create(ctx context.Context, user domain.User) error
}

// User repository structure.
type UserRepository struct{ coll *mongo.Collection }

// Creating a new user repository.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection(userCollection)}
}

// Creating a new user.
func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)

	return err
}
