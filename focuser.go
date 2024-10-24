package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *ApiServer) configureFocuserAPI(router *httprouter.Router) {
	// ASCOM Methods specifc to the Focuser API
	router.GET("/setup/v1/focuser/1/setup", srv.handleFocuserSetup)
	router.GET("/api/v1/focuser/1/absolute", srv.handleAbsolute)
	router.GET("/api/v1/focuser/1/ismoving", srv.handleIsMoving)
	router.GET("/api/v1/focuser/1/maxincrement", srv.handleMaxIncrement)
	router.GET("/api/v1/focuser/1/maxstep", srv.handleMaxStep)
	router.GET("/api/v1/focuser/1/position", srv.handlePosition)
	router.GET("/api/v1/focuser/1/stepsize", srv.handleStepSize)
	router.GET("/api/v1/focuser/1/tempcomp", srv.handleTempComp)
	router.PUT("/api/v1/focuser/1/tempcomp", srv.handleNotSupported)
	router.GET("/api/v1/focuser/1/tempcompavailable", srv.handleTempCompAvailable)
	router.GET("/api/v1/focuser/1/temperature", srv.handleNotSupported)
	router.PUT("/api/v1/focuser/1/halt", srv.handleHalt)
	router.PUT("/api/v1/focuser/1/move", srv.handleMove)
}

// Handlers below are specific for the Focuser API

// Setup page.
func (srv *ApiServer) handleFocuserSetup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Alpaca MHP focuser server")
}

// True if the focuser is capable of absolute position; that is, being commanded to a specific step location.
func (srv *ApiServer) handleAbsolute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := booleanResponse{
		Value: true,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// True if the focuser is currently moving to a new position. False if the focuser is stationary.
func (srv *ApiServer) handleIsMoving(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := booleanResponse{
		Value: false, // JMC need to find a better way to determine if focuser is moving
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Maximum increment size allowed by the focuser; i.e. the maximum number of steps allowed in one move operation.
func (srv *ApiServer) handleMaxIncrement(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := int32Response{
		Value: int32(MhpGetMaxIncrement()),
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Maximum step position permitted.
func (srv *ApiServer) handleMaxStep(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := int32Response{
		Value: int32(MhpGetMaxStep()),
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Current focuser position, in steps.
func (srv *ApiServer) handlePosition(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := int32Response{
		Value: MhpGetPosition(),
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Step size (microns) for the focuser
func (srv *ApiServer) handleStepSize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := int32Response{
		Value: 1, // JMC TODO
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Gets the state of temperature compensation mode (if available), else always False.
func (srv *ApiServer) handleTempComp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := booleanResponse{
		Value: false,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// True if focuser has temperature compensation available.
func (srv *ApiServer) handleTempCompAvailable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := booleanResponse{
		Value: false,
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Immediately stop any focuser motion due to a previous Move(Int32) method call.
func (srv *ApiServer) handleHalt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := stringResponse{
		Value: "",
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Moves the focuser by the specified amount or to the specified position depending on the value of the Absolute property.
func (srv *ApiServer) handleMove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	value, err := getPositionFromRequest(r)
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

	// Move the focuser:
	err = MhpMove(value)
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
	resp := stringResponse{
		Value: "",
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
