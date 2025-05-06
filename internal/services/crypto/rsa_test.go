package crypto

import (
	"bytes"
	"testing"
)

func TestRSAService_EncryptAndDecrypt(t *testing.T) {
	data := []byte("Какой-то очень секретный текст. Прям ну очень секретный. Отвечаю")

	r := &RSAService{}
	pub, priv, err := r.GenerateKeys(2048)
	if err != nil {
		t.Errorf("RSAService.GenerateKeys() has err = %v", err.Error())
	}
	c, err := r.Encrypt(pub, data)
	if err != nil {
		t.Errorf("RSAService.Encrypt() has err = %v", err.Error())
	}
	m, err := r.Decrypt(priv, c)
	if err != nil {
		t.Errorf("RSAService.Decrypt() has err = %v", err.Error())
	}
	if !bytes.Equal(data, m) {
		t.Errorf("RSAService want data == m, got data = %v, m = %v", data, m)
	}
}
