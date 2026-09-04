package onboarding

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ecosistema/core/src/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateTenantDB(slugTenant string) error {

	dbName := fmt.Sprintf("prismar_%s", slugTenant)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error al conectar: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error al hacer ping: %v", err)
	}

	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error al crear BD: %v", err)
	}

	//EJECUTAR MIGRACIONES
	err = ExecuteTenantMigrations(dbName)
	if err != nil {
		return fmt.Errorf("error ejecutando migraciones: %v", err)
	}

	return nil

}

func ExecuteTenantMigrations(dbName string) error {

	// Conectar a la BD del tenant
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
		dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Archivos que NO se deben ejecutar en tenants
	excludedFiles := []string{
		"007_init.up.sql",
		"011_init.up.sql",
		"015_init.up.sql",
		"017_init.up.sql",
		"018_init.up.sql",
		"020_init.up.sql",
		"022_init.up.sql",
		"023_init.up.sql",
		"024_init.up.sql",
		"025_init.up.sql",
		"026_init.up.sql",
	}

	// Ruta de las migraciones
	migrationsPath := "./migrations"

	// Leer todos los archivos .sql
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("error leyendo carpeta de migraciones: %v", err)
	}

	if len(files) == 0 {
		return nil
	}

	// Ordenar archivos alfabéticamente
	sort.Strings(files)

	for _, file := range files {
		fileName := filepath.Base(file)

		skip := false
		for _, excluded := range excludedFiles {
			if fileName == excluded {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		sqlContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %v", file, err)
		}

		statements := splitSQLStatements(string(sqlContent))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("error ejecutando migración %s: %v", fileName, err)
			}

		}

	}

	return nil

}

func splitSQLStatements(content string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDollarQuote := false
	dollarTag := ""

	i := 0
	for i < len(content) {
		char := content[i]

		// ── Saltar comentarios de línea (--) ──────────────────────────
		// Solo si no estamos dentro de ningún string
		if char == '-' && !inSingleQuote && !inDollarQuote &&
			i+1 < len(content) && content[i+1] == '-' {
			// Avanzar hasta el fin de línea sin escribir en current
			for i < len(content) && content[i] != '\n' {
				i++
			}
			// Escribir el salto de línea si existe (preserva formato)
			if i < len(content) {
				current.WriteByte(content[i])
				i++
			}
			continue
		}

		// ── Saltar comentarios de bloque (/* ... */) ──────────────────
		if char == '/' && !inSingleQuote && !inDollarQuote &&
			i+1 < len(content) && content[i+1] == '*' {
			i += 2
			for i < len(content) {
				if content[i] == '*' && i+1 < len(content) && content[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
			continue
		}

		// ── Manejar comillas simples ───────────────────────────────────
		if char == '\'' && !inDollarQuote {
			if i+1 < len(content) && content[i+1] == '\'' {
				// Escape de comilla: ''
				current.WriteByte(char)
				i++
				current.WriteByte('\'')
				i++
				continue
			}
			inSingleQuote = !inSingleQuote
			current.WriteByte(char)
			i++
			continue
		}

		// ── Manejar dollar quotes ($$ o $tag$) ────────────────────────
		if char == '$' && !inSingleQuote {
			tag := "$"
			j := i + 1
			for j < len(content) && content[j] != '$' {
				tag += string(content[j])
				j++
			}
			if j < len(content) {
				tag += "$"
			}

			if !inDollarQuote {
				inDollarQuote = true
				dollarTag = tag
				current.WriteString(tag)
				i = j + 1
			} else if tag == dollarTag {
				inDollarQuote = false
				dollarTag = ""
				current.WriteString(tag)
				i = j + 1
			} else {
				current.WriteByte(char)
				i++
			}
			continue
		}

		// ── Separador de statements ───────────────────────────────────
		if char == ';' && !inSingleQuote && !inDollarQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			i++
			continue
		}

		current.WriteByte(char)
		i++
	}

	// Último statement sin ; final
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func ConnectToTenantDB(dbName string) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
		dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error conectando a BD tenant: %v", err)
	}

	return db, nil
}
