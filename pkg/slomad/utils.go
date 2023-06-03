package slomad

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// StringValOr returns the value of a string pointer if it's not nil, or a
// default value otherwise.
func StringValOr(sp *string, val string) string {
	if sp != nil {
		return *sp
	}
	return val
}

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}

// IntValOr returns the value of an int pointer if it's not nil, or a default.
func IntValOr(ip *int, val int) int {
	if ip != nil {
		return *ip
	}
	return val
}
