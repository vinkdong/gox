package scram_sha_256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Hi function as defined in RFC5802.
func Hi(password, salt []byte, iterations int) ([]byte, error) {
	mac := hmac.New(sha256.New, password)
	if _, err := mac.Write(salt); err != nil {
		return nil, err
	}
	if _, err := mac.Write([]byte{0, 0, 0, 1}); err != nil {
		return nil, err
	}
	u := mac.Sum(nil)
	f := make([]byte, len(u))
	copy(f, u)
	for i := 1; i < iterations; i++ {
		mac.Reset()
		if _, err := mac.Write(u); err != nil {
			return nil, err
		}
		u = mac.Sum(nil)
		for j := 0; j < len(f); j++ {
			f[j] ^= u[j]
		}
	}
	return f, nil
}

func ComputeClientProof(password, saltBase64 string, iterations int, authMessage string) (string, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		panic("Invalid base64 encoding for salt")
	}
	saltedPassword, err := Hi([]byte(password), salt, iterations)
	if err != nil {
		return "", err
	}
	return ComputeClientSaltedPasswordProof(saltedPassword, authMessage)
}

func ComputeClientSaltedPasswordProof(saltedPassword []byte, authMessage string) (string, error) {
	// ClientKey
	clientKey := hmac.New(sha256.New, saltedPassword)
	if _, err := clientKey.Write([]byte("Client Key")); err != nil {
		return "", err
	}

	storedKey := sha256.Sum256(clientKey.Sum(nil))

	// AuthMessage
	authMsg := []byte(authMessage)

	// ClientSignature
	clientSignature := hmac.New(sha256.New, storedKey[:])
	if _, err := clientSignature.Write(authMsg); err != nil {
		return "", err
	}

	// ClientProof
	clientProof := xor(clientKey.Sum(nil), clientSignature.Sum(nil))

	// Encode to base64
	encodedClientProof := base64.StdEncoding.EncodeToString(clientProof)
	return encodedClientProof, nil
}

func xor(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func ComputeServerSignature(password, saltBase64 string, authMessage string, iterations int) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		panic("Invalid base64 encoding for salt")
	}
	saltedPassword, err := Hi([]byte(password), salt, iterations)
	if err != nil {
		return nil, err
	}
	serverKey := hmac.New(sha256.New, saltedPassword)
	if _, err := serverKey.Write([]byte("Server Key")); err != nil {
		return nil, err
	}

	serverSignature := hmac.New(sha256.New, serverKey.Sum(nil))
	if _, err := serverSignature.Write([]byte(authMessage)); err != nil {
		return nil, err
	}

	return serverSignature.Sum(nil), nil
}

func ComputeServerSaltedPasswordSignature(saltedPassword []byte, authMessage string) ([]byte, error) {
	serverKey := hmac.New(sha256.New, saltedPassword)
	if _, err := serverKey.Write([]byte("Server Key")); err != nil {
		return nil, err
	}

	serverSignature := hmac.New(sha256.New, serverKey.Sum(nil))
	if _, err := serverSignature.Write([]byte(authMessage)); err != nil {
		return nil, err
	}

	return serverSignature.Sum(nil), nil
}

func VerifyServerProof(serverProofBase64 string, saltBase64, password string, iterations int, authMessage string) (bool, error) {
	serverProof, err := base64.StdEncoding.DecodeString(serverProofBase64)
	if err != nil {
		return false, err
	}

	expectedServerSignature, err := ComputeServerSignature(password, saltBase64, authMessage, iterations)
	if err != nil {
		return false, err
	}

	return hmac.Equal(serverProof, expectedServerSignature), nil
}
