package codec

import (
	"encoding"
	"reflect"
	"sync"

	"git.sr.ht/~whereswaldon/forest-go"
	"git.sr.ht/~whereswaldon/forest-go/serialize"
	"github.com/shamaton/msgpack"
)

// Codec can encode and decode forest types.
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(b []byte, v interface{}) error
}

// Default initializes to the safe Arbor codec, unless overriden by internal
// code such as test code.
type Default struct {
	Inner Codec
	sync.Once
}

func (d *Default) Encode(v interface{}) ([]byte, error) {
	d.Once.Do(func() {
		if d.Inner == nil {
			d.Inner = Arbor{}
		}
	})
	return d.Inner.Encode(v)
}
func (d *Default) Decode(b []byte, v interface{}) error {
	d.Once.Do(func() {
		if d.Inner == nil {
			d.Inner = Arbor{}
		}
	})
	return d.Inner.Decode(b, v)
}

// Arbor serializes forest types using the safe arbor serialization that
// validates the data.
type Arbor struct{}

func (ac Arbor) Encode(v interface{}) ([]byte, error) {
	buf, err := serialize.ArborSerialize(reflect.ValueOf(v))
	if err != nil {
		if b, ok := v.(encoding.BinaryMarshaler); ok {
			return b.MarshalBinary()
		}
		return nil, err
	}
	return buf, nil
}

func (ac Arbor) Decode(by []byte, v interface{}) error {
	copied := make([]byte, len(by))
	copy(copied, by)
	by = copied
	switch v := v.(type) {
	case *forest.Identity:
		return v.UnmarshalBinary(by)
	case *forest.Community:
		return v.UnmarshalBinary(by)
	case *forest.Reply:
		return v.UnmarshalBinary(by)
	}
	if _, err := serialize.ArborDeserialize(reflect.ValueOf(v), by); err != nil {
		if b, ok := v.(encoding.BinaryUnmarshaler); ok {
			return b.UnmarshalBinary(by)
		}
		return err
	}
	return nil
}

// Unsafe serializes forest types without validating them for the purposes of
// fast testing.
//
// This avoids branches within the orchard code that could lead to nodes
// not being validated.
//
// NOTE(jfm) msgpack is used because stdlib encodings have coincidental issues:
//
// - Gob calls into encoding.BinaryMarshaler method set, which calls the
// node validation code paths.
//
// - JSON encodes types like fields.Version into an escaped hexidecimal string,
// but refuses to decode back into fields.Version without some help.
//
// In an attempt to avoid touching the fields types, I've opted to go with a
// encoding that works "out of the box".
type Unsafe struct{}

func (Unsafe) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}
func (Unsafe) Decode(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}
