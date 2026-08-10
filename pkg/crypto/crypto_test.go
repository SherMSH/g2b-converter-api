package crypto

import "testing"

// Транспортный ключ процессинга лежит в PKCS#1, а не в PKIX. Разбор только
// через PKIX ронял SetPIN до обращения к D8.
func TestReadPublicKeyPKCS1(t *testing.T) {
	key, err := ReadPublicKey("../../internal/app/files/transport_setpin.der")
	if err != nil {
		t.Fatalf("транспортный ключ не прочитан: %v", err)
	}
	if key.N.BitLen() == 0 {
		t.Error("ключ пустой")
	}
	t.Logf("размер ключа: %d бит", key.N.BitLen())
}

// Generate3DESKey отдаёт двойной ключ в 16 байт, а des.NewTripleDESCipher
// принимает только 24: без разворачивания K1|K2|K1 шифрование PIN-блока
// падало с "invalid key size 16".
func TestEncryptWith3DESAcceptsDoubleLengthKey(t *testing.T) {
	key, err := Generate3DESKey()
	if err != nil {
		t.Fatalf("генерация ключа: %v", err)
	}
	if len(key) != 16 {
		t.Fatalf("ожидался двойной ключ в 16 байт, получено %d", len(key))
	}

	block := make([]byte, 8)
	encrypted, err := EncryptWith3DES(key, block)
	if err != nil {
		t.Fatalf("шифрование двойным ключом: %v", err)
	}
	if len(encrypted) != len(block) {
		t.Errorf("длина результата %d, ожидалась %d", len(encrypted), len(block))
	}
}

// Длина данных обязана быть кратна блоку DES
func TestEncryptWith3DESRejectsBadLength(t *testing.T) {
	key, err := Generate3DESKey()
	if err != nil {
		t.Fatalf("генерация ключа: %v", err)
	}
	if _, err := EncryptWith3DES(key, make([]byte, 7)); err == nil {
		t.Error("ожидалась ошибка на некратной длине данных")
	}
}
