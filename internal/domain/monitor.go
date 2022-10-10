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

package domain

import "time"

// Max promo epoch.
const MaxMonitorEpoch int = 5

// Promo monitor epochs.
var Epochs map[int]*Monitor = map[int]*Monitor{
	1: {Id: 1, Reward: 1000, UsageLimit: 500},
	2: {Id: 2, Reward: 900, UsageLimit: 2000},
	3: {Id: 3, Reward: 800, UsageLimit: 2500},
	4: {Id: 4, Reward: 700, UsageLimit: 10000},
	5: {Id: 5, Reward: 600, UsageLimit: 10000},
}

// Promo monitor structure.
type Monitor struct {
	// Promo epoch id.
	Id int `bson:"_id"`
	// Promo reward.
	Reward int `bson:"reward"`
	// Promo usage limit.
	UsageLimit int `bson:"usageLimit"`
	// Promo epoch started in.
	StartedIn time.Time `bson:"startedIn,omitempty"`
	// Updated at promo monitor.
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
}
