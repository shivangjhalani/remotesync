package security

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "github.com/golang-jwt/jwt"
)

type Credentials struct {
    Username string
    Password string
}

type TokenManager struct {
    secretKey []byte
}

func NewTokenManager(secret string) *TokenManager {
    return &TokenManager{
        secretKey: []byte(secret),
    }
}

func (tm *TokenManager) GenerateToken(username string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    return token.SignedString(tm.secretKey)
}

func (tm *TokenManager) ValidateToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return tm.secretKey, nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["username"].(string), nil
    }

    return "", fmt.Errorf("invalid token")
}

func HashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}