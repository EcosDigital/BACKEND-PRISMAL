package core

import (
	empresa "github.com/ecosistema/core/src/modules/_core/gestion_empresa"
	modular "github.com/ecosistema/core/src/modules/_core/gestion_modular"
	onboarding "github.com/ecosistema/core/src/modules/_core/onboarding"
	"github.com/ecosistema/core/src/modules/_core/referenciales"
	"github.com/ecosistema/core/src/modules/_core/roles"
	"github.com/ecosistema/core/src/modules/_core/search"
	"github.com/ecosistema/core/src/modules/_core/seguridad"
	"github.com/ecosistema/core/src/modules/_core/terceros"
	uploads "github.com/ecosistema/core/src/modules/_core/up"
	"github.com/ecosistema/core/src/modules/_core/usuarios"
	"github.com/gofiber/fiber/v2"
)

func RegisterCoreRoutes(r fiber.Router) {

	core := r.Group("/core")

	//rutas onBoarding
	onboarding.Rutas_onBoarding(core)
	uploads.Rutas_Uploads(core)

	seguridad.Rutas_seguridad(core)

	//busquedas dinamicas
	search.SearchDynamicsCore(core)
	referenciales.SetupRoutes(core) // referenciales
	modular.Rutas_Modular(core)     //gestion modular
	empresa.Rutas_empresa(core)
	empresa.Rutas_Sede(core)

	terceros.Rutas_Terceros(core)
	roles.Rutas_Roles(core)
	usuarios.Rutas_Usuarios(core) //usuarios

}
