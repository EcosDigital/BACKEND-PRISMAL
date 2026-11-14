package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

const tokenPrefix = "prm_"

// HashToken devuelve el SHA-256 en hex del token (lo que se guarda en la base).
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreateToken genera un token aleatorio, guarda su hash y lo devuelve completo una sola vez.
func CreateToken(db *gorm.DB, idUsuario int, nombre string) (*TokenCreatedResp, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, err
	}
	token := tokenPrefix + hex.EncodeToString(raw)
	prefijo := token[:12]

	id, err := InsertToken(db, idUsuario, nombre, HashToken(token), prefijo)
	if err != nil {
		return nil, err
	}

	return &TokenCreatedResp{
		TokenResp: TokenResp{
			ID:        id,
			Nombre:    nombre,
			Prefijo:   prefijo,
			CreatedAt: time.Now(),
		},
		Token: token,
	}, nil
}

// ValidateToken valida un token de conexión y devuelve su usuario (nil si no es válido).
func ValidateToken(db *gorm.DB, token string) (*TokenUser, error) {
	user, err := FindUserByTokenHash(db, HashToken(token))
	if err != nil || user == nil {
		return nil, err
	}
	_ = TouchToken(db, user.TokenID)
	return user, nil
}
