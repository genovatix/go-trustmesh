package encoders

import (
	"go.dedis.ch/kyber/v3"
	"io"
)

// KyberPoint wraps kyber.Point to implement custom encoding/decoding
type KyberPoint struct {
	Point kyber.Point
}

// Implementing the BinaryMarshaler interface
func (kp *KyberPoint) MarshalBinary() ([]byte, error) {
	return kp.Point.MarshalBinary()
}

// Implementing the BinaryUnmarshaler interface
func (kp *KyberPoint) UnmarshalBinary(data []byte) error {
	return kp.Point.UnmarshalBinary(data)
}

// String returns a human-readable representation of the KyberPoint
func (kp *KyberPoint) String() string {
	return kp.Point.String()
}

// MarshalSize returns the encoded length of the KyberPoint
func (kp *KyberPoint) MarshalSize() int {
	return kp.Point.MarshalSize()
}

// MarshalTo encodes the KyberPoint and writes it to an io.Writer
func (kp *KyberPoint) MarshalTo(w io.Writer) (int, error) {
	data, err := kp.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

// UnmarshalFrom decodes the KyberPoint by reading from an io.Reader
func (kp *KyberPoint) UnmarshalFrom(r io.Reader) (int, error) {
	data := make([]byte, kp.MarshalSize())
	n, err := r.Read(data)
	if err != nil {
		return n, err
	}
	return n, kp.UnmarshalBinary(data)
}
