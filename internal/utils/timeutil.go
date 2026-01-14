package utils

import "database/sql"

func FormatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("2006-01-02T15:04:05Z07:00")
}

func FormatNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}
