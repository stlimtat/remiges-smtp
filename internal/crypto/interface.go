package crypto

import (
	"context"
	"crypto"
)

const (
	KeyTypeEd25519 = "ed25519"
	KeyTypeRSA     = "rsa"
)

//go:generate mockgen -destination=mock.go -package=crypto . IKeyGenerator,IKeyLoader,IKeyWriter
type IKeyGenerator interface {
	GenerateKey(ctx context.Context, bitSize int, id string, keyType string) (publicKeyPEM, privateKeyPEM []byte, err error)
}

type IKeyLoader interface {
	LoadPrivateKey(ctx context.Context, keyType string, privateKeyPath string) (crypto.Signer, error)
}

type IKeyWriter interface {
	WriteKey(ctx context.Context, id string, publicKeyPEM, privateKeyPEM []byte) (publicKeyPath, privateKeyPath string, err error)
}
