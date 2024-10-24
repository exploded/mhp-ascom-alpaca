package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *ApiServer) configureSwitchAPI(router *httprouter.Router) {
	// ASCOM Methods specifc to the Switch API
	router.GET("/setup/v1/switch/1/setup", srv.handleSwitchSetup)
	router.GET("/api/v1/switch/1/maxswitch", srv.handleMaxSwitch)
	router.GET("/api/v1/switch/1/canwrite", srv.handleCanWrite)
	router.GET("/api/v1/switch/1/getswitch", srv.handleGetSwitch)
	router.GET("/api/v1/switch/1/getswitchdescription", srv.handleGetSwitchDescription)
	router.GET("/api/v1/switch/1/getswitchname", srv.handleGetSwitchName)
	router.GET("/api/v1/switch/1/getswitchvalue", srv.handleGetSwitchValue)
	router.GET("/api/v1/switch/1/minswitchvalue", srv.handleMinSwitchValue)
	router.GET("/api/v1/switch/1/maxswitchvalue", srv.handleMaxSwitchValue)
	router.PUT("/api/v1/switch/1/setswitch", srv.handleSetSwitch)
	router.PUT("/api/v1/switch/1/setswitchname", srv.handleSetSwitchName)
	router.PUT("/api/v1/switch/1/setswitchvalue", srv.handleSetSwitchValue)
	router.GET("/api/v1/switch/1/switchstep", srv.handleSwitchStep)
}

// Handlers below are specific for the Switch API

// Setup page.
func (srv *ApiServer) handleSwitchSetup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Alpaca MHP switch server")
}

// Returns the number of switch devices managed by this driver. Devices are numbered from 0 to MaxSwitch - 1
func (srv *ApiServer) handleMaxSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := int32Response{
		Value: NumOnOffSwitch + NumVarSwitch,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Reports if the specified switch device can be written to, default true.
// This is false if the device cannot be written to, for example a limit switch or a sensor.
// Devices are numbered from 0 to MaxSwitch - 1
func (srv *ApiServer) handleCanWrite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := booleanResponse{
		Value: true, // all switches can be changed
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Return the state of switch device id as a boolean. Devices are numbered from 0 to MaxSwitch - 1
func (srv *ApiServer) handleGetSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		result, err := MhpGetOnOff(sn)
		if err != nil {
			resp := stringResponse{
				Value: err.Error(),
			}
			srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := booleanResponse{
				Value: result,
			}
			srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}
	}
}

// Gets the description of the specified switch device.
// This is to allow a fuller description of the device to be returned,
// for example for a tool tip. Devices are numbered from 0 to MaxSwitch - 1
func (srv *ApiServer) handleGetSwitchDescription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)

	} else {
		result := MhpGetName(sn)
		resp := stringResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// Gets the name of the specified switch device. Devices are numbered from 0 to MaxSwitch - 1
func (srv *ApiServer) handleGetSwitchName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		result := MhpGetName(sn)
		resp := stringResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// Gets the value of the specified switch device as a double.
// Devices are numbered from 0 to MaxSwitch - 1,
// The value of this switch is expected to be between MinSwitchValue and MaxSwitchValue.
func (srv *ApiServer) handleGetSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		result := MhpGetValue(sn)
		resp := doubleResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

}

// Gets the minimum value of the specified switch device as a double. Devices are numbered from 0 to MaxSwitch - 1.
func (srv *ApiServer) handleMinSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		result := MhpGetMin(sn)
		resp := doubleResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

// Gets the maximum value of the specified switch device as a double. Devices are numbered from 0 to MaxSwitch - 1.
func (srv *ApiServer) handleMaxSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		result := MhpGetMax(sn)
		resp := doubleResponse{
			Value: result,
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

}

// Sets a switch controller device to the specified state, true or false.
func (srv *ApiServer) handleSetSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Note this is a PUT
	// Get the switch number from the the request
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		// Get the required value from the request
		sv, err := getSwitchStateFromRequest(r)
		if err != nil {
			resp := stringResponse{
				Value: err.Error(),
			}
			srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			err := MhpSetOnOff(sn, sv)
			if err != nil {
				resp := stringResponse{
					Value: err.Error(),
				}
				srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			} else {
				resp := stringResponse{
					Value: "Sucess",
				}
				srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}
}

// Sets a switch device name to the specified value.
func (srv *ApiServer) handleSetSwitchName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Note this is a PUT
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	sna, err := getSwitchNameFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = MhpSetName(sn, sna)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := putResponse{}

	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Sets a switch device value to the specified value.
func (srv *ApiServer) handleSetSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Note this is a PUT
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	sv, err := getValueFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = MhpSetValue(sn, sv)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := putResponse{}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the step size that this device supports (the difference between successive values of the device).
// Devices are numbered from 0 to MaxSwitch - 1.
func (srv *ApiServer) handleSwitchStep(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	result := MhpGetStep(sn) // should always be 1

	resp := doubleResponse{
		Value: result,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
