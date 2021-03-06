package httpServer

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/koestler/go-ve-sensor/dataflow"
	"github.com/koestler/go-ve-sensor/storage"
	"net/http"
)

func HandleDeviceIndex(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	devices := storage.GetAll()

	// cache device index for 5 minutes
	w.Header().Set("Cache-Control", "public, max-age=300")
	writeJsonHeaders(w)

	b, err := json.MarshalIndent(devices, "", "    ")
	if err != nil {
		return StatusError{500, err}
	}
	w.Write(b)

	return nil;
}

func HandleDeviceGetRoundedValues(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	vars := mux.Vars(r)

	device, err := storage.GetByName(vars["DeviceId"])
	if err != nil {
		return StatusError{404, err}
	}

	roundedValues := env.RoundedStorage.GetMap(dataflow.Filter{Devices: map[*storage.Device]bool{device: true}})
	roundedValuesEssential := roundedValues.ConvertToEssential()

	writeJsonHeaders(w)
	b, err := json.MarshalIndent(roundedValuesEssential, "", "    ")
	if err != nil {
		return StatusError{500, err}
	}
	w.Write(b)
	return nil
}

func HandleDeviceGetPictureThumb(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	return HandleDeviceGetPicture(env, w, r, true)
}

func HandleDeviceGetPictureRaw(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	return HandleDeviceGetPicture(env, w, r, false)
}

func HandleDeviceGetPicture(env *Environment, w http.ResponseWriter, r *http.Request, thumb bool) Error {
	vars := mux.Vars(r)

	device, err := storage.GetByName(vars["DeviceId"])

	if err != nil {
		return StatusError{404, err}
	}

	picture, err := storage.PictureDb.GetPicture(device)
	if err != nil {
		return StatusError{404, err}
	}

	var jpeg []byte
	if thumb {
		jpeg = picture.JpegThumb
	} else {
		jpeg = picture.JpegRaw
	}

	writeJpegHeaders(w)
	if _, err = w.Write(jpeg); err != nil {
		return StatusError{500, err}
	}
	return nil
}
