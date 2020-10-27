package scrypt

import (
	"bytes"
	"testing"

	"gbchain-org/go-gbchain/common/hexutil"
)

// Tests whether the ScryptHash lookup works
func TestScryptHash(t *testing.T) {
	// Create a block to verify
	hash := hexutil.MustDecode("0x885c778d7eedb68876b1377e216ed1d2c2417b0fca06b66ca4facae79ae5330d")
	nonce := uint64(3249874452068615500)

	wantDigest := hexutil.MustDecode("0xa926c4799edcb96b973634888e610fa9f0ca66b4d170903f80fe99487785414e")
	wantResult := hexutil.MustDecode("0xec9aa0657969e59514b6546d36c706f5aa1625b1f471950a9e6a009452308297")

	digest, result := ScryptHash(hash, nonce)
	if !bytes.Equal(digest, wantDigest) {
		t.Errorf("ScryptHash digest mismatch: have %x, want %x", digest, wantDigest)
	}
	if !bytes.Equal(result, wantResult) {
		t.Errorf("ScryptHash result mismatch: have %x, want %x", result, wantResult)
	}

}

// Benchmarks the verification performance
func BenchmarkScryptHash(b *testing.B) {
	hash := hexutil.MustDecode("0x885c778d7eedb68876b1377e216ed1d2c2417b0fca06b66ca4facae79ae5330d")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScryptHash(hash, 0)
	}
}
