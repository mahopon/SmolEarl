package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"math/big"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// encodeBase62 encodes into a base62 string
func encodeBase62(num *big.Int) string {
	if num.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	var result string
	zero := big.NewInt(0)
	base := big.NewInt(62)

	for num.Cmp(zero) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, base, remainder)
		result = string(base62Chars[remainder.Int64()]) + result
	}
	slog.Debug("Encoded value into base62", "original", num, "new", result)
	return result
}

// generateShortCode generates a short code from a URL using hashing and base62 encoding
func generateShortCode(originalURL string) string {
	nonce, _ := generateNonce() // Nonce avoids getting same output for same input
	urlWithNonce := originalURL + nonce
	hash := sha256.Sum256([]byte(urlWithNonce))

	hashInt := new(big.Int)
	hashInt.SetBytes(hash[:])

	shortCode := encodeBase62(hashInt)

	if len(shortCode) > 6 {
		return shortCode[:6]
	}
	slog.Debug("Generated short code", "url", originalURL, "code", shortCode)
	return shortCode
}

// generateNonce generates a random 16 byte value
func generateNonce() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	slog.Debug("Generating nonce...")
	_, err := rand.Read(bytes)
	if err != nil {
		slog.Error("Error generating context", "error", err)
		return "", err
	}
	hex_val := hex.EncodeToString(bytes)
	slog.Debug("Generated nonce", "value", hex_val)
	return hex_val, nil
}
