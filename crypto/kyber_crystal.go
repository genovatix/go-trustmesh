package crypto

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

// PublicKey Define PublicKey and PrivateKey types for better type safety.
type PublicKey kyber.Point
type PrivateKey kyber.Scalar

type KyberCrystal struct{}

const KYBER_VERSION = "3"

// GenerateKeys generates a public-private key pair using Kyber library.
func (k KyberCrystal) GenerateKeys() (PublicKey, PrivateKey, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	private := suite.Scalar().Pick(suite.RandomStream()) // Generate the private key.
	public := suite.Point().Mul(private, nil)            // Derive the public key from the private key.

	// Casting to defined types for better type safety.
	return PublicKey(public), PrivateKey(private), nil
}

// Usage returns the cryptographic capabilities of the KyberCrystal implementation.
func (k KyberCrystal) Usage() []AlgUsage {
	return []AlgUsage{KEX, KEM, ENCRYPTION}
}

// Name returns the name of the cryptographic implementation.
func (k KyberCrystal) Name() string {
	return "kyber-crystal"
}

// Version returns the version of the Kyber library used.
func (k KyberCrystal) Version() string {
	return KYBER_VERSION
}

// GetKyberAlgorithm returns an instance of the KyberCrystal algorithm.
func GetKyberAlgorithm() KyberCrystal {
	return KyberCrystal{}
}
