package a

func slice_valid(s []string) bool {
	if s != nil {
		s[0] = "a"
		return true
	}
	return false
}

func slice_invalid(s []string) bool {
	s[0] = "a" // want "nil check leakage"
	return true
}

func slice_valid_len(s []string) bool{
	if len(s) != 0{
		s[0] = "a"
		return true
	}
	return false

}