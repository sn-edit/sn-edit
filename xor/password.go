package xor

// EncryptDecrypt runs a XOR encryption on the input string, encrypting it if it hasn't already been,
// and decrypting it if it has, using the key provided.
// @Source: https://github.com/KyleBanks/XOREncryption/blob/master/Go/xor.go
func EncryptDecrypt(input, key string) (output string) {
	for i := range input {
		output += string(input[i] ^ key[i%len(key)])
	}
	return output
}
