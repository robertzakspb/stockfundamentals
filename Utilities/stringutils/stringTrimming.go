package stringutils

func TrimFirstCharacter(str string) string {
	for i := range str {
		if i > 0 {
			return str[i:]
		}
	}

	return str
}

func TrimLastCharacter(str string) string {
	for i := range str {
		if i > 0 {
			return str[:len(str)-1]
		}
	}

	return str
}
