package tokens

import "time"

type TokenCreateRequest struct {
	Nombre string `json:"nombre" validate:"required,max=100"`
}

type TokenResp struct {
	ID        string     `json:"id"`
	Nombre    string     `json:"nombre"`
	Prefijo   string     `json:"prefijo"`
	Revocado  bool       `json:"revocado"`
	UltimoUso *time.Time `json:"ultimo_uso"`
	CreatedAt time.Time  `json:"created_at"`
}

// TokenCreatedResp solo se devuelve al crear: es la única vez que se ve el token completo.
type TokenCreatedResp struct {
	TokenResp
	Token string `json:"token"`
}

// TokenUser es el usuario dueño de un token válido (con su rol).
type TokenUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	IDRol     int    `json:"id_rol"`
	Rol       string `json:"rol"`
	IDEmpresa *int   `json:"id_empresa"`
	TokenID   string `json:"token_id"`
}
