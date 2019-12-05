// Package helpers contains manually written functions for transforming custom
// types. Such package should be written by end-user because plugin has no idea
// which types could be used. It generates just function names for non-standard
// types.
package helpers

import (
	"strconv"
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

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// StringToInt32 converts string to int32. It doesn't return an error for now,
// and if string is not correct or value is aut of range of int32 it will return
// default int32 value which is 0.
func StringToInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}

	return int32(i)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 converts string to int63. For details see comments for
// StringToInt32 function.
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}
