package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/google/uuid"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var UploadDir string

func init() {
	wd, _ := os.Getwd()
	UploadDir = filepath.Join(wd, "src", "public", "uploads")
}

// crear directorio (si no existe)
func EnsureUploadDir() error {

	dir := filepath.Clean(UploadDir)

	//verificar existencia
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		//crear con permisos
		if err := os.Mkdir(dir, 0755); err != nil {
			logging.Error.Printf("❌ Error creando directorio: %v", err)
		}

		logging.Info.Println("✅ Directorio creado exitosamente")
	}

	return nil
}

// validar imagen
func ValidteImage(file *multipart.FileHeader) error {
	// validar tamaño (maximo 10 mb)
	if file.Size > MaxFileSize {
		logging.Error.Printf("Archivo muy grande (máx 10MB)")
		return fmt.Errorf("Archivo muy grande (máx 10MB)")
	}

	//validar extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		logging.Error.Printf("Archivo con extension incorrecta")
		return fmt.Errorf("solo se permiten .jpg, .jpeg, .png")
	}

	return nil

}

// guardar imagen
func SaveImage(file *multipart.FileHeader) (string, error) {
	//generar nombre unico
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext

	//existe direvtorio
	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		return "", err
	}

	//ruta completa
	filepath := filepath.Join(UploadDir, filename)

	//abrir un archivo subido
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	//crear archivo destino
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	//copiar contenido
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil

}

// Eliminar imagen
func DeleteImage(filename string) error {
	filepath := filepath.Join(UploadDir, filename)
	return os.Remove(filepath)
}
