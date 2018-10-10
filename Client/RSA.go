package main

import (
	//"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {

	// Generate RSA Key

	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	raulPublicKey := &raulPrivateKey.PublicKey


	fmt.Println("Private Key : ", raulPrivateKey)
	fmt.Println("Public key ", raulPublicKey)

	//Encrypt Miryan Message
	message := []byte("the code must be like a piece of music")
	label := []byte("")
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, message, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("The message encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, raulPrivateKey, ciphertext, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("The message decrypted [%x] to \n[%s]\n", ciphertext, plainText)

}