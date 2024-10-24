package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

func (srv *ApiServer) configureCommonAPI(router *httprouter.Router) {
	// ASCOM Methods Common To All Devices
	router.PUT("/api/v1/switch/1/action", srv.handleNotSupported)
	router.PUT("/api/v1/focuser/1/action", srv.handleNotSupported)

	router.PUT("/api/v1/switch/1/commandblind", srv.handleNotSupported)
	router.PUT("/api/v1/focuser/1/commandblind", srv.handleNotSupported)

	router.PUT("/api/v1/switch/1/commandbool", srv.handleNotSupported)
	router.PUT("/api/v1/focuser/1/commandbool", srv.handleNotSupported)

	router.PUT("/api/v1/switch/1/commandstring", srv.handleNotSupported)
	router.PUT("/api/v1/focuser/1/commandstring", srv.handleNotSupported)

	router.GET("/api/v1/switch/1/connected", srv.handleConnected)
	router.GET("/api/v1/focuser/1/connected", srv.handleConnected)

	router.PUT("/api/v1/switch/1/connected", srv.handleConnect)
	router.PUT("/api/v1/focuser/1/connected", srv.handleConnect)

	router.GET("/api/v1/switch/1/description", srv.handleDescriptionCommon)
	router.GET("/api/v1/focuser/1/description", srv.handleDescriptionCommon)

	router.GET("/api/v1/switch/1/driverinfo", srv.handleDriverinfo)
	router.GET("/api/v1/focuser/1/driverinfo", srv.handleDriverinfo)

	router.GET("/api/v1/switch/1/driverversion", srv.handleDriverVersion)
	router.GET("/api/v1/focuser/1/driverversion", srv.handleDriverVersion)

	router.GET("/api/v1/switch/1/interfaceversion", srv.handleInterfaceVersion)
	router.GET("/api/v1/focuser/1/interfaceversion", srv.handleInterfaceVersion)

	router.GET("/api/v1/switch/1/name", srv.handleName)
	router.GET("/api/v1/focuser/1/name", srv.handleName)

	router.GET("/api/v1/switch/1/supportedactions", srv.handleSupportedActions)
	router.GET("/api/v1/focuser/1/supportedactions", srv.handleSupportedActions)
}

// ASCOM Common API handlers

// Retrieves the connected state of the device (GET)
func (srv *ApiServer) handleConnected(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// This is the GET commend
	result := MhpGetConnected()
	resp := booleanResponse{
		Value: result,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Sets the connected state of the device (PUT)
func (srv *ApiServer) handleConnect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cn, err := getConnectedFromRequest(r)
	if err != nil {
		resp := stringResponse{
			Value: err.Error(),
		}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	MhpSetConnect(cn)
	resp := stringResponse{
		Value: "",
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// The description of the device
func (srv *ApiServer) handleDescriptionCommon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result := "Mount Hub Pro"
	resp := stringResponse{
		Value: result,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// A string containing only the major and minor version of the driver.
func (srv *ApiServer) handleDriverinfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result := "Alpaca Mount Hub Pro Driver https://github.com/exploded/"
	resp := stringResponse{
		Value: result,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// A string containing only the major and minor version of the driver.
func (srv *ApiServer) handleDriverVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bi, _ := debug.ReadBuildInfo()
	resp := stringResponse{
		Value: bi.Main.Version,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// This method returns the version of the ASCOM device interface contract to which this
// device complies. Only one interface version is current at a moment in time and all new
// devices should be built to the latest interface version. Applications can choose which
// device interface versions they support and it is in their interest to support previous
// versions as well as the current version to ensure thay can use the largest number of devices.
func (srv *ApiServer) handleInterfaceVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := int32Response{
		Value: 2,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// The name of the device
func (srv *ApiServer) handleName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{
		Value: "Mount Hub Pro",
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the list of action names supported by this driver.
func (srv *ApiServer) handleSupportedActions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringlistResponse{
		Value: []string{""},
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
