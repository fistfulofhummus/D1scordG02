package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

// RemoveDuplicatesValues: A helper function to remove duplicate items in a list
func RemoveDuplicatesValues(arrayToEdit []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range arrayToEdit {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// https://stackoverflow.com/questions/28828440/is-there-a-way-to-write-generic-code-to-find-out-whether-a-slice-contains-specif
func Find(slice, elem interface{}) bool {
	sv := reflect.ValueOf(slice)

	// Check that slice is actually a slice/array.
	// you might want to return an error here
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		return false
	}

	// iterate the slice
	for i := 0; i < sv.Len(); i++ {

		// compare elem to the current slice element
		if elem == sv.Index(i).Interface() {
			return true
		}
	}

	// nothing found
	return false
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func SendGET(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return -1
	}
	return resp.StatusCode
}

// Will get around to this later
// func sendPOST(url string, payload string) int {
// 	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload))
// 	if err != nil {
// 		return 2
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		return 0
// 	}
// 	return 1
// }

func UpdateStats([]int) {

}

// func EncrFile(filepath string, aesKey string) (string, error) {
// 	plaintext, err := ioutil.ReadFile(filepath)
// 	aesKeyUnstring, _ := hex.DecodeString(aesKey)
// 	if err != nil {
// 		return "0", err
// 	}
// 	block, err := aes.NewCipher(aesKeyUnstring)
// 	if err != nil {
// 		return "0", err
// 	}
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return "0", err
// 	}
// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		return "0", err
// 	}
// 	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
// 	os.WriteFile("New.txt", ciphertext, 0644)
// 	encrFilePath := "New.txt"
// 	return encrFilePath, err
// }

func EncrFile(filepath string, aesKey string) (string, error) {
	// Read the plaintext from the file
	plaintext, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// Decode the AES key from hex string to []byte
	aesKeyUnstring, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(aesKeyUnstring)
	if err != nil {
		return "", err
	}

	// Create a GCM (Galois/Counter Mode) cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce (unique value) for encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext using GCM
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	encFilePath := filepath + ".enc"
	// Write the ciphertext to a new file
	err = os.WriteFile(encFilePath, ciphertext, 0644)
	if err != nil {
		return "", err
	}
	// Return the path of the encrypted file
	return encFilePath, nil
}
