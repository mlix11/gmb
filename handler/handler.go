package handler

import (
	"fmt"
	"github.com/loose11/gmb/database"
	"net/http"
)

type HandlerConfig struct {
	Database database.AppDatabase
}

func (h *HandlerConfig) Standard(w http.ResponseWriter, req *http.Request) {
	/*	files, _ := h.Database.GetAll()
		if len(files) > 0 {
			for _, value := range files {*/
	fmt.Fprintf(w, "test\n")
	/*}
	}*/
}
