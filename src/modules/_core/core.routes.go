package core

import (
	asistente "github.com/ecosistema/core/src/modules/_core/asistente_ia"
	"github.com/ecosistema/core/src/modules/_core/comprobantes"
	estructurafisica "github.com/ecosistema/core/src/modules/_core/estructura_fisica"
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
	"github.com/ecosistema/core/src/modules/gestiones"
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
	empresa.Rutas_empresa(core)     // gestion empresa
	empresa.Rutas_Sede(core)        // gestion sedes

	terceros.Rutas_Terceros(core) //gestion terceros
	roles.Rutas_Roles(core)       // gestion roles
	usuarios.Rutas_Usuarios(core) //usuarios

	estructurafisica.Rutas_EstructuraFisica(core) // estructura fisica (mesas)

	gestiones.Rutas_Gestiones(core)

	comprobantes.Rutas_comprobante(core) //gestion comprobantes

	asistente.Rutas_Asistente(core) //asistente de IA (Copilot)

}
