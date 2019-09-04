package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

// App encapsulates Env,Router and middleware
type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL		string `json:"url" validate:"nonzero"`
	ExpirationInMinutes		int64 `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	shortlink 	string `json:"shortlink"`
}


// Initialize is initialization of app
func (a *App) Initialize() {
	//set log formatter
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes()  {
	a.Router.HandleFunc("/api/shorten",a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info",a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}",a.redirect).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter,r *http.Request)  {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req);err != nil {


		respondWithError(w,StatusError{http.StatusBadRequest,
		fmt.Errorf("parse parameters failed %v",r.Body)})

		return
	}

	if err := validator.Validate(req);err != nil{
		respondWithError(w,StatusError{http.StatusBadRequest,
			fmt.Errorf("validate parameters failed %v",req)})
		return
	}

	defer r.Body.Close()

	fmt.Printf("%v\n",req)
}

func (a *App) getShortlinkInfo(w http.ResponseWriter,r *http.Request)  {

	vals := r.URL.Query()

	s := vals.Get("shortlink")

	fmt.Printf("%s\n",s)

}

func (a *App) redirect(w http.ResponseWriter,r *http.Request)  {
	vars := mux.Vars(r)
	fmt.Printf("%s\n",vars["shortlink"])
}

//Run starts listen and server
func (a *App) Run(addr string)  {
	log.Fatal(http.ListenAndServe(addr,a.Router))
}

func respondWithError(w http.ResponseWriter,err error)  {

	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s",e.Status(),e)
		respondWithJson(w,e.Status(),e.Error())
	default:
		log.Printf("err: %s",e)
		respondWithJson(w,http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))

	}

}

func respondWithJson(w http.ResponseWriter,code int,payload interface{})  {
	resp,_ := json.Marshal(payload)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(resp)
}



