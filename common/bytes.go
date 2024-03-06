package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

func BytesToPublicKey(bytes []byte) (*ecdsa.PublicKey, error) {
	if len(bytes) != 64 {
		return nil, fmt.Errorf("invalid public key length")
	}

	x := new(big.Int).SetBytes(bytes[:32])
	y := new(big.Int).SetBytes(bytes[32:])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}
