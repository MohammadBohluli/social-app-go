package main

import (
	"net/http"

	"github.com/MohammadBohluli/social-app-go/pkg"
)

func (app application) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"status": "every thing is OK 😊",
	}

	if err := pkg.WriteJson(w, http.StatusOK, data); err != nil {
		pkg.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
