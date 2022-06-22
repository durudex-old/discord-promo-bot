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

package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/durudex/discord-promo-bot/internal/config"
)

// Test initialize config.
func TestConfig_Init(t *testing.T) {
	// Environment configurations.
	type env struct{ configPath, botToken, mongoUri, mongoUsername, mongoPassword string }

	// Testing args.
	type args struct{ env env }

	// Set environments configurations.
	setEnv := func(env env) {
		os.Setenv("CONFIG_PATH", env.configPath)
		os.Setenv("BOT_TOKEN", env.botToken)
		os.Setenv("MONGO_URI", env.mongoUri)
		os.Setenv("MONGO_USERNAME", env.mongoUsername)
		os.Setenv("MONGO_PASSWORD", env.mongoPassword)
	}

	// Tests structures.
	tests := []struct {
		name    string
		args    args
		want    *config.Config
		wantErr bool
	}{
		{
			name: "OK",
			args: args{env: env{
				configPath:    "fixtures/main",
				botToken:      "123",
				mongoUri:      "mongodb://localhost:27017",
				mongoUsername: "admin",
				mongoPassword: "qwerty",
			}},
			want: &config.Config{
				Bot: config.BotConfig{Token: "123"},
				Database: config.DatabaseConfig{
					Mongodb: config.MongodbConfig{
						URI:      "mongodb://localhost:27017",
						Username: "admin",
						Password: "qwerty",
						Timeout:  time.Second * 10,
						Database: "durudex",
					},
				},
				User:  config.UserConfig{MinAge: time.Hour * 1440},
				Promo: config.PromoConfig{Award: 100},
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environments configurations.
			setEnv(tt.args.env)

			// Initialize config.
			got, err := config.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("error initialize config: %s", err.Error())
			}

			// Check for similarity of a config.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("error config are not similar")
			}
		})
	}
}
