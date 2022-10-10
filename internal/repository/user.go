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

package repository

import (
	"context"

	"github.com/durudex/discord-promo-bot/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongodb database collection.
const userCollection string = "test"

// User repository interface.
type User interface {
	// Creating a new user.
	Create(ctx context.Context, user domain.User) error
	// Getting a user.
	Get(ctx context.Context, id string) (domain.User, error)
	// Updating a user promo code.
	UpdatePromo(ctx context.Context, id, promo string) error
	// Using a promo code.
	UsePromo(ctx context.Context, id, promo string, reward int) error
	// Updating a user balance.
	UpdateBalance(ctx context.Context, id string, amount int) error
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
	if mongo.IsDuplicateKeyError(err) {
		return &domain.Error{Code: domain.CodeAlreadyExists, Message: "You are registered."}
	}

	return err
}

// Getting a user.
func (r *UserRepository) Get(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, &domain.Error{Code: domain.CodeNotFound, Message: "User not found."}
		}

		return domain.User{}, err
	}

	return user, nil
}

// Updating a user promo code.
func (r *UserRepository) UpdatePromo(ctx context.Context, id, promo string) error {
	if err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id, "promo": nil},
		bson.M{"$set": bson.M{"promo": promo}},
	).Err(); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &domain.Error{Code: domain.CodeAlreadyExists, Message: "The promo already exists."}
		} else if err == mongo.ErrNoDocuments {
			return &domain.Error{Code: domain.CodeNotFound, Message: "User does not exist or has already created a promo code."}
		}

		return err
	}

	return nil
}

// Using a promo code.
func (r *UserRepository) UsePromo(ctx context.Context, id, promo string, reward int) error {
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

		// Check if is author.
		if user.Id == id {
			return nil, &domain.Error{Code: domain.CodeInvalidArgument, Message: "You can't use your own promo code."}
		}

		// Update a user used promo and increment balance.
		if err := r.coll.FindOneAndUpdate(
			sessCtx,
			bson.M{"_id": id, "used": nil},
			bson.M{"$set": bson.M{"used": promo}, "$inc": bson.M{"balance": reward}},
		).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, &domain.Error{Code: domain.CodeNotFound, Message: "User does not exist or has already used the promo code."}
			}

			return nil, err
		}

		// Increment promo author balance.
		_, err = r.coll.UpdateByID(sessCtx, user.Id, bson.M{"$inc": bson.M{"balance": reward}})
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

// Updating a user balance.
func (r *UserRepository) UpdateBalance(ctx context.Context, discordId string, amount int) error {
	// Update a user used balance.
	if err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": discordId}, bson.M{"$inc": bson.M{"balance": amount}}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.Error{Code: domain.CodeNotFound, Message: "User does not exist."}
		}

		return err
	}

	return nil
}
