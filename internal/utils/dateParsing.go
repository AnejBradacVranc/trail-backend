package utils

import "time"

func ParseTimeString(value string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05.999999-07",
		time.RFC3339,
		"2006-01-02",
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, err
}
