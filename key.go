package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/curve25519"
)

const (
	WireguardPublicKeySize  = 32 // The size of a Wireguard public key in bytes.
	WireguardPrivateKeySize = 32 // The size of a Wireguard private key in bytes.
)

type (
	WireguardPublicKey  [WireguardPublicKeySize]byte  // A Wireguard public key.
	WireguardPrivateKey [WireguardPrivateKeySize]byte // A Wireguard private key.
)

// clamp clamps the private key to a valid curve25519 scalar.
func (sk *WireguardPrivateKey) clamp() {
	sk[0] &= 248
	sk[31] = (sk[31] & 127) | 64
}

// newWireguardPrivateKey generates a new random private key and clamps it.
func newWireguardPrivateKey() (sk WireguardPrivateKey, err error) {
	_, err = rand.Read(sk[:])
	sk.clamp()
	return
}

// publicKey derives the corresponding public key from the private key.
func (sk *WireguardPrivateKey) publicKey() (pk WireguardPublicKey) {
	apk := (*[WireguardPublicKeySize]byte)(&pk)
	ask := (*[WireguardPrivateKeySize]byte)(sk)
	curve25519.ScalarBaseMult(apk, ask)
	return
}

// base64PrivateKey returns the private key encoded in base64 format.
func (sk *WireguardPrivateKey) base64PrivateKey() (pks string) {
	pks = base64.StdEncoding.EncodeToString(sk[:])
	return
}

// base64PublicKey returns the public key encoded in base64 format.
func (sk *WireguardPrivateKey) base64PublicKey() (pks string) {
	pk := sk.publicKey()
	pks = base64.StdEncoding.EncodeToString(pk[:])
	return
}
