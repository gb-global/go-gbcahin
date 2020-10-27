package core

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/crypto"
	"gbchain-org/go-gbchain/rlp"
)

var (
	emptyCtx = NewCrossTransaction(
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(1024),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		common.Address{},
		nil,
	)

	rightvrsCtx, _ = NewCrossTransaction(
		big.NewInt(1e18),
		big.NewInt(2e18),
		big.NewInt(1024),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToHash("0b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca"),
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		common.Address{},
		nil,
	).WithSignature(
		NewEIP155CtxSigner(big.NewInt(1)),
		common.Hex2Bytes("fff9e65e751407a69c5125a0e0dafd2e0048ce9b60d39bb0b58c251b4a72d3822005750b091faae17b20a0a966b0c40ca44134fb214d0df0fed1c10646141f7000"),
	)
)

func TestCrossTransactionSigHash(t *testing.T) {
	signer := NewEIP155CtxSigner(big.NewInt(1))
	if signer.Hash(emptyCtx) != common.HexToHash("4f4a4cb2a67de52c0036f3ab5046388144471c13d42d62e2390e0c6eaa7e3d50") {
		t.Errorf("empty transaction hash mismatch, got %x", emptyCtx.Hash())
	}
	if signer.Hash(rightvrsCtx) != common.HexToHash("ba84e225625b9b69304e867a2c9af871b37086b7238190dd44dbd57e5429304e") {
		t.Errorf("RightVRS transaction hash mismatch, got %x", rightvrsCtx.Hash())
	}
}

func TestCrossTransactionEncode(t *testing.T) {
	ctxb, err := rlp.EncodeToBytes(rightvrsCtx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := common.FromHex("f8e8f8e6880de0b6b3a7640000a00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbcaa00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca94095e7baea6a6c7c4c2dfeb977efac326af552d87940000000000000000000000000000000000000000a00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca820400881bc16d674ec800008025a0fff9e65e751407a69c5125a0e0dafd2e0048ce9b60d39bb0b58c251b4a72d382a02005750b091faae17b20a0a966b0c40ca44134fb214d0df0fed1c10646141f70")
	if !bytes.Equal(ctxb, should) {
		t.Errorf("encoded RLP mismatch, got %x", ctxb)
	}
}

func decodeCtx(data []byte) (*CrossTransaction, error) {
	var tx CrossTransaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

func TestCtxRecipient(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeCtx(common.Hex2Bytes("f8e8f8e6880de0b6b3a7640000a00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbcaa00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca94095e7baea6a6c7c4c2dfeb977efac326af552d87940000000000000000000000000000000000000000a00b2aa4c82a3b0187a087e030a26b71fc1a49e74d3776ae8e03876ea9153abbca820400881bc16d674ec800008025a0fff9e65e751407a69c5125a0e0dafd2e0048ce9b60d39bb0b58c251b4a72d382a02005750b091faae17b20a0a966b0c40ca44134fb214d0df0fed1c10646141f70"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	from, err := CtxSender(NewEIP155CtxSigner(big.NewInt(1)), tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}
