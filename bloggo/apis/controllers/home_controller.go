package controllers

import (
	"net/http"

	"github.com/ha-wk/bloggo/apis/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcoming you in this API")

}
