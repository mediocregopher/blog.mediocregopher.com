package post

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date represents a calendar date with no timezone information attached.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// DateFromTime converts a Time into a Date, truncating all non-date
// information.
func DateFromTime(t time.Time) Date {
	t = t.UTC()
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

// ToTime converts a Date into a Time. The returned time will be UTC midnight of
// the Date.
func (d *Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(src interface{}) error {

	if src == nil {
		*d = Date{}
		return nil
	}

	ts, ok := src.(int64)

	if !ok {
		return fmt.Errorf("cannot scan value %#v into Date", src)
	}

	*d = DateFromTime(time.Unix(ts, 0))
	return nil
}

// Value implements the driver.Valuer interface.
func (d Date) Value() (driver.Value, error) {

	if d == (Date{}) {
		return nil, nil
	}

	return d.ToTime().Unix(), nil
}
