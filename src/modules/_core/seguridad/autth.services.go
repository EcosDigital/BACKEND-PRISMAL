package seguridad

import (
	"fmt"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserID    int
	Email     string
	Nombre    string
	Rol       string
	RolID     *int
	EmpresaID *int
	SedeID    *int
}

var jwtSecret = []byte(core.Cfg.Jwt_secret)

func GenerateBasicToken(user *UserInfo) (string, error) {

	jti := utils.GenerateJTI()
	ExpiresAt := time.Now().Add(10 * time.Minute)

	claims := utils.JWTClaims{
		UserID:            user.UserID,
		Email:             user.Email,
		RolID:             user.RolID,
		FullAuthenticated: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := signToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil

}

func GenerateFinalToken(user *UserInfo) (string, error) {

	jti := utils.GenerateJTI()
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := utils.JWTClaims{
		UserID:            user.UserID,
		Email:             user.Email,
		Nombre:            user.Nombre,
		Rol:               user.Rol,
		RolID:             user.RolID,
		EmpresaID:         user.EmpresaID,
		SedeID:            user.SedeID,
		FullAuthenticated: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := signToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil

}

func signToken(claims utils.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func UserAccessAuthorized(db *gorm.DB, claims *utils.JWTClaims) ([]AccessAutorized, error) {
	// Verificar que db no sea nil
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Comprobar accesos por empresa
	companies, err := GetCompanysAccess(db, *claims.RolID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas: %v", err)
	}

	// Comprobar accesos por sede
	sedes, err := GetSedesAccess(db, *claims.RolID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener sedes: %v", err)
	}

	// Construir resultado
	result := make([]AccessAutorized, 0, len(companies))

	for _, company := range companies {
		var sedesEmpresa []SedesAccess
		for _, s := range sedes {
			if s.IdEmpresa == company.IdEmpresa {
				sedesEmpresa = append(sedesEmpresa, s)
			}
		}

		result = append(result, AccessAutorized{
			EmpresaID:   company.IdEmpresa,
			RazonSocial: company.RazonSocial,
			Descripcion: company.Descripcion,
			UsaSedes:    company.UsaSedes,
			Sedes:       sedesEmpresa,
		})
	}
	return result, nil

}
