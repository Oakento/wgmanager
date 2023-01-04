package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func newPrivateKey() [32]byte {
	var k [32]byte
	_, err := rand.Read(k[:])
	if err != nil {
		fmt.Println(err)
	}
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return k
}

func newPublicKey(prikey [32]byte) [32]byte {
	var k [32]byte
	curve25519.ScalarBaseMult(&k, &prikey)
	return k
}

func NewWgKeyPairs() (publicKey string, privateKey string) {
	prikey := newPrivateKey()
	pubkey := newPublicKey(prikey)
	return base64.StdEncoding.EncodeToString(pubkey[:]), base64.StdEncoding.EncodeToString(prikey[:])
}
