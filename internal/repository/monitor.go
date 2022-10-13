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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongodb database collection.
const monitorCollection string = "monitor"

type Monitor interface {
	// Getting promo monitor.
	Get(ctx context.Context, id int, last bool) (domain.Monitor, error)
	// Updating promo monitor.
	Update(ctx context.Context, monitor domain.Monitor) error
}

// Monitor repository structure.
type MonitorRepository struct{ coll *mongo.Collection }

// Creating a new monitor repository.
func NewMonitorRepository(db *mongo.Database) *MonitorRepository {
	return &MonitorRepository{coll: db.Collection(monitorCollection)}
}

// Getting promo monitor.
func (r *MonitorRepository) Get(ctx context.Context, id int, last bool) (domain.Monitor, error) {
	var (
		monitor domain.Monitor
		filter  interface{}
		opts    *options.FindOneOptions
	)

	// Checking is last options specified.
	if last {
		filter = bson.M{}
		opts = options.FindOne().SetSort(bson.M{"_id": -1})
	} else {
		filter = bson.M{"_id": id}
	}

	// Find monitor epoch.
	if err := r.coll.FindOne(
		ctx,
		filter,
		opts,
	).Decode(&monitor); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Monitor{}, &domain.Error{Code: domain.CodeNotFound, Message: "Epoch not found."}
		}

		return domain.Monitor{}, err
	}

	return monitor, nil
}

// Updating promo monitor.
func (r *MonitorRepository) Update(ctx context.Context, monitor domain.Monitor) error {
	updateQuery := bson.M{}

	// Checking is reward specified.
	if monitor.Reward != 0 {
		updateQuery["reward"] = monitor.Reward
	}
	// Checking is started in specified.
	if !monitor.StartedIn.IsZero() {
		updateQuery["startedIn"] = monitor.StartedIn
	}

	updateQuery["usageLimit"] = monitor.UsageLimit
	updateQuery["updatedAt"] = monitor.UpdatedAt

	// Update promo monitor.
	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": monitor.Id},
		bson.M{"$set": updateQuery},
		options.Update().SetUpsert(true),
	)

	return err
}
