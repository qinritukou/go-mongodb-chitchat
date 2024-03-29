package main

import (
	"go-mongodb-chitchat/chichat/dao"
	"net/http"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "main.layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "main.layout", "private.navbar", "error")
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	threads, err := dao.Threads()

	if err != nil {
		errorMessage(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "main.layout", "public.navbar", "index")
		} else {
			generateHTML(writer, threads, "main.layout", "private.navbar", "index")
		}
	}
}
