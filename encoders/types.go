package encoders

type Encodable interface {
	Encode(input any) (any, error)
	Decode(input any, output any) error
}

type EncodingType int

const (
	UNKNOWN EncodingType = iota
	JSON
	PROTO
	QMesh
)

type Wrapper struct {
	Item     interface{}
	Data     []byte
	Encoding EncodingType
	JSONEncoder
}

func (*Wrapper) Encode(input any) (any, error) {
	return nil, nil
}
func (*Wrapper) Decode(input any, output any) error {
	return nil
}

func NewWrapper(item interface{}, encoding EncodingType) Encodable {
	w := &Wrapper{Item: item, Data: nil, Encoding: encoding}
	return w
}
