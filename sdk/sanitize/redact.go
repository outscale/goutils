package sanitize

type redactFunc func(string) string

var Redacted = "[REDACTED]"

func RedactAll(s string) string {
	return Redacted
}

func KeepFirst2Last2(s string) string {
	if len(s) <= 4 {
		return Redacted
	}
	return s[:2] + "..." + s[len(s)-2:]
}
