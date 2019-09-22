## Cipher library ##

### Example for using algorithm AES (GCM) ###

```go
package main

import (
	"bitbucket.org/sparkmaker/gohelper/cipher/aes"
	"bitbucket.org/sparkmaker/gohelper/logger/stdout"
)

func main() {
	// Secret key
	secret := "password"

	// New cipher algorithm AES
	cipher := aes.New(secret)

	// Encryption
	sampletext := "sample words"
	ciphertext := cipher.Encrypt(sampletext)
	stdout.Debug(ciphertext)

	// Decryption
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		stdout.Error(err.Error())
	}
	stdout.Debug(plaintext)
}
```