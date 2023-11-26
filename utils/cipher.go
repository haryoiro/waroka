package utils

import (
	"golang.org/x/crypto/blowfish"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func Encrypt(pt []byte) ([]byte, error) {
	key := []byte(GetSecret())
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

func Decrypt(ct []byte) ([]byte, error) {
	key := []byte(GetSecret())
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		return nil, err
	}
	return pt, nil
}
