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

import "regexp"

// Regular expression for promo code.
const Promo string = "^[a-z0-9-_.]{3,12}$"

var RxPromo = regexp.MustCompile(Promo)

// User structure.
type User struct {
	// User discord id.
	Id string `bson:"_id"`
	// User own promo code.
	Promo string `bson:"promo,omitempty"`
	// User used promo code.
	Used string `bson:"used,omitempty"`
	// User token balance.
	Balance int `bson:"balance,omitempty"`
}

// Validating a user.
func (u User) Validate() error {
	switch {
	case !RxPromo.MatchString(u.Promo):
		return &Error{Code: CodeInvalidArgument, Message: "The promo code is invalid."}
	default:
		return nil
	}
}
