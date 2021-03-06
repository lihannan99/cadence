// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	metrics "github.com/uber/cadence/common/metrics"

	tally "github.com/uber-go/tally"

	time "time"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// AddCounter provides a mock function with given fields: scope, counter, delta
func (_m *Client) AddCounter(scope int, counter int, delta int64) {
	_m.Called(scope, counter, delta)
}

// IncCounter provides a mock function with given fields: scope, counter
func (_m *Client) IncCounter(scope int, counter int) {
	_m.Called(scope, counter)
}

// RecordTimer provides a mock function with given fields: scope, timer, d
func (_m *Client) RecordTimer(scope int, timer int, d time.Duration) {
	_m.Called(scope, timer, d)
}

// StartTimer provides a mock function with given fields: scope, timer
func (_m *Client) StartTimer(scope int, timer int) tally.Stopwatch {
	ret := _m.Called(scope, timer)

	var r0 tally.Stopwatch
	if rf, ok := ret.Get(0).(func(int, int) tally.Stopwatch); ok {
		r0 = rf(scope, timer)
	} else {
		r0 = ret.Get(0).(tally.Stopwatch)
	}

	return r0
}

// Tagged provides a mock function with given fields: tags
func (_m *Client) Tagged(tags map[string]string) metrics.Client {
	ret := _m.Called(tags)

	var r0 metrics.Client
	if rf, ok := ret.Get(0).(func(map[string]string) metrics.Client); ok {
		r0 = rf(tags)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metrics.Client)
		}
	}

	return r0
}

// UpdateGauge provides a mock function with given fields: scope, gauge, value
func (_m *Client) UpdateGauge(scope int, gauge int, value float64) {
	_m.Called(scope, gauge, value)
}
