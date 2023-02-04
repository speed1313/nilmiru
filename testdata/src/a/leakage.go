package a

func leakage(gopher *int) {
	// The pattern can be written in regular expression.
	print(gopher)  // want "nil check leakage"
}

