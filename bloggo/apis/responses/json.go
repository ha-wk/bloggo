package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, statuscode int, data interface{}) {
	w.WriteHeader(statuscode)              //writing response of server
	err := json.NewEncoder(w).Encode(data) //encoding the data to json
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statuscode int, err error) {
	if err != nil {
		JSON(w, statuscode, struct { //3rd parameter of JSON is method signature i.e interface
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil) //NIL SPECIFIES ZERO METHODS
}
