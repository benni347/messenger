# **Cryptographic Framework for Communication Platform**

This repository contains a bespoke cryptographic framework that is intended to be used for a communication platform. It provides an implementation of various cryptographic functions, including:

- Elliptic Curve Cryptography (ECC) for digital signature generation and verification
- Dilithium for digital signature generation and verification
- Kyber for key generation, encryption, and decryption
- Blake2b for hash calculation

The ECC and Kyber functions are built-in to Go's standard library, while the Dilithium and Blake2b functions are imported from external packages.

## Installation

To install this package, simply run:

```bash
go get github.com/benni347/encryption
```

## Usage

To use the functions provided in this package, first import the package:

```go
import (
    "github.com/benni347/encryption"
)
```

Then, call the relevant functions as needed. Here are some examples:

### Elliptic Curve Cryptography (ECC)

```go
privateKey, publicKey, err := encryption.GenerateECCKeyPair()
if err != nil {
    // Handle error
}

message := []byte("Hello, world!")
hash := encryption.CalculateHash(message)
signature, err := encryption.SignEcc(privateKey, hash)
if err != nil {
    // Handle error
}

isValid := encryption.VerifyEcc(publicKey, hash, signature)
```

### Dilithium

```go
modeName := "Dilithium2"
publicKey, privateKey, err := encryption.GenerateDilithiumKeyPair(modeName)
if err != nil {
    // Handle error
}

packedPublicKey, packedPrivateKey := encryption.PackDilithiumKeys(publicKey, privateKey)

// ...

signature, _, err := encryption.SignDilithium(privateKey, message, modeName)
if err != nil {
    // Handle error
}

isValid, err := encryption.VerifyDilithium(publicKey, message, signature, modeName)
if err != nil {
    // Handle error
}
```

### Kyber

```go
privateKey, publicKey, err := encryption.GenerateKyberKeyPair()
if err != nil {
    // Handle error
}

ciphertext, sharedSecret, err := encryption.EncryptKyber(&publicKey)
if err != nil {
    // Handle error
}

plaintext, err := encryption.DecryptKyber(&ciphertext, &privateKey)
if err != nil {
    // Handle error
}
```

### Blake2b

```go
message := []byte("Hello, world!")
hash := encryption.CalculateHash(message)
```

## API

### `encryption` package

#### Functions

##### Elliptic Curve Cryptography (ECC) Functions

- `GenerateECCKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)`: generates a private and public key pair for ECC.
- `SignEcc(privateKey *ecdsa.PrivateKey, messageHash []byte) ([]byte, error)`: signs a message hash using a private key.
- `VerifyEcc(publicKey *ecdsa.PublicKey, messageHash, signature []byte) bool`: verifies a message hash's signature using a public key.

##### Dilithium Functions

- `GenerateDilithiumKeyPair(modeName string) (dilithium.PublicKey, dilithium.PrivateKey, error)`: generates a private and public key pair for Dilithium, using a specified mode.
- `PackDilithiumKeys(publicKey dilithium.PublicKey, privateKey dilithium.PrivateKey) ([]byte, []byte)`: packs the public and private key into byte slices.
- `UnpackDilithiumKeys(modeName string, packedPublicKey []byte, packedPrivateKey []byte) (dilithium.PublicKey, dilithium.PrivateKey)`: unpacks the public and private key from byte slices.
- `SignDilithium(privateKey dilithium.PrivateKey, msg []byte, modeName string) ([]byte, int, error)`: signs a message using a private key and mode, returning the signature and its size.
- `VerifyDilithium(publicKey dilithium.PublicKey, msg []byte, signature []byte, modeName string) (bool, error)`: verifies a message's signature using a public key and mode.

##### Kyber Functions

- `GenerateKyberKeyPair() ([kyberk2so.Kyber1024SKBytes]byte, [kyberk2so.Kyber1024PKBytes]byte, error)`: generates a private and public key pair for Kyber.
- `EncryptKyber(publicKey *[kyberk2so.Kyber1024PKBytes]byte) ([kyberk2so.Kyber1024CTBytes]byte, [kyberk2so.KyberSSBytes]byte, error)`: encrypts a message using a public key.
- `DecryptKyber(ciphertext *[kyberk2so.Kyber1024CTBytes]byte, privateKey *[kyberk2so.Kyber1024SKBytes]byte) ([kyberk2so.KyberSSBytes]byte, error)`: decrypts a ciphertext using a private key.

##### Blake2b Functions

- `CalculateHash(message []byte) []byte`: calculates the hash of a given message using Blake2b.

## License

The ECC and Kyber functions used in this repository are part of the Go standard library and are licensed under a BSD-style license.

The Dilithium functions used in this repository are imported from the github.com/cloudflare/circl/sign/dilithium package and are licensed under a BSD 3-clause "New" or "Revised" License.

The Blake2b function used in this repository is imported from the golang.org/x/crypto/blake2b package and is licensed under a BSD-style license.

The rest of the code in this repository is licensed under the MIT License.

## Contributing

Contributions to this repository are welcome. If you have any suggestions, bug reports, or feature requests, please open an issue on GitHub.

## Acknowledgments

This repository was created as part of a culminating academic endeavor.

### Contributors

<a href="https://github.com/benni347/encryption/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=benni347/encryption" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
