package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ClientId string

type ConnectedClient struct {
	ClientId            ClientId
	ClientTransactionID uint32
	Connected           bool
}

type ApiServer struct {
	ApiPort             uint32
	ServerTransactionID uint32
}

func NewApiServer(apiPort uint32) *ApiServer {
	return &ApiServer{
		ApiPort: apiPort,
	}
}

func (srv *ApiServer) Start() {

	router := httprouter.New()
	srv.configureManagementAPI(router)
	srv.configureCommonAPI(router)
	srv.configureSwitchAPI(router)
	srv.configureFocuserAPI(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", srv.ApiPort), router))

}

func (srv *ApiServer) handleNotSupported(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{
		Value: "Command not supported",
	}

	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func (srv *ApiServer) prepareAlpacaResponse(r *http.Request, resp *alpacaResponse) {
	ctid := getClientTransactionId(r)
	if ctid < 0 {
		ctid = 0
	}
	srv.ServerTransactionID += 1
	resp.ClientTransactionID = uint32(ctid)
	resp.ServerTransactionID = srv.ServerTransactionID

}

func (srv *ApiServer) validAlpacaRequest(r *http.Request) bool {
	//if _, err := strconv.Atoi(c.Param("device_id")); err != nil {
	//	return false
	//}
	cid := getClientId(r)
	if cid < 0 {
		return false
	}

	ctidv := getClientTransactionId(r)
	return ctidv >= 0
}

func getClientId(r *http.Request) int {
	cidv := ""
	if r.Method == "GET" {
		cidv = r.URL.Query().Get("ClientID")
		if cidv == "" {
			cidv = r.URL.Query().Get("clientid")
			if cidv == "" {
				return -1
			}
		}
	} else {
		cidv = r.PostFormValue("ClientID")
		if cidv == "" {
			cidv = r.PostFormValue("clientid")
			if cidv == "" {
				return -1
			}
		}
	}
	cid, err := strconv.Atoi(cidv)
	if err != nil {
		return -1
	}
	if cid < 0 {
		return -1
	}
	return cid
}

func getClientTransactionId(r *http.Request) int {
	ctidv := ""
	if r.Method == "GET" {
		ctidv = r.URL.Query().Get("ClientTransactionID")
		if ctidv == "" {
			ctidv = r.URL.Query().Get("clienttransactionid")
			if ctidv == "" {
				return -1
			}
		}
	} else { // PUT
		ctidv = r.PostFormValue("ClientTransactionID")
		if ctidv == "" {
			ctidv = r.PostFormValue("clienttransactionid")
			if ctidv == "" {
				return -1
			}
		}
	}
	ctid, err := strconv.Atoi(ctidv)
	if err != nil {
		return -1
	}
	if ctid < 0 {
		return -1
	}
	return ctid
}

func getIdFromRequest(r *http.Request) (result int32, err error) {
	sid := ""
	if r.Method == "GET" {
		sid = r.URL.Query().Get("Id")
		if sid == "" {
			return -1, errors.New("id parameter missing")
		}
	} else { // i.e. PUT command
		//sid = c.Request.FormValue("Id")
		sid = r.PostFormValue("Id")
		if sid == "" {
			return -1, errors.New("id parameter missing")
		}
	}

	iid, err := strconv.ParseInt(sid, 10, 32)
	if err != nil {
		return -1, errors.New("id parameter not numeric")
	}
	result = int32(iid)
	if result < 0 {
		return -1, errors.New("id parameter out of range")
	}

	result += 1 // offset for shared device

	return result, nil
}

func getValueFromRequest(r *http.Request) (result int64, err error) {
	svalue := ""
	if r.Method == "PUT" {
		// svalue = c.Request.FormValue("Value")
		svalue = r.PostFormValue("Value")
		if svalue == "" {
			return -1, errors.New("value parameter missing")
		}
	} else {
		return -1, errors.New("expected PUT")
	}
	ivalue, err := strconv.ParseInt(svalue, 10, 64)
	if err != nil {
		return -1, errors.New("value parameter not numeric")
	}
	result = int64(ivalue)
	if result < 0 {
		return -1, errors.New("value parameter out of range")
	}
	//log.Println("Value:", result)
	return result, nil
}

func getPositionFromRequest(r *http.Request) (result int32, err error) {
	sposition := ""
	if r.Method != "PUT" {
		return -1, errors.New("expected PUT")
	}

	// svalue = c.Request.FormValue("Value")
	sposition = r.PostFormValue("Position")
	if sposition == "" {
		return -1, errors.New("position parameter missing")
	}

	ivalue, err := strconv.ParseInt(sposition, 10, 32)
	if err != nil {
		return -1, errors.New("position parameter not numeric")
	}
	result = int32(ivalue)
	if result == 0 {
		return -1, errors.New("position parameter out of range")
	}
	log.Println("Position:", result)
	return result, nil
}

func getSwitchStateFromRequest(r *http.Request) (bstate bool, err error) {
	sstate := ""
	if r.Method == "GET" {
		sstate = r.URL.Query().Get("State")
		if sstate == "" {
			err = errors.New("state parameter missing")
			return
		}
	} else { // e.g PUT command
		sstate = r.PostFormValue("State")
		if sstate == "" {
			err = errors.New("state parameter missing")
			return
		}
	}
	bstate, err = strconv.ParseBool(sstate)
	if err != nil {
		return
	}
	return
}

func getSwitchNameFromRequest(r *http.Request) (sname string, err error) {
	if r.Method == "GET" {
		sname = r.URL.Query().Get("Name")
		if sname == "" {
			err = errors.New("name parameter missing")
			return
		}
	} else { // e.g PUT command
		sname = r.PostFormValue("Name")
		if sname == "" {
			err = errors.New("name parameter missing")
			return
		}
	}
	return
}

func getConnectedFromRequest(r *http.Request) (connect bool, err error) {
	// PUT command
	connect, err = strconv.ParseBool(r.PostFormValue("Connected"))
	if err != nil {
		err = errors.New("connected parameter missing")
		return
	}
	return
}
