package main

import (
	"net/http"

	"github.com/MohammadBohluli/social-app-go/pkg"
)

// HealthCheck godoc
//
//	@Summary		Health check endpoint
//	@Description	Checks if the API is running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"API is healthy"
//	@Router			/health [get]
func (app application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "every thing is OK 😊",
	}

	if err := pkg.WriteJson(w, http.StatusOK, data); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}
