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

package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongodb config structure.
type MongoConfig struct {
	URI      string
	Username string
	Password string
	Timeout  time.Duration
}

// Creating a new mongodb client.
func NewClient(cfg *MongoConfig) (*mongo.Client, error) {
	// Creating a new mongodb client options.
	opts := options.Client().ApplyURI(cfg.URI)
	if cfg.Username != "" && cfg.Password != "" {
		// Set client auth options.
		opts.SetAuth(options.Credential{Username: cfg.Username, Password: cfg.Password})
	}

	// Creating a new mongodb client.
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Connecting to the mongodb server.
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the mongodb client is connected.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
