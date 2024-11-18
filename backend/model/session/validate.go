package session

var idCharsetSet = func() map[rune]bool {
	s := make(map[rune]bool, idCharsetLen)
	for _, r := range idCharset {
		s[r] = true
	}
	return s
}()

// ValidateID returns true if the provided sessionID
// meets the application requirements, else false.
func ValidateID(id string) bool {
	if len(id) != idLen {
		return false
	}
	for _, r := range id {
		if !idCharsetSet[r] {
			return false
		}
	}

	return true
}
