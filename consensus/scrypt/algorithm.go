package scrypt

import (
	"gbchain-org/go-gbchain/crypto"
	"gbchain-org/go-gbchain/crypto/scrypt"
)

func ScryptHash(hash []byte, nonce uint64) ([]byte, []byte) {
	hashT := make([]byte, 80)
	copy(hashT[0:32], hash[:])
	copy(hashT[32:64], hash[:])
	copy(hashT[72:], []byte{
		byte(nonce >> 56),
		byte(nonce >> 48),
		byte(nonce >> 40),
		byte(nonce >> 32),
		byte(nonce >> 24),
		byte(nonce >> 16),
		byte(nonce >> 8),
		byte(nonce),
	})

	if digest, err := scrypt.Key(hashT, hashT, 1024, 1, 1, 32, ScryptMode); err == nil {
		return crypto.Keccak256(digest), digest
	} else {
		panic(err.Error())
	}
}
