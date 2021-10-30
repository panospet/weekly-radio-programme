package common

import (
	"fmt"
	"strings"
	"time"
)

func HasConflict(input1 string, input2 string) (bool, error) {
	start1, end1, err := GetStartEnd(input1)
	if err != nil {
		return false, err
	}

	start2, end2, err := GetStartEnd(input2)
	if err != nil {
		return false, err
	}

	return start1.Before(end2) && end1.After(start2), nil
}

func GetStartEnd(input string) (time.Time, time.Time, error) {
	if len(input) != 11 {
		return time.Time{}, time.Time{}, fmt.Errorf(`bad input "%s"`, input)
	}

	parts := strings.Split(input, "-")
	startStr := parts[0]
	endStr := parts[1]

	if endStr == "00:00" {
		endStr = "23:59"
	}

	start, err := SmartParse(startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("cannot parse start string")
	}
	end, err := SmartParse(endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("cannot parse end string")
	}
	return start, end, nil
}

func SmartParse(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02")+" "+input)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse")
	}
	return t, nil
}
