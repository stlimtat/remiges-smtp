package crypto

import "context"

const (
	KeyTypeEd25519 = "ed25519"
	KeyTypeRSA     = "rsa"
)

//go:generate mockgen -destination=mock.go -package=crypto . IKeyGenerator,IKeyWriter
type IKeyGenerator interface {
	GenerateKey(ctx context.Context, bitSize int, id string) (publicKeyPEM, privateKeyPEM []byte, err error)
}

type IKeyWriter interface {
	WriteKey(ctx context.Context, id string, publicKeyPEM, privateKeyPEM []byte) error
}
