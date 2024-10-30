package lib

// StringUtf8ToIso88591 convert a UTF-8 string in a ISO-8859-1 string.
func StringUtf8ToIso88591(value string) string {
	bytes := []byte(value)
	return string(bytesUtf8ToIso88591(bytes))
}

// StringIso88591ToUtf8 convert a ISO-8859-1 string in a UTF-8 string.
func StringIso88591ToUtf8(value string) string {
	bytes := []byte(value)
	return string(bytesIso88591ToUtf8(bytes))
}

// bytesIso88591ToUtf8 converts a byte slice in ISO-8859-1 encoding to a UTF-8 encoded byte slice.
func bytesIso88591ToUtf8(value []byte) []byte {
	var utf8Output []byte
	for _, b := range value {
		if b < 0x80 { //nolint:mnd
			// ASCII characters are the same in UTF-8 and ISO-8859-1
			utf8Output = append(utf8Output, b)
		} else {
			// For characters in the ISO-8859-1 range (0x80 - 0xFF),
			// we need to convert them to two-byte UTF-8 encoding.
			utf8Output = append(utf8Output, 0xC0|(b>>6), 0x80|(b&0x3F)) //nolint:mnd
		}
	}
	return utf8Output
}

// bytesUtf8ToIso88591 convert a UTF-8 string in a sequence of bytes ISO-8859-1.
func bytesUtf8ToIso88591(value []byte) []byte {
	var isoOutput []byte
	for i := 0; i < len(value); i++ {
		b := value[i]
		switch {
		case b < 0x80: //nolint:mnd
			// ASCII compatible with ISO-8859-1
			isoOutput = append(isoOutput, b)
		case (b&0xE0) == 0xC0 && i+1 < len(value) && (value[i+1]&0xC0) == 0x80: //nolint:mnd
			// converts from UTF-8 to ISO-8859-1: two bytes (110xxxxx 10xxxxxx)
			isoByte := ((b & 0x1F) << 6) | (value[i+1] & 0x3F) //nolint:mnd
			isoOutput = append(isoOutput, isoByte)
			i++ // jump to next byte
		default:
			// characters without of ISO-8859-1 cannot be converted, ignore and replace by '?'
			isoOutput = append(isoOutput, '?')
		}
	}
	return isoOutput
}
