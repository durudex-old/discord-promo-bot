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
	"sync"
	"time"

	"github.com/durudex/discord-promo-bot/internal/domain"
	"github.com/durudex/discord-promo-bot/internal/repository"

	"github.com/rs/zerolog/log"
)

// Monitor service interface.
type Monitor interface {
	// Getting a promo monitor.
	Get(ctx context.Context, id int, current bool) (domain.Monitor, error)
	// Saving promo monitor.
	Save(ctx context.Context, skip bool, monitor ...domain.Monitor) error
	// Sync promo monitor with database.
	Sync(ctx context.Context) error
	// Using a promo code with monitor.
	Use() (int, error)
	// De using promo code with monitor.
	DeUse()
}

// Monitor service structure.
type MonitorService struct {
	// Monitor repository.
	repos repository.Monitor
	// Monitor structure.
	monitor *domain.Monitor
	// Sync monitor mutex.
	mutex sync.Mutex
	// Updated monitor status.
	updated bool
}

// Creating a new monitor service.
func NewMonitorService(repos repository.Monitor) *MonitorService {
	return &MonitorService{repos: repos, mutex: sync.Mutex{}}
}

// Getting a promo monitor.
func (s *MonitorService) Get(ctx context.Context, id int, current bool) (domain.Monitor, error) {
	if id < domain.MaxMonitorEpoch {
		return domain.Monitor{}, &domain.Error{
			Code:    domain.CodeInvalidArgument,
			Message: "There can be no more than 5 epochs.",
		}
	}

	return s.repos.Get(ctx, id, current)
}

// Saving promo monitor.
func (s *MonitorService) Save(ctx context.Context, skip bool, monitor ...domain.Monitor) error {
	if monitor != nil {
		// Updating promo monitor.
		err := s.repos.Update(ctx, monitor[0])
		return err
	} else if s.updated || skip {
		// Updating promo monitor.
		if err := s.repos.Update(ctx, domain.Monitor{
			Id:         s.monitor.Id,
			Reward:     s.monitor.Reward,
			UsageLimit: s.monitor.UsageLimit,
			StartedIn:  s.monitor.StartedIn,
			UpdatedAt:  time.Now(),
		}); err != nil {
			return err
		}

		s.updated = false
	}

	return nil
}

// Sync promo monitor with database.
func (s *MonitorService) Sync(ctx context.Context) error {
	// Getting current promo monitor.
	monitor, err := s.repos.Get(ctx, 0, true)
	if err != nil {
		return err
	}

	s.monitor = &monitor

	return nil
}

// Using a promo code with monitor.
func (s *MonitorService) Use() (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Checking for the end of the limit of using the promo code in the epoch.
	if s.monitor.UsageLimit == 0 {
		// Checking is max epoch.
		if s.monitor.Id == domain.MaxMonitorEpoch {
			return 0, &domain.Error{Code: domain.CodeNotFound, Message: "Rewards are over!"}
		}

		go func(mon domain.Monitor) {
			// Saving promo monitor.
			if err := s.Save(context.Background(), false, domain.Monitor{
				Id:         mon.Id,
				UsageLimit: mon.UsageLimit,
				UpdatedAt:  time.Now(),
			}); err != nil {
				log.Error().Err(err).Msg("error saving monitor")
			}
		}(*s.monitor)

		monitor := domain.Epochs[s.monitor.Id+1]

		s.monitor = &domain.Monitor{
			Id:         monitor.Id,
			Reward:     monitor.Reward,
			UsageLimit: monitor.UsageLimit,
			StartedIn:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	s.monitor.UsageLimit--
	s.updated = true

	return s.monitor.Reward, nil
}

// De using promo code with monitor.
func (s *MonitorService) DeUse() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.monitor.UsageLimit++
}
