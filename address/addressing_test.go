package address

import (
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"testing"
)

func Test(t *testing.T) {
	t.Log("start testing")
	ctx := signing.CreateContext("secp256k1")
	key := ctx.NewRandomPrivateKey()

	t.Logf("key:" + key.AsHex())
	signer := signing.NewCryptoFactory(ctx).NewSigner(key)
	address := MakeAccountAddress(signer.GetPublicKey().AsHex())
	t.Log("address:" + address)
}
