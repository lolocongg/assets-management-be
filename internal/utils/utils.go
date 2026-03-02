package utils

import (
	"strings"
	"time"
)

func NormalizeOrder(order string) string {
	switch strings.ToUpper(order) {
	case "ASC":
		return "ASC"
	default:
		return "DESC"
	}
}

func FormatDate(raw string) string {
	layouts := []string{
		"2006-01-02 15:04:05.000 -0700",
		"2006-01-02 15:04:05 -0700",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Format("02/01/2006")
		}
	}

	return raw // fallback nếu parse fail
}
