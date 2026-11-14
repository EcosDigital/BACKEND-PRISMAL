package tokens

import "gorm.io/gorm"

// InsertToken guarda el hash del token y devuelve el id (uuid) del registro.
func InsertToken(db *gorm.DB, idUsuario int, nombre, hash, prefijo string) (string, error) {
	var id string
	err := db.Raw(`
		INSERT INTO seguridad.cfg_tokens_conexion
			(id_usuario, nombre, token_hash, prefijo, created_by)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id`,
		idUsuario, nombre, hash, prefijo, idUsuario,
	).Row().Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// ListTokensByUser lista los tokens del usuario (sin el hash).
func ListTokensByUser(db *gorm.DB, idUsuario int) ([]TokenResp, error) {
	var results []TokenResp
	err := db.Raw(`
		SELECT id, nombre, prefijo, revocado, ultimo_uso, created_at
		FROM seguridad.cfg_tokens_conexion
		WHERE id_usuario = ?
		ORDER BY created_at DESC`, idUsuario,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []TokenResp{}
	}
	return results, nil
}

// RevokeToken revoca un token solo si pertenece al usuario. Devuelve false si no se encontró.
func RevokeToken(db *gorm.DB, id string, idUsuario int) (bool, error) {
	res := db.Exec(`
		UPDATE seguridad.cfg_tokens_conexion
		SET revocado = true
		WHERE id = ? AND id_usuario = ? AND revocado = false`,
		id, idUsuario,
	)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// FindUserByTokenHash devuelve el usuario dueño del token solo si el token no
// está revocado y el usuario sigue activo. Si no, devuelve nil.
func FindUserByTokenHash(db *gorm.DB, hash string) (*TokenUser, error) {
	var results []TokenUser
	err := db.Raw(`
		SELECT u.id, u.email, u.id_rol, r.nombre AS rol, r.id_empresa, t.id AS token_id
		FROM seguridad.cfg_tokens_conexion t
		INNER JOIN seguridad.cfg_usuarios u ON u.id = t.id_usuario
		INNER JOIN seguridad.cfg_roles_usuario r ON r.id = u.id_rol
		WHERE t.token_hash = ? AND t.revocado = false AND u.activo = true
		LIMIT 1`, hash,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[0], nil
}

// TouchToken registra la fecha del último uso.
func TouchToken(db *gorm.DB, id string) error {
	return db.Exec(`UPDATE seguridad.cfg_tokens_conexion SET ultimo_uso = NOW() WHERE id = ?`, id).Error
}
