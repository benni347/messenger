package fields

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	// CurrentVersion is the Forest version that this library writes
	CurrentVersion Version = 1

	// HashDigestLengthSHA512_256 is the length of the digest produced by the SHA512/256 hash algorithm
	HashDigestLengthSHA512_256 ContentLength = 32
)

// MultiByteSerializationOrder defines the order in which multi-byte
// integers are serialized into binary
var MultiByteSerializationOrder binary.ByteOrder = binary.BigEndian

// fundamental types
type genericType uint8

const sizeofgenericType = 1

func (g genericType) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, MultiByteSerializationOrder, g)
	return b.Bytes(), err
}

func (g *genericType) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, MultiByteSerializationOrder, g)
}

func (g *genericType) BytesConsumed() int {
	return sizeofgenericType
}

func (g *genericType) Equals(g2 *genericType) bool {
	if g == nil {
		return g2 == nil
	}
	return *g == *g2
}

// ContentLength represents the length of a piece of data in the Forest
type ContentLength uint16

const sizeofContentLength = 2

const (
	// MaxContentLength is the maximum representable content length in this
	// version of the Forest
	MaxContentLength = math.MaxUint16
)

func NewContentLength(size int) (*ContentLength, error) {
	if size > MaxContentLength {
		return nil, fmt.Errorf("Cannot represent content of size %d, max is %d", size, MaxContentLength)
	}
	c := ContentLength(size)
	return &c, nil
}

// MarshalBinary converts the ContentLength into its binary representation
func (c ContentLength) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, MultiByteSerializationOrder, c)
	return b.Bytes(), err
}

const contentLengthTextFormat = "B%d"

func (c ContentLength) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(contentLengthTextFormat, c)), nil
}

func (c *ContentLength) UnmarshalText(b []byte) error {
	_, err := fmt.Sscanf(string(b), contentLengthTextFormat, c)
	if err != nil {
		return fmt.Errorf("failed unmarshalling content length: %v", err)
	}
	return nil
}

// UnmarshalBinary converts from the binary representation of a ContentLength
// back to its structured form
func (c *ContentLength) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, MultiByteSerializationOrder, c)
}

func (c *ContentLength) BytesConsumed() int {
	return sizeofContentLength
}

func (c *ContentLength) Equals(c2 *ContentLength) bool {
	if c == nil {
		return c2 == nil
	}
	return *c == *c2
}

// TreeDepth represents the depth of a node within a tree
type TreeDepth uint32

const sizeofTreeDepth = 4

// MarshalBinary converts the TreeDepth into its binary representation
func (t TreeDepth) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, MultiByteSerializationOrder, t)
	return b.Bytes(), err
}

func (t TreeDepth) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("L%d", t)), nil
}

// UnmarshalBinary converts from the binary representation of a TreeDepth
// back to its structured form
func (t *TreeDepth) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, MultiByteSerializationOrder, t)
}

func (t *TreeDepth) BytesConsumed() int {
	return sizeofTreeDepth
}

func (t *TreeDepth) Equals(t2 *TreeDepth) bool {
	if t == nil {
		return t2 == nil
	}
	return *t == *t2
}

// Blob represents a quantity of arbitrary binary data in the Forest
type Blob []byte

// Contains checks if a substing of bytes exists in a Blob
func (v Blob) Contains(c []byte) bool {
	return strings.Contains(string(v), string(c))
}

// ContainsString checks if a substing exists in a Blob
func (v Blob) ContainsString(s string) bool {
	return v.Contains([]byte(s))
}

// MarshalBinary converts the Blob into its binary representation
func (v Blob) MarshalBinary() ([]byte, error) {
	return v, nil
}

func (v Blob) MarshalText() ([]byte, error) {
	based := base64.RawURLEncoding.EncodeToString([]byte(v))
	return []byte(based), nil
}

func (v *Blob) UnmarshalText(b []byte) error {
	if []byte(*v) == nil {
		*v = make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	}
	_, err := base64.RawURLEncoding.Decode(*v, b)
	if err != nil {
		return fmt.Errorf("failed decoding blob: %v", err)
	}
	return nil
}

// UnmarshalBinary converts from the binary representation of a Blob
// back to its structured form
func (v *Blob) UnmarshalBinary(b []byte) error {
	*v = b
	return nil
}

func (v *Blob) BytesConsumed() int {
	return len([]byte(*v))
}

func (v *Blob) Equals(v2 *Blob) bool {
	if v == nil {
		return v2 == nil
	}
	return bytes.Equal([]byte(*v), []byte(*v2))
}

// Version represents the version of the Arbor Forest Schema used to construct
// a particular node
type Version uint16

const sizeofVersion = 2

// MarshalBinary converts the Version into its binary representation
func (v Version) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, MultiByteSerializationOrder, v)
	return b.Bytes(), err
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("V%d", v)), nil
}

// UnmarshalBinary converts from the binary representation of a Version
// back to its structured form
func (v *Version) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, MultiByteSerializationOrder, v)
}

func (v *Version) BytesConsumed() int {
	return sizeofVersion
}

func (v *Version) Equals(v2 *Version) bool {
	if v == nil {
		return v2 == nil
	}
	return *v == *v2
}

// Timestamp represents the time at which a node was created. It is measured as milliseconds
// since the start of the UNIX epoch.
type Timestamp uint64

const sizeofTimestamp = 8

const nanosPerMilli = 1000000

func TimestampFrom(t time.Time) Timestamp {
	return Timestamp(t.UnixNano() / nanosPerMilli)
}

func (t Timestamp) Time() time.Time {
	sec := (uint(t) / 1000)
	nsec := (uint(t) % 1000) * nanosPerMilli
	return time.Unix(int64(sec), int64(nsec))
}

