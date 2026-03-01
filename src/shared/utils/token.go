package utils

import (
	"errors"

	"github.com/ecosistema/core/src/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID            int    `json:"user_id"`
	Email             string `json:"email"`
	Nombre            string `json:"nombre,omitempty"`
	Rol               string `json:"rol_omitempty"`
	RolID             *int   `json:"rol_id,omitempty"`
	EmpresaID         *int   `json:"id_empresa,omitempty"`
	SedeID            *int   `json:"id_sede,omitempty"`
	FullAuthenticated bool   `json:"full_authenticated"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(core.Cfg.Jwt_secret)

// ValidateJWT → valida token sin importar si es final o básico
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token inválido")
	}
	return token.Claims.(*JWTClaims), nil
}

// DecodeTokenFinal → valida que sea token final
func DecodeTokenFinal(tokenString string) (*JWTClaims, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}
	if !claims.FullAuthenticated || claims.EmpresaID == nil || *claims.EmpresaID == 0 {
		return nil, errors.New("token final inválido")
	}
	return claims, nil
}

// funcion for redis
func GenerateJTI() string {
	return uuid.New().String()
}
