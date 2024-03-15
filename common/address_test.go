package common_test

import (
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"testing"
	"trustmesh/common" // Adjust this to the actual path.
)

func TestNetworkAddressSerializationDeserialization(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Creating a NetworkAddress
	lat, lon := 37.7749, -122.4194
	networkAddress, err := common.NewNetworkAddress(lat, lon)
	require.NoError(t, err, "Creating NetworkAddress should not fail.")

	// Manually set the suite for testing purposes
	networkAddress.Suite = suite

	// Generating mock ZKP for testing; assume this function works correctly.
	err = networkAddress.GenerateZKP(256)
	require.NoError(t, err, "Generating ZKP should not fail.")

	// Marshaling NetworkAddress to JSON
	jsonData, err := networkAddress.MarshalJSON()
	require.NoError(t, err, "Marshaling NetworkAddress to JSON should not fail.")

	// Unmarshaling JSON back to NetworkAddress
	var decodedAddress common.NetworkAddress
	decodedAddress.Suite = suite // Setting suite to decode the LocationCommitment
	err = decodedAddress.UnmarshalJSON(jsonData)
	require.NoError(t, err, "Unmarshaling JSON to NetworkAddress should not fail.")

	// Comparing original and decoded LocationCommitment
	originalCommitmentBytes, err := networkAddress.LocationCommitment.MarshalBinary()
	require.NoError(t, err, "Marshaling original LocationCommitment should not fail.")

	decodedCommitmentBytes, err := decodedAddress.LocationCommitment.MarshalBinary()
	require.NoError(t, err, "Marshaling decoded LocationCommitment should not fail.")

	// Using base64 encoding for comparison due to potential +/= in base64 output
	require.Equal(t, base64.StdEncoding.EncodeToString(originalCommitmentBytes), base64.StdEncoding.EncodeToString(decodedCommitmentBytes), "Original and decoded LocationCommitment should match.")

	// Extend this test to check other fields as necessary.
}
