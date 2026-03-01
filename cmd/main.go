package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	loader "github.com/ecosistema/core/src"
	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/database"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	//cargar config
	core.LoadConfig()

	//inicaliar logger
	logging.InitLogger()

	//Crear directorio uploads
	if err := utils.EnsureUploadDir(); err != nil {
		logging.Error.Printf("Error creando directorio: %v", err)
	}

	//servdidor con fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 5MB
	})

	// Middlewares globales
	app.Use(logger.New())  //logg de todas las peticiones
	app.Use(recover.New()) //atrapa cualquier panic

	/*app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie, X-Requested-With, X-Tenant-ID",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           43200,

		AllowOriginsFunc: func(origin string) bool {

			// SINGLE TENANT
			if core.Cfg.TenancyMode == "single" {
				return origin == core.Cfg.Frontend_url
			}

			// MULTI TENANT
			baseURL, err := url.Parse(core.Cfg.Frontend_url)
			if err != nil {
				return false
			}

			reqURL, err := url.Parse(origin)
			if err != nil {
				return false
			}

			baseHost := baseURL.Hostname()
			reqHost := reqURL.Hostname()

			if reqHost == baseHost {
				return true
			}

			if core.Cfg.Allow_subdomains {
				return strings.HasSuffix(reqHost, "."+baseHost)
			}

			return false
		},

		AllowOrigins: core.Cfg.Frontend_url,
	}))*/

	// Reemplaza todo el bloque app.Use(cors.New(...)) por esto:

	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		if isAllowedOrigin(origin) {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Cookie, X-Requested-With, X-Tenant-ID")
			c.Set("Access-Control-Expose-Headers", "Set-Cookie")
			c.Set("Access-Control-Max-Age", "43200")
		}

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	//conexión DB
	db, err := database.ConnectDB(&core.Cfg)
	if err != nil {
		logging.Error.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// Servir imágenes públicas
	wd, _ := os.Getwd()
	uploadPath := filepath.Join(wd, "src", "public", "uploads")
	app.Static("/uploads", uploadPath)

	//ejecutar migraciones
	if core.Cfg.TenancyMode == "single" {
		database.RunMigrations(db)
	}

	//app.Use(middlewares.TenancyMiddleware) //tenantMode

	app.Use(func(c *fiber.Ctx) error {
		// Rutas públicas del backend admin que no requieren tenant
		if strings.HasPrefix(c.Path(), "/licence") {
			return c.Next()
		}
		return middlewares.TenancyMiddleware(c)
	})

	//registro de rutas (configuraciones)
	loader.SetupRoutes(app)

	//construir direcion con el puerto del arhivo de configuracion
	puerto := core.Cfg.Backend_port
	logging.Info.Printf("🚀 Servidor corriendo en el puerto: %s", puerto)

	//levantar servidor
	if err := app.Listen(puerto); err != nil {
		logging.Error.Printf("❌ Error al iniciar servidor: %v", err)

	}

}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	reqURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	baseURL, err := url.Parse(core.Cfg.Frontend_url)
	if err != nil {
		return false
	}

	reqHost := reqURL.Hostname()
	baseHost := baseURL.Hostname()

	// Dominio exacto siempre permitido en cualquier modo
	if reqHost == baseHost {
		return true
	}

	// Subdominios permitidos si está habilitado (aplica para AMBOS modos)
	if core.Cfg.Allow_subdomains {
		return strings.HasSuffix(reqHost, "."+baseHost)
	}

	return false
}
