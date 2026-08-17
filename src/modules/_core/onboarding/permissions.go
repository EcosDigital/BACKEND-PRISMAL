package onboarding

import (
	"fmt"

	modular "github.com/ecosistema/core/src/modules/_core/gestion_modular"
	"github.com/ecosistema/core/src/modules/_core/roles"
	"gorm.io/gorm"
)

// BuildDefaultAdminPermissions arma el árbol completo de módulos, funciones y
// subfunciones contratados en la licencia (código de licencia recién creado),
// con todo marcado, para que el rol Administrador del onboarding quede con
// acceso total desde el primer login sin pasos manuales adicionales.
//
// Debe invocarse mientras "db" sigue apuntando a la base de datos admin, ya
// que el catálogo de módulos/funciones/subfunciones vive únicamente ahí (la
// base del tenant solo guarda una referencia liviana en ref_modulos_tenant).
func BuildDefaultAdminPermissions(db *gorm.DB, codigoLicencia string) ([]roles.Item, error) {

	modulos, err := modular.FilterModulesByCodeLicence(codigoLicencia)
	if err != nil {
		return nil, fmt.Errorf("error consultando módulos de la licencia: %v", err)
	}

	tree := make([]roles.Item, 0, len(modulos))

	for _, m := range modulos {
		funciones, err := modular.FilterFuncionByIDModule(db, int64(m.ID))
		if err != nil {
			return nil, fmt.Errorf("error consultando funciones del módulo '%s': %v", m.Nombre, err)
		}

		modID := m.ID
		modTitle := m.Nombre
		modIcon := m.Icono
		modBgColor := m.BgColor
		modBrColor := m.BorderColor

		tree = append(tree, roles.Item{
			ID:       &modID,
			Title:    &modTitle,
			Icon:     &modIcon,
			BgColor:  &modBgColor,
			BrColor:  &modBrColor,
			Children: buildFuncionItems(funciones),
		})
	}

	return tree, nil
}

// buildFuncionItems convierte funciones (y sus subfunciones anidadas) al
// formato Item que espera la configuración de permisos de roles.
func buildFuncionItems(funciones []modular.FuncionesResponse) []roles.Item {

	items := make([]roles.Item, 0, len(funciones))

	for _, f := range funciones {
		fID := f.ID
		fTitle := f.Title

		item := roles.Item{
			ID:    &fID,
			Title: &fTitle,
			Path:  f.Path,
		}

		if len(f.Children) > 0 {
			item.Children = buildFuncionItems(f.Children)
		}

		items = append(items, item)
	}

	return items
}
