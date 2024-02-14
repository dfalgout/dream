package helpers

import "strings"

func CleanEmail(email *string) string {
	if email == nil {
		return ""
	}
	checkedEmail := strings.ToLower(*email)
	return strings.TrimSpace(checkedEmail)
}
