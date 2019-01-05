package read

import (
	"fmt"
	"time"
)

func ParseDateOrTime(input string) (out time.Time, err error) {
	if out, err = time.Parse("2006-01-02T15:04:05 MST", input); err == nil {
		return
	}
	if out, err = time.Parse("2006-01-02T15:04:05Z07:00", input); err == nil {
		return
	}
	if out, err = time.Parse("2006-01-02T15:04:05Z", input); err == nil {
		return
	}
	if out, err = time.Parse("2006-01-02T15:04:05", input); err == nil {
		return
	}
	if out, err = time.Parse("2006-01-02T15:04", input); err == nil {
		return
	}
	if out, err = time.Parse("2006-01-02", input); err == nil {
		return
	}
	return out, fmt.Errorf("could not parse timestamp `%s`", input)
}
