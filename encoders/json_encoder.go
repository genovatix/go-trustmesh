package encoders

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-json"
	"io"
)

type JSONEncoder struct {
	encoder json.Marshaler
	decoder json.Unmarshaler
}

func (j *JSONEncoder) EncodeToIO(ctx context.Context, e Encodable, w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.EncodeContext(ctx, e)
	if err != nil {
		errStr := fmt.Sprintf("encoder error (EncodeToIO): %s", ERROR_JSON_ENCODE.Error())
		return errors.New(errStr)
	}
	return nil
}

func (j *JSONEncoder) EncodeToBytes(ctx context.Context, e Encodable) ([]byte, error) {
	data, err := json.Marshal(e)

	if err != nil {
		errStr := fmt.Sprintf("encoder error (EncodeToBytes): %s", ERROR_JSON_ENCODE.Error())
		return nil, errors.New(errStr)
	}
	return data, nil
}

func NewJsonEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

func DumpJson(data interface{}) {
	spew.Dump(data)
}

func (j *JSONEncoder) DecodeFromBytes(e []byte, dest Encodable) error {
	err := json.Unmarshal(e, &dest)

	if err != nil {
		return err
	}

	return nil

}
