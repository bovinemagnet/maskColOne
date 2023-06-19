package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	//key  = []byte("thisis32byteslongpassphrase12345") // Key size must be 16, 24, or 32 bytes long.
	key  []byte
	mode string
	in   string
	out  string
)

func init() {
	// Define and parse command-line arguments
	var keyString string
	flag.StringVar(&mode, "mode", "", "Mode for operation: --mode=e for encryption and --mode=d for decryption")
	flag.StringVar(&in, "in", "input.tsv", "Input file")
	flag.StringVar(&out, "out", "output.tsv", "Output file")
	flag.StringVar(&keyString, "key", "thisis32byteslongpassphrase12345", "Encryption Key (must be 32 characters long)")
	flag.Parse()

	// convert string to byte slice for key
	key = []byte(keyString)
}

func main() {
	switch mode {
	case "e":
		processFile(in, out, encrypt)
	case "d":
		processFile(in, out, decrypt)
	default:
		fmt.Println("Unknown mode. Use --mode=e for encryption or --mode=d for decryption.")
	}
}

func processFile(inputPath string, outputPath string, processFunc func(string, []byte) (string, error)) {
	// Read TSV file
	records := readTSV(inputPath)

	// Process records
	for i, record := range records {
		if len(record) > 0 {
			text := record[0]

			// Process (encrypt or decrypt)
			result, err := processFunc(text, key)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Replace first column with processed data
			record[0] = result
			records[i] = record
		}
	}

	// Write TSV file
	writeTSV(outputPath, records)
}

func readTSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Use tab-delimited instead of comma

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}

func writeTSV(path string, records [][]string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t' // Use tab-delimited instead of comma

	err = writer.WriteAll(records) // Write all records
	if err != nil {
		panic(err)
	}

	writer.Flush()
}

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	return fmt.Sprintf("%x", ciphertext), nil
}

func decrypt(ciphertext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
