package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koestler/go-ve-sensor/vedata"
	"github.com/koestler/go-ve-sensor/vehttp"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:generate ./frontend_to_bindata.sh

var HttpRoutes = vehttp.Routes{
	vehttp.Route{
		"DeviceIndex",
		"GET",
		"/api/v0/device/",
		HttpHandleDeviceIndex,
	},
	vehttp.Route{
		"DeviceIndex",
		"GET",
		"/api/v0/device/{DeviceId:[a-zA-Z0-9\\-]{1,32}}",
		HttpHandleDeviceGet,
	},
	vehttp.Route{
		"Index",
		"GET",
		"/",
		HttpHandleAssetsGet,
	},
	vehttp.Route{
		"Assets",
		"GET",
		"/{Path:.+}",
		HttpHandleAssetsGet,
	},
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func writeJsonHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
}

func jsonWriteNotFound(w http.ResponseWriter) {
	jsonWriteError(w, "Object Not Found")
}

func jsonWriteError(w http.ResponseWriter, text string) {
	writeJsonHeaders(w, http.StatusNotFound)

	ret := jsonErr{Code: http.StatusNotFound, Text: text}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func HttpHandleDeviceIndex(w http.ResponseWriter, r *http.Request) {
	deviceIds := vedata.ReadDeviceIds()

	// cache device index for 5 minutes
	w.Header().Set("Cache-Control", "public, max-age=300")
	writeJsonHeaders(w, http.StatusOK)

	b, err := json.MarshalIndent(deviceIds, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func HttpHandleDeviceGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deviceId := vedata.DeviceId(vars["DeviceId"])
	device, err := deviceId.ReadDevice()

	// cache date for 5 seconds
	w.Header().Set("Cache-Control", "public, max-age=5")

	if err == nil {
		writeJsonHeaders(w, http.StatusOK)

		b, err := json.MarshalIndent(device, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(b)
	} else {
		jsonWriteNotFound(w)
	}
}

func HttpHandleAssetsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	path := vars["Path"]
	if path == "" {
		path = "index.html"
	}

	// cache static files and 404 for one day
	w.Header().Set("Cache-Control", "public, max-age=86400")

	if bs, err := Asset(path); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 asset not found\n")
		log.Printf("handlers: %v", err)
	} else {
		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}

		w.WriteHeader(http.StatusOK)
		var reader = bytes.NewBuffer(bs)
		io.Copy(w, reader)
	}
}
