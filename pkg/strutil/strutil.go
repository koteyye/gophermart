package strutil

// OnlyDigits возвращает true, если строка не пуста и содержит только цифры.
func OnlyDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
