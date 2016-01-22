// Copyright © 2016 Matthew Holt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package issuance

import (
	"log"
	"time"
)

/*
type RateLimiter interface {
	BackOff()
	Resume()
	Wait()
}
*/

// rateLimiter implements a very simple rate limiter
// that throttles according to this schedule, with
// each interval used a certain number of times:
//
//     10s - 3x
//     30s - 2x
//     1m  - 5x
//     5m  - 2x
//     10m - 3x
//     30m - 2x
//     1h  - ∞x
//
// Upon resuming, the interval will be reset to 0.
type rateLimiter struct {
	interval time.Duration // how long to wait
	count    int           // how many times we've waited at this interval
}

// BackOff tells the rate limiter to throttle another step.
func (rl *rateLimiter) BackOff() {
	// Switch on current interval to determine the new one
	switch rl.interval {
	case 0:
		rl.count = 1
		rl.interval = 10 * time.Second
	case 10 * time.Second:
		rl.count++
		if rl.count > 3 {
			rl.count = 1
			rl.interval = 30 * time.Second
		}
	case 30 * time.Second:
		rl.count++
		if rl.count > 2 {
			rl.count = 1
			rl.interval = 1 * time.Minute
		}
	case 1 * time.Minute:
		rl.count++
		if rl.count > 5 {
			rl.count = 1
			rl.interval = 5 * time.Minute
		}
	case 5 * time.Minute:
		rl.count++
		if rl.count > 2 {
			rl.count = 1
			rl.interval = 10 * time.Minute
		}
	case 10 * time.Minute:
		rl.count++
		if rl.count > 3 {
			rl.count = 1
			rl.interval = 30 * time.Minute
		}
	case 30 * time.Minute:
		rl.count++
		if rl.count > 2 {
			rl.count = 1
			rl.interval = 1 * time.Hour
		}
	case 1 * time.Hour:
		rl.count++
	default:
		log.Println("[ERROR] Unexpected interval:", rl.interval)
	}
}

// Wait waits the duration of the interval.
func (rl *rateLimiter) Wait() {
	time.Sleep(rl.interval)
}

// Resume resets the interval back to 0.
func (rl *rateLimiter) Resume() {
	rl.interval = 0
}
