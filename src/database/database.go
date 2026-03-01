package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"

	_ "github.com/lib/pq"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
)

var tenantGormDB = map[string]*gorm.DB{}

// ConnectDB abre la conexión a PostgreSQL y devuelve *sql.DB
func ConnectDB(cfg *core.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db_host,
		cfg.Db_port,
		cfg.Db_user,
		cfg.Db_pass,
		cfg.Db_name,
	)

	// 1. Conexión base database/sql
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logging.Error.Printf("Error al abrir conexión con DB: %v", err)
		return nil, err
	}

	//configuracion de pool de conexiones
	db.SetMaxOpenConns(25)                 //maximo de conexiones abiertas
	db.SetMaxIdleConns(25)                 //maximo de conexiones inactivas
	db.SetConnMaxLifetime(5 * time.Minute) // recicla conexiones cada 5 minutos

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logging.Info.Printf("Conexión exitosa a PostgreSQL 🚀")

	// 2. Inicializar GORM SOBRE la misma conexión
	gormDB, err := gorm.Open(
		gormpostgres.New(gormpostgres.Config{
			Conn: db,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return nil, err
	}

	DB = db
	GormDB = gormDB

	logging.Info.Printf("Conexión PostgreSQL + GORM inicializada correctamente ✅")

	return db, nil

}

func ConnectTenantDB(dbName string) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		dbName,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Pool (igual que el principal)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(
		gormpostgres.New(gormpostgres.Config{
			Conn: sqlDB,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return nil, err
	}

	return gdb, nil

}

func GetGormDBByTenant(tenant string) (*gorm.DB, error) {

	// 1. Si no es multi-tenant, no cambia nada
	if os.Getenv("TENANCY_MODE") != "multi" {
		return GormDB, nil
	}

	// 2. Si el pool ya existe, reutilizarlo
	if db, ok := tenantGormDB[tenant]; ok {
		return db, nil
	}

	// 3. Construir nombre de la DB
	dbName := "prismar_" + tenant

	// 4. Crear pool nuevo
	gdb, err := ConnectTenantDB(dbName)
	if err != nil {
		return nil, err
	}

	// 5. Guardar pool en memoria
	tenantGormDB[tenant] = gdb

	return gdb, nil
}
