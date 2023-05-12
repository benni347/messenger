package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	utils "github.com/benni347/messengerutils"
	"github.com/cloudflare/circl/sign/dilithium"
	kyberk2so "github.com/symbolicsoft/kyber-k2so"
	"golang.org/x/crypto/blake2b"
)

func GenerateECCKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func SignEcc(privateKey *ecdsa.PrivateKey, messageHash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, messageHash)
	if err != nil {
		return nil, err
	}

	curveBits := privateKey.PublicKey.Curve.Params().BitSize
	keyBytes := (curveBits + 7) / 8

	signature := make([]byte, keyBytes*2)
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	copy(signature[keyBytes-len(rBytes):], rBytes)
	copy(signature[keyBytes*2-len(sBytes):], sBytes)

	return signature, nil
}

func VerifyEcc(publicKey *ecdsa.PublicKey, messageHash, signature []byte) bool {
	curveBits := publicKey.Curve.Params().BitSize
	keyBytes := (curveBits + 7) / 8

	r := new(big.Int).SetBytes(signature[:keyBytes])
	s := new(big.Int).SetBytes(signature[keyBytes:])

	return ecdsa.Verify(publicKey, messageHash, r, s)
}

// From here to the lines which is a comment which contains --- the functions used are under the BSD3-Clause license.
// https://pkg.go.dev/github.com/cloudflare/circl/sign/dilithium

func GenerateDilithiumKeyPair(modeName string) (dilithium.PublicKey, dilithium.PrivateKey, error) {
	mode := dilithium.ModeByName(modeName)
	if mode == nil {
		return nil, nil, fmt.Errorf("mode not supported")
	}

	publicKey, privateKey, err := mode.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating key pair: %v", err)
	}

	return publicKey, privateKey, nil
}

func PackDilithiumKeys(
	publicKey dilithium.PublicKey,
	privateKey dilithium.PrivateKey,
) ([]byte, []byte) {
	return publicKey.Bytes(), privateKey.Bytes()
}

func UnpackDilithiumKeys(
	modeName string,
	packedPublicKey []byte,
	packedPrivateKey []byte,
) (dilithium.PublicKey, dilithium.PrivateKey) {
	mode := dilithium.ModeByName(modeName)

	return mode.PublicKeyFromBytes(packedPublicKey), mode.PrivateKeyFromBytes(packedPrivateKey)
}

func SignDilithium(
	privateKey dilithium.PrivateKey,
	msg []byte,
	modeName string,
) ([]byte, int, error) {
	mode := dilithium.ModeByName(modeName)
	if mode == nil {
		return nil, -1, fmt.Errorf("mode not supported")
	}

	signatureSize := mode.SignatureSize()

	return mode.Sign(privateKey, msg), signatureSize, nil
}

func VerifyDilithium(
	publicKey dilithium.PublicKey,
	msg []byte,
	signature []byte,
	modeName string,
) (bool, error) {
	mode := dilithium.ModeByName(modeName)
	if mode == nil {
		return false, fmt.Errorf("mode not supported")
	}

	return mode.Verify(publicKey, msg, signature), nil
}

// ---

func GenerateKyberKeyPair() ([kyberk2so.Kyber1024SKBytes]byte, [kyberk2so.Kyber1024PKBytes]byte, error) {
	privateKey, publicKey, err := kyberk2so.KemKeypair1024()
	if err != nil {
		utils.PrintError("During the Creation of the Kyber KeyPair", err)
		return privateKey, publicKey, err
	}

	return privateKey, publicKey, nil
}

func EncryptKyber(
	publicKey *[kyberk2so.Kyber1024PKBytes]byte,
) ([kyberk2so.Kyber1024CTBytes]byte, [kyberk2so.KyberSSBytes]byte, error) {
	return kyberk2so.KemEncrypt1024(*publicKey)
}

func DecryptKyber(
	ciphertext *[kyberk2so.Kyber1024CTBytes]byte,
	privateKey *[kyberk2so.Kyber1024SKBytes]byte,
) ([kyberk2so.KyberSSBytes]byte, error) {
	return kyberk2so.KemDecrypt1024(*ciphertext, *privateKey)
}

func CalculateHash(message []byte) []byte {
	hash, err := blake2b.New512(nil)
	if err != nil {
		utils.PrintError("During the Creation of the Hash", err)
		return nil
	}
	hash.Write(message)
	return hash.Sum(nil)
}
