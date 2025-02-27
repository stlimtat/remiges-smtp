package crypto

import "context"

const (
	KeyTypeRSA     = "rsa"
	KeyTypeEd25519 = "ed25519"
)

//go:generate mockgen -destination=mock.go -package=crypto . IKeyGenerator
type IKeyGenerator interface {
	GenerateKey(ctx context.Context, bitSize int, id string) (publicKeyPEM, privateKeyPEM []byte, err error)
}
