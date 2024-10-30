package lib

// stringUtf8ToIso88591 convert a UTF-8 string in a sequence of bytes ISO-8859-1.
func stringUtf8ToIso88591(value string) string {
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
	return string(isoOutput)
}
