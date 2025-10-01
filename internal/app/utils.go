package app

import "net/http"

func (app *Application) commitHeadersAndWriteStatus(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
}
