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
	"testing"
	"time"
)

func TestRateLimitBackOff(t *testing.T) {
	rl := rateLimiter{}
	for i := 0; i < 3; i++ {
		rl.BackOff()
		expectInterval(t, rl, 10*time.Second)
	}
	for i := 0; i < 2; i++ {
		rl.BackOff()
		expectInterval(t, rl, 30*time.Second)
	}
	for i := 0; i < 5; i++ {
		rl.BackOff()
		expectInterval(t, rl, 1*time.Minute)
	}
	for i := 0; i < 2; i++ {
		rl.BackOff()
		expectInterval(t, rl, 5*time.Minute)
	}
	for i := 0; i < 3; i++ {
		rl.BackOff()
		expectInterval(t, rl, 10*time.Minute)
	}
	for i := 0; i < 2; i++ {
		rl.BackOff()
		expectInterval(t, rl, 30*time.Minute)
	}
	for i := 0; i < 100; i++ {
		rl.BackOff()
		expectInterval(t, rl, 1*time.Hour)
	}
}

func TestRateLimitResume(t *testing.T) {
	rl := rateLimiter{interval: 1 * time.Second}
	rl.Resume()
	expectInterval(t, rl, 0)
}

func TestRateLimitWait(t *testing.T) {
	rl := rateLimiter{interval: 500 * time.Millisecond}
	start := time.Now()
	rl.Wait()
	if since := time.Since(start); since < rl.interval {
		t.Errorf("Expected wait for %v, but has only been %v", rl.interval, since)
	}
}

func expectInterval(t *testing.T, rl rateLimiter, expected time.Duration) {
	if rl.interval != expected {
		t.Errorf("Expected interval to be %v but was actually %v", expected, rl.interval)
	}
}
