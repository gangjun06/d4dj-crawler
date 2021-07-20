package crypto_test

import (
	"io/ioutil"
	"testing"

	"github.com/gangjun06/d4dj-info-server/utils/crypto"
)

func TestDecrypt(t *testing.T) {
	file, err := ioutil.ReadFile("./test.png.enc")
	if err != nil {
		t.Error("Error import test file. place test.png.enc file to utils/crypto/")
	}
	_, err = crypto.New().Decrypt(file)
	if err != nil {
		t.Error("Error decrypt file: ", err.Error())
	}
}
