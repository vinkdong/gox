package encrypt

import "encoding/base64"

func Base64StringEncode(origin string) string {
	return base64.StdEncoding.EncodeToString([]byte(origin))
}

func Base64StringDecode(encryptedString string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}
	return string(data[:]), err
}

func Base64Encode(origin []byte) string {
	return base64.StdEncoding.EncodeToString(origin)
}
