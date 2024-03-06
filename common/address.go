package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/goccy/go-json"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"
	"io"
	"trustmesh/encoders"
)

type NetworkAddress struct {
	AnonGeoLocation SafeLatitudeLongitude
	LocationKey     kyber.Point
	PrivateKey      kyber.Scalar
	PublicKey       kyber.Point
}

func NewNetworkAddress(lat, long float64) *NetworkAddress {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	private := suite.Scalar().Pick(random.New())
	public := suite.Point().Mul(private, nil)
	var sLatLon SafeLatitudeLongitude

	clat, clon := convertToPrecisionGrid(lat, long, 5)
	sLatLon.Set(float64(clat), float64(clon), 5)

	return &NetworkAddress{
		AnonGeoLocation: sLatLon,
		LocationKey:     CommitLocation(private, sLatLon.Bytes()),
		PrivateKey:      private,
		PublicKey:       public,
	}
}

func (na *NetworkAddress) PublicKeyBase64() string {
	publicKeyBytes, _ := na.PublicKey.MarshalBinary()
	return base64.StdEncoding.EncodeToString(publicKeyBytes)
}

// Encryption utility function
func encrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (na *NetworkAddress) EncodeToString(secretKey []byte) (string, error) {
	publicKeyBytes, err := na.PublicKey.MarshalBinary()
	if err != nil {
		return "", err
	}
	data := map[string]interface{}{
		"PublicKey": base64.StdEncoding.EncodeToString(publicKeyBytes),
		"Nonce":     GenerateOrUpdateNonce(na.PublicKeyBase64()),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return encrypt(jsonData, secretKey)
}

// GenerateSharedSecret takes another party's public key and generates the shared secret
func (na *NetworkAddress) GenerateSharedSecret(peerPublicKey kyber.Point) []byte {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	sharedSecret := suite.Point().Mul(na.PrivateKey, peerPublicKey)

	secretBytes, _ := sharedSecret.MarshalBinary()
	return secretBytes
}

func (na *NetworkAddress) Encode(input any) (any, error) {
	encoder := encoders.NewJsonEncoder()
	encoded, err := encoder.EncodeToBytes(nil, na)
	if err != nil {
		return nil, err
	}
	return encoded, nil

}

func (na *NetworkAddress) Decode(input any, output any) error {
	return nil
}

type Address interface {
	Network() *NetworkAddress
	Nonce() Nonce
	PubKey() []byte
	String() string
	Hex() string
	Bytes() []byte
	Len() int
	Version() int
	NetCode() byte
	LatLon() SafeLatitudeLongitude
}

type address struct{}
