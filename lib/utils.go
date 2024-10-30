package lib

// // sanitizeUTF8 verify se a string is UTF-8. If not, convert to a valid UTF-8 Sequence.
// func sanitizeUTF8(s string) string {
// 	if utf8.ValidString(s) {
// 		return s
// 	}
// 	return string([]rune(s))
// }

// // convertToUTF8 Convert uma string ISO-8859-1 para UTF-8
// func convertToUTF8(s string) (string, error) {
// 	reader := transform.NewReader(bytes.NewReader([]byte(s)), charmap.ISO8859_1.NewDecoder())
// 	utf8Bytes, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(utf8Bytes), nil
// }
