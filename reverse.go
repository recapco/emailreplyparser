package emailreplyparser

func reverse(input string) string { // note: this will not work well with combining characters
	reversed := []rune(input)

	for i, j := 0, len(reversed)-1; i < len(reversed)/2; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return string(reversed)
}
