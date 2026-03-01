package uploads

import (
	"github.com/gofiber/fiber/v2"
)

func Rutas_Uploads(r fiber.Router) {
	api := r.Group("/uploads")

	api.Post("/images", UploadImageController)
	api.Delete("/images/:filename", DeleteImageController)

}
