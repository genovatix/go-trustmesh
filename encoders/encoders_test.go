package encoders_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"testing"
	"trustmesh/encoders"
)

func TestKyberPointEncodingDecoding(t *testing.T) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	point := suite.Point().Base()

	kp := &encoders.KyberPoint{Point: point}

	// Test Binary Marshal/Unmarshal
	data, err := kp.MarshalBinary()
	require.NoError(t, err, "Binary marshal should not fail")

	var kpUnmarshaled encoders.KyberPoint
	kpUnmarshaled.Point = suite.Point() // Initialize the point
	err = kpUnmarshaled.UnmarshalBinary(data)
	require.NoError(t, err, "Binary unmarshal should not fail")

	require.Equal(t, kp.String(), kpUnmarshaled.String(), "Unmarshaled point should match original")

	// Test MarshalTo/UnmarshalFrom
	var buf bytes.Buffer
	n, err := kp.MarshalTo(&buf)
	require.NoError(t, err, "MarshalTo should not fail")
	require.Equal(t, kp.MarshalSize(), n, "Written bytes should match marshal size")

	var kpFrom encoders.KyberPoint
	kpFrom.Point = suite.Point() // Initialize the point
	n, err = kpFrom.UnmarshalFrom(&buf)
	require.NoError(t, err, "UnmarshalFrom should not fail")
	require.Equal(t, kp.String(), kpFrom.String(), "Unmarshaled point from io.Reader should match original")
}
