package nulls

import "time"

// Time represents custom time type.
type Time struct {
	Time  time.Time
	Valid bool
}
