package crypto_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gangjun06/d4dj-info-server/utils/crypto"
)

func TestDecrypt(t *testing.T) {
	fileName := "./test.png.enc"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Error("Error import test file. place test file to utils/crypto/")
	}
	data, err := crypto.New().Decrypt(file)
	if err != nil {
		t.Error("Error decrypt file: ", err.Error())
	}
	ioutil.WriteFile(strings.ReplaceAll(fileName, ".enc", ""), data, 0644)
}
