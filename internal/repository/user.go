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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo db database collection.
const userCollection string = "user"

// User repository interface.
type User interface {
	CreateUser(ctx context.Context, user domain.User) error
	UpdatePromo(ctx context.Context, discordId, promo string) error
	UsePromo(ctx context.Context, discordId, promo string, award int) error
}

// User repository structure.
type UserRepository struct{ coll *mongo.Collection }

// Creating a new user repository.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection(userCollection)}
}

// Creating a new user.
func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)

	return err
}

// Updating a user promo code.
func (r *UserRepository) UpdatePromo(ctx context.Context, discordId, promo string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"discordId": discordId}, bson.M{"$set": bson.M{"promo": promo}})

	return err
}

// Use a promo code.
func (r *UserRepository) UsePromo(ctx context.Context, discordId, promo string, award int) error {
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		var user domain.User

		// Find a user promo code.
		err := r.coll.FindOne(sessCtx, bson.M{"promo": promo}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, &domain.Error{Code: domain.CodeNotFound, Message: "Promo code not found."}
			}

			return nil, err
		}

		// Update a user uses promo and increment balance.
		if err := r.coll.FindOneAndUpdate(
			sessCtx,
			bson.M{"discordId": discordId},
			bson.M{"$set": bson.M{"used": promo}, "$inc": bson.M{"balance": award}},
		).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, &domain.Error{Code: domain.CodeNotFound, Message: "User not found."}
			}

			return nil, err
		}

		// Increment promo author balance.
		_, err = r.coll.UpdateByID(sessCtx, user.Id, bson.M{"$inc": bson.M{"balance": award}})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Creating a new mongodb session.
	session, err := r.coll.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Executing the callback.
	_, err = session.WithTransaction(ctx, callback)

	return err
}
