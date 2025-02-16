package main

import (
	"net/http"

	"github.com/MohammadBohluli/social-app-go/pkg"
)

func (app application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "every thing is OK 😊",
	}

	if err := pkg.WriteJson(w, http.StatusOK, data); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}
