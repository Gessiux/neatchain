package keystore

import (
	"github.com/Gessiux/neatchain/common"
)

func KeyFileName(keyAddr common.Address) string {
	return keyFileName(keyAddr)
}

func WriteKeyStore(filepath string, keyjson []byte) error {
	return writeKeyFile(filepath, keyjson)
}
