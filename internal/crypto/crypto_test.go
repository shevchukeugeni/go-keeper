package crypto

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	_, err := Encrypt("test_pass", "test_data")
	require.NoError(t, err)
}

func TestDecrypt(t *testing.T) {
	res, err := Decrypt("test_pass", "d3c614d4892893917ef052b9be654ad39e9930a681c1b17004")
	require.NoError(t, err)
	require.Equal(t, "test_data", res)

	res, err = Decrypt("test_pass", "nnnn")
	require.Equal(t, err, hex.InvalidByteError(("n")[0]))
	require.Equal(t, "", res)

	res, err = Decrypt("test_pass2", "d3c614d4892893917ef052b9be654ad39e9930a681c1b17004")
	require.Equal(t, err, errors.New("cipher: message authentication failed"))
	require.Equal(t, "", res)
}
