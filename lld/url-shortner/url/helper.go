package url

func encodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	var result []byte
	base := int64(len(base62Chars))

	for num > 0 {
		rem := num % base
		result = append([]byte{base62Chars[rem]}, result...)
		num = num / base
	}

	return string(result)
}

func toFixedLength(code string, length int) string {
	if len(code) >= length {
		return code
	}

	padding := make([]byte, length-len(code))
	for i := range padding {
		padding[i] = '0'
	}

	return string(padding) + code
}
