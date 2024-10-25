package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func GenerateRequestHash(request interface{}) (string, error) {
    requestJSON, err := json.Marshal(request)
    if err != nil {
        return "", err
    }
    hash := sha256.Sum256(requestJSON)
    return hex.EncodeToString(hash[:]), nil
}
