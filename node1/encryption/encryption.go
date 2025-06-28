package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	// "fmt"
	// "bytes"
)

// EncryptFile encrypts the input file and writes the result to output file using the given key.
func GenrateKey() (string, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
func EncryptFile(key []byte, inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	_, err = io.Copy(writer, inFile)
	return err
}


func DecryptFile(key []byte, reader io.Reader) (*os.File, error) {
	// Create a temp file to write decrypted output
	tempFile, err := os.CreateTemp("", "decrypted_*")
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Read IV from the beginning of the encrypted stream
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(reader, iv); err != nil {
		return nil, err
	}

	// Setup stream decryption
	stream := cipher.NewCTR(block, iv)
	streamReader := &cipher.StreamReader{S: stream, R: reader}

	// Decrypt and write directly to temp file
	if _, err := io.Copy(tempFile, streamReader); err != nil {
		tempFile.Close()
		return nil, err
	}

	// Seek to start so frontend or next handler can read it
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		tempFile.Close()
		return nil, err
	}

	return tempFile, nil
}

// func DecryptBytes(key []byte, encryptedData []byte) ([]byte, error) {
// 	if len(encryptedData) < aes.BlockSize {
// 		return nil, fmt.Errorf("encrypted data too short")
// 	}

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	iv := encryptedData[:aes.BlockSize]
// 	ciphertext := encryptedData[aes.BlockSize:]

// 	stream := cipher.NewCTR(block, iv)
// 	reader := &cipher.StreamReader{
// 		S: stream,
// 		R: bytes.NewReader(ciphertext),
// 	}

// 	var decrypted bytes.Buffer
// 	_, err = io.Copy(&decrypted, reader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return decrypted.Bytes(), nil
// }


