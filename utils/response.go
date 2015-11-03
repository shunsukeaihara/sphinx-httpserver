package utils

import (
	"net/http"

	"github.com/unrolled/render"
)

var renderer = render.New(render.Options{})

func WriteJsonErrorResponse(w http.ResponseWriter, message string, code int) {
	renderer.JSON(w, code, map[string]string{
		"message": message,
	})
}
