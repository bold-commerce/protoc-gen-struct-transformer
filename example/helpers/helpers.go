// Package helpers contains manually written functions for transforming custom
// types. Such package should be written by end-user because plugin has no idea
// which types could be used. It generates just function names for non-standard
// types.
package helpers

import (
	"time"

	"github.com/bold-commerce/protoc-gen-struct-transformer/example/nulls"
)

func TimeToNullsTime(t time.Time) nulls.Time {
	return nulls.Time{Time: t}
}

func NullsTimeToTime(nt nulls.Time) time.Time {
	return nt.Time
}

func TimePtrToNullsTimePtr(t *time.Time) *nulls.Time {
	return &nulls.Time{Time: *t}
}

func NullsTimePtrToTimePtr(nt *nulls.Time) *time.Time {
	return &nt.Time
}

func TimePtrToNullsTime(t *time.Time) nulls.Time {
	return nulls.Time{Time: *t}
}

func NullsTimeToTimePtr(nt nulls.Time) *time.Time {
	return &nt.Time
}
