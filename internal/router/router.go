package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-task/config"
	"test-task/internal/auth"
	"test-task/internal/auth/tokens"

	"github.com/gorilla/mux"
)

var router *mux.Router

func getRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/refresh", RefreshHandler).Methods("POST")
	router.HandleFunc("/protected", ProtectedHandler).Methods("GET")

	return router
}

func Run() error {
	router := getRouter()

	fmt.Println("Starting the server")
	err := http.ListenAndServe(
		fmt.Sprintf("%v:%v", config.ServerHost, config.ServerPort),
		router,
	)
	if err != nil { return err }
	
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginData auth.LoginData
	json.NewDecoder(r.Body).Decode(&loginData)

	err := tokens.CheckGuid(loginData.GUID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, err)
		return
	}

	tokensJson, err := tokens.CreateTokenPair(loginData.GUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal error")
		return
	}

	w.Write(tokensJson)
	fmt.Fprint(w)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	accessToken = accessToken[len("Bearer "):]

	var accessPayload auth.AccessPayload
	json.NewDecoder(r.Body).Decode(&accessPayload)

	err := tokens.VerifyAccess(accessPayload, accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var refreshData auth.RefreshData
	json.NewDecoder(r.Body).Decode(&refreshData)

	w.Header().Set("Content-Type", "application/json")
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	accessToken = accessToken[len("Bearer "):]

	tokensJson, err := tokens.VerifyRefresh(refreshData)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	w.Write(tokensJson)
	fmt.Fprint(w)
}