// MarshalBinary converts the Timestamp into its binary representation
func (v Timestamp) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, MultiByteSerializationOrder, v)
	return b.Bytes(), err
}

func (v Timestamp) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("V%d", v)), nil
}

// UnmarshalBinary converts from the binary representation of a Timestamp
// back to its structured form
func (v *Timestamp) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, MultiByteSerializationOrder, v)
}

func (v *Timestamp) BytesConsumed() int {
	return sizeofTimestamp
}

func (v *Timestamp) Equals(v2 *Timestamp) bool {
	if v == nil {
		return v2 == nil
	}
	return *v == *v2
}

// specialized types
type NodeType genericType

const (
	NodeTypeIdentity NodeType = iota
	NodeTypeCommunity
	NodeTypeReply

	sizeofNodeType = sizeofgenericType
)

var ValidNodeTypes = map[NodeType]struct{}{
	NodeTypeIdentity:  struct{}{},
	NodeTypeCommunity: struct{}{},
	NodeTypeReply:     struct{}{},
}

var NodeTypeNames = map[NodeType]string{
	NodeTypeIdentity:  "identity",
	NodeTypeCommunity: "community",
	NodeTypeReply:     "reply",
}

func (t NodeType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t NodeType) MarshalText() ([]byte, error) {
	return []byte(NodeTypeNames[t]), nil
}

func (t *NodeType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidNodeTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid node type", *t)
	}
	return nil
}

func (t *NodeType) BytesConsumed() int {
	return sizeofNodeType
}

func (t *NodeType) Equals(t2 *NodeType) bool {
	if t == nil {
		return t2 == nil
	}
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

func (t NodeType) String() string {
	return NodeTypeNames[t]
}

type HashType genericType

const (
	HashTypeNullHash HashType = iota
	HashTypeSHA512

	sizeofHashType = sizeofgenericType
)

// map to valid lengths
var ValidHashTypes = map[HashType][]ContentLength{
	HashTypeNullHash: []ContentLength{0},
	HashTypeSHA512:   []ContentLength{HashDigestLengthSHA512_256},
}

var HashNames = map[HashType]string{
	HashTypeNullHash: "NullHash",
	HashTypeSHA512:   "SHA512",
}

func (t HashType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t HashType) MarshalText() ([]byte, error) {
	return []byte(HashNames[t]), nil
}

func (t *HashType) UnmarshalText(b []byte) error {
	for hashType, hashName := range HashNames {
		if hashName == string(b) {
			*t = hashType
			return nil
		}
	}
	return fmt.Errorf("no such hash type %s", string(b))
}

func (t *HashType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidHashTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid hash type", *t)
	}
	return nil
}

func (t *HashType) BytesConsumed() int {
	return sizeofHashType
}

func (t *HashType) Equals(t2 *HashType) bool {
	if t == nil {
		return t2 == nil
	}
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type ContentType genericType

const (
	sizeofContentType                 = sizeofgenericType
	ContentTypeUTF8String ContentType = 1
	ContentTypeTwig       ContentType = 2
)

var ValidContentTypes = map[ContentType]struct{}{
	ContentTypeUTF8String: struct{}{},
	ContentTypeTwig:       struct{}{},
}

var ContentNames = map[ContentType]string{
	ContentTypeUTF8String: "UTF-8",
	ContentTypeTwig:       "Twig",
}

func (t ContentType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t ContentType) MarshalText() ([]byte, error) {
	return []byte(ContentNames[t]), nil
}

func (t *ContentType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidContentTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid content type", *t)
	}
	return nil
}

func (t *ContentType) BytesConsumed() int {
	return sizeofContentType
}

func (t *ContentType) Equals(t2 *ContentType) bool {
	if t == nil {
		return t2 == nil
	}
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type KeyType genericType

const (
	sizeofKeyType             = sizeofgenericType
	KeyTypeNoKey      KeyType = 0
	KeyTypeOpenPGPRSA KeyType = 1
)

var ValidKeyTypes = map[KeyType]struct{}{
	KeyTypeNoKey:      struct{}{},
	KeyTypeOpenPGPRSA: struct{}{},
}

var KeyNames = map[KeyType]string{
	KeyTypeNoKey:      "None",
	KeyTypeOpenPGPRSA: "OpenPGP-RSA",
}

func (t KeyType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t KeyType) MarshalText() ([]byte, error) {
	return []byte(KeyNames[t]), nil
}

func (t *KeyType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidKeyTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid key type", *t)
	}
	return nil
}

func (t *KeyType) BytesConsumed() int {
	return sizeofKeyType
}

func (t *KeyType) Equals(t2 *KeyType) bool {
	if t == nil {
		return t2 == nil
	}
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type SignatureType genericType

const (
	sizeofSignatureType                   = sizeofgenericType
	SignatureTypeOpenPGPRSA SignatureType = 1
)

var ValidSignatureTypes = map[SignatureType]struct{}{
	SignatureTypeOpenPGPRSA: struct{}{},
}

var SignatureNames = map[SignatureType]string{
	SignatureTypeOpenPGPRSA: "OpenPGP-RSA",
}

func (t SignatureType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t SignatureType) MarshalText() ([]byte, error) {
	return []byte(SignatureNames[t]), nil
}

func (t *SignatureType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidSignatureTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid signature type", *t)
	}
	return nil
}

func (t *SignatureType) BytesConsumed() int {
	return sizeofSignatureType
}

func (t *SignatureType) Equals(t2 *SignatureType) bool {
	if t == nil {
		return t2 == nil
	}
	return ((*genericType)(t)).Equals((*genericType)(t2))
}
