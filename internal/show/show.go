package show

import (
	"fmt"
	"strings"
	"time"
	"weekly-radio-programme/common"
)

const (
	Monday    = "Mon"
	Tuesday   = "Tue"
	Wednesday = "Wed"
	Thursday  = "Thu"
	Friday    = "Fri"
	Saturday  = "Sat"
	Sunday    = "Sun"
)

var weekdays = []string{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}

type Show struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Weekday     string    `json:"weekday" db:"weekday"`
	Timeslot    string    `json:"timeslot" db:"timeslot"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

func (o Show) Validate() error {
	if !common.Contains(weekdays, o.Weekday) {
		return fmt.Errorf("invalid weekday")
	}
	if len(o.Timeslot) != 11 {
		return fmt.Errorf("invalid timeslot. Acceptable format is hh:mm-hh:mm")
	}
	parts := strings.Split(o.Timeslot, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid timeslot. Acceptable format is hh:mm-hh:mm")
	}
	if _, err := time.Parse("15:04", parts[0]); err != nil {
		return fmt.Errorf("invalid timeslot. %s not a valid time", parts[0])
	}
	if _, err := time.Parse("15:04", parts[1]); err != nil {
		return fmt.Errorf("invalid timeslot. %s not a valid time", parts[1])
	}

	return nil
}
