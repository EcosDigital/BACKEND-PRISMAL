package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logging.Error.Printf("Error al crear el driver de migraciones: %v", err)
		log.Fatalf("Error creando driver para migraciones: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", //carpeta de migraciones
		"postgres", driver,
	)
	if err != nil {
		logging.Error.Printf("Error al procesar las migraciones: %v", err)
		log.Fatalf("Error cargando migraciones: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logging.Error.Printf("Error al aplicar las migraciones: %v", err)
		log.Fatalf("Error aplicando migraciones %v", err)
	}

	logging.Info.Print("Migraciones aplicadas correctamente ✅")
	fmt.Println("✅ Migraciones exitosaas...")

}
