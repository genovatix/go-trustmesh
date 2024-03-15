package common

import (
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"
	"math/big"
	"trustmesh/libzk13"
)

// SafeLatitudeLongitude represents an anonymized geographical location.
type SafeLatitudeLongitude []int

// NetworkAddress includes cryptographic elements and an anonymized location.
type NetworkAddress struct {
	AnonGeoLocation    SafeLatitudeLongitude
	LocationCommitment kyber.Point   `json:"locationCommitment"`
	ZKP                *libzk13.ZK13 `json:"-"`
	PrivateKey         kyber.Scalar  `json:"-"`
	PublicKey          kyber.Point   `json:"public_key"`
	r, P               *big.Int
	Suite              kyber.Group
}

// AddressInfo provides a serializable and usable representation of NetworkAddress.
type AddressInfo struct {
	PublicKey          string `json:"publicKey"`
	LocationCommitment string `json:"locationCommitment"`
	ZKPProof           string `json:"zkpProof"`
}

// GenerateCryptoKeys creates a pair of cryptographic keys using the Kyber library.
func GenerateCryptoKeys() (kyber.Group, kyber.Scalar, kyber.Point, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey := suite.Scalar().Pick(random.New())
	publicKey := suite.Point().Mul(privateKey, nil)
	return suite, privateKey, publicKey, nil
}

// NewNetworkAddress initializes a NetworkAddress with given latitude and longitude.
func NewNetworkAddress(lat, lon float64) (*NetworkAddress, error) {
	// Assume ConvertToPrecisionGrid and CommitLocation functions are defined in latlon.go

	suite, privateKey, publicKey, err := GenerateCryptoKeys()

	if err != nil {
		return nil, fmt.Errorf("error generating crypto keys: %v", err)
	}

	precision, err := GetDynamicPrecision()
	anonGeoLocation, err := ConvertToPrecisionGrid(lat, lon, precision)
	if err != nil {
		return nil, fmt.Errorf("error converting to precision grid: %v", err)
	}

	anonGeoBytes, err := anonGeoLocation.Bytes()
	_, locationCommitment, err := CommitLocation(privateKey, anonGeoBytes)
	if err != nil {
		return nil, fmt.Errorf("error creating location commitment: %v", err)
	}

	na := &NetworkAddress{
		AnonGeoLocation:    anonGeoLocation,
		LocationCommitment: locationCommitment,
		PrivateKey:         privateKey,
		PublicKey:          publicKey,
		Suite:              suite,
	}

	return na, nil
}

// GenerateZKP generates a Zero-Knowledge Proof for the NetworkAddress.
func (na *NetworkAddress) GenerateZKP(bits int) error {
	if na.AnonGeoLocation == nil || len(na.AnonGeoLocation) == 0 {
		return fmt.Errorf("AnonGeoLocation is empty. Cannot generate ZKP")
	}

	secretBaggage := fmt.Sprintf("%v", na.AnonGeoLocation)
	na.ZKP = libzk13.NewZK13(secretBaggage, bits)
	r, P := na.ZKP.Prover() // Assume Prover method exists and returns big.Int values
	na.r = r
	na.P = P

	return nil
}

// GenerateAddress creates a new NetworkAddress and encapsulates it into AddressInfo.
func GenerateAddress(lat, lon float64, bits int) (*AddressInfo, error) {
	na, err := NewNetworkAddress(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to create network address: %v", err)
	}

	err = na.GenerateZKP(bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ZKP: %v", err)
	}

	// Serialize public key, location commitment, and ZKP for usability
	publicKeyStr, err := na.PublicKey.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize public key: %v", err)
	}
	locationCommitmentStr, err := na.LocationCommitment.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize location commitment: %v", err)
	}
	zkpProofStr := fmt.Sprintf("%s|%s", na.r.Text(16), na.P.Text(16))

	addressInfo := &AddressInfo{
		PublicKey:          string(publicKeyStr),
		LocationCommitment: string(locationCommitmentStr),
		ZKPProof:           zkpProofStr,
	}

	return addressInfo, nil
}

func (na *NetworkAddress) MarshalJSON() ([]byte, error) {
	commitmentBytes, err := na.LocationCommitment.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LocationCommitment: %v", err)
	}
	locationCommitmentStr := base64.StdEncoding.EncodeToString(commitmentBytes)

	// Serialize the rest of NetworkAddress, converting LocationCommitment to a base64 string
	return json.Marshal(&struct {
		LocationCommitment string `json:"locationCommitment"`
	}{
		LocationCommitment: locationCommitmentStr,
	})
}

// UnmarshalJSON customizes the JSON unmarshaling for NetworkAddress.
func (na *NetworkAddress) UnmarshalJSON(data []byte) error {
	// Temporary struct to extract LocationCommitment as a base64-encoded string
	temp := struct {
		LocationCommitment string `json:"locationCommitment"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal NetworkAddress: %v", err)
	}

	// Decode the base64-encoded LocationCommitment string to a byte slice
	commitmentBytes, err := base64.StdEncoding.DecodeString(temp.LocationCommitment)
	if err != nil {
		return fmt.Errorf("failed to decode LocationCommitment from base64: %v", err)
	}

	// Convert bytes back to kyber.Point for LocationCommitment
	commitment := na.Suite.Point() // Ensure the suite is properly initialized
	if err := commitment.UnmarshalBinary(commitmentBytes); err != nil {
		return fmt.Errorf("failed to unmarshal LocationCommitment to kyber.Point: %v", err)
	}
	na.LocationCommitment = commitment

	// Unmarshal the rest of the struct as usual
	// You might need to handle other fields similarly, especially if they use kyber types
	return nil
}
