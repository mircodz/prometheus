// Copyright 2021 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textparse

// This file implements operations on google.protobuf.Timestamp.

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

// validateTimestamp determines whether a Timestamp is valid.
// A valid timestamp represents a time in the range
// [0001-01-01, 10000-01-01) and has a Nanos field
// in the range [0, 1e9).
//
// If the Timestamp is valid, validateTimestamp returns nil.
// Otherwise, it returns an error that describes
// the problem.
//
// Every valid Timestamp can be represented by a time.Time, but the converse is not true.
func validateTimestamp(ts *timestamppb.Timestamp) error {
	if ts == nil {
		return errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < minValidSeconds {
		return fmt.Errorf("timestamp: %#v before 0001-01-01", ts)
	}
	if ts.Seconds >= maxValidSeconds {
		return fmt.Errorf("timestamp: %#v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return fmt.Errorf("timestamp: %#v: nanos not in range [0, 1e9)", ts)
	}
	return nil
}

// timestampFromProto converts a google.protobuf.Timestamp proto to a time.Time.
// It returns an error if the argument is invalid.
//
// Unlike most Go functions, if Timestamp returns an error, the first return value
// is not the zero time.Time. Instead, it is the value obtained from the
// time.Unix function when passed the contents of the Timestamp, in the UTC
// locale. This may or may not be a meaningful time; many invalid Timestamps
// do map to valid time.Times.
//
// A nil Timestamp returns an error. The first return value in that case is
// undefined.
func timestampFromProto(ts *timestamppb.Timestamp) (time.Time, error) {
	// Don't return the zero value on error, because corresponds to a valid
	// timestamp. Instead return whatever time.Unix gives us.
	var t time.Time
	if ts == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}
	return t, validateTimestamp(ts)
}
