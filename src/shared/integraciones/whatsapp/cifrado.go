package whatsapp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/ecosistema/core/src/core"
)

func nuevoCifrador() (cipher.AEAD, error) {
	if core.Cfg.Whatsapp_clave_cifrado == "" {
		return nil, errors.New("falta la clave de cifrado de WhatsApp")
	}

	clave := sha256.Sum256([]byte(core.Cfg.Whatsapp_clave_cifrado))

	bloque, err := aes.NewCipher(clave[:])
	if err != nil {
		return nil, errors.New("no se pudo preparar el cifrado")
	}

	return cipher.NewGCM(bloque)
}

// CifrarLlave convierte la llave de Meta en un texto ilegible para guardarlo.
func CifrarLlave(llave string) (string, error) {
	cifrador, err := nuevoCifrador()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, cifrador.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("no se pudo cifrar la llave")
	}

	sellado := cifrador.Seal(nonce, nonce, []byte(llave), nil)
	return base64.StdEncoding.EncodeToString(sellado), nil
}

// DescifrarLlave recupera la llave original justo antes de usarla.
func DescifrarLlave(cifrada string) (string, error) {
	cifrador, err := nuevoCifrador()
	if err != nil {
		return "", err
	}

	datos, err := base64.StdEncoding.DecodeString(cifrada)
	if err != nil || len(datos) < cifrador.NonceSize() {
		return "", errors.New("la llave guardada no es válida")
	}

	nonce, contenido := datos[:cifrador.NonceSize()], datos[cifrador.NonceSize():]
	llave, err := cifrador.Open(nil, nonce, contenido, nil)
	if err != nil {
		return "", errors.New("la llave guardada no es válida")
	}

	return string(llave), nil
}
