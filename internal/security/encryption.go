package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

type Encryptor struct {
    gcm cipher.AEAD
}

func NewEncryptor(key []byte) (*Encryptor, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
    nonce := make([]byte, e.gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return e.gcm.Seal(nonce, nonce, data, nil), nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
    nonceSize := e.gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return e.gcm.Open(nil, nonce, ciphertext, nil)
}