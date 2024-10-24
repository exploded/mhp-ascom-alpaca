package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

// Management API router
func (srv *ApiServer) configureManagementAPI(router *httprouter.Router) {
	router.GET("/", srv.handleRoot)
	router.GET("/management/apiversions", srv.handleApiVersions)
	router.GET("/management/v1/description", srv.handleDescription)
	router.GET("/management/v1/configureddevices", srv.handleConfiguredDevices)
}

// Returns root web page.
func (srv *ApiServer) handleRoot(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Alpaca MHP server")
}

// Returns an integer array of supported Alpaca API version numbers.
func (srv *ApiServer) handleApiVersions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := uint32listResponse{
		Value: []uint32{1},
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Returns cross-cutting information that applies to all devices available at this URL:Port.
func (srv *ApiServer) handleDescription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bi, _ := debug.ReadBuildInfo()
	resp := managementDescriptionResponse{
		Value: ServerDescription{
			ServerName:          "Alpaca MHP",
			Manufacturer:        "https://github.com/exploded/mhp-ascom-alpaca",
			ManufacturerVersion: bi.Main.Version,
			Location:            Location,
		},
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Returns an array of device description objects, providing unique information for each served device,
// enabling them to be accessed through the Alpaca Device API.
func (srv *ApiServer) handleConfiguredDevices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	val := MhpGetInit()
	resp := managementDevicesListResponse{
		Value: val,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
