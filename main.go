package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Cache is a simple abstraction how to cache a request
type Cache interface {
	RequestInCache(*http.Request) bool
	GetDataFromCache(http.ResponseWriter, *http.Request)
	SendRequestToService(*http.Request) *http.Response
	StoreRespond(*http.Request, io.Reader)
}

var cache Cache

func main() {
	ic := &ImaginaryCache{
		rootPath:          Conf.CacheRoot,
		imaginaryHostPort: Conf.ImaginaryHostPort,
	}
	cache = ic
	router := mux.NewRouter()
	router.HandleFunc(`/{type}`, GetRequest).Methods("GET")
	log.Println("Listen at port: " + Conf.ServerPort)
	log.Println("Imaginary server: " + Conf.ImaginaryHostPort)
	log.Fatal(http.ListenAndServe(Conf.ServerPort, router))
}

// GetRequest is the handler for all GET methods.
func GetRequest(w http.ResponseWriter, r *http.Request) {

	if vars := mux.Vars(r); vars["type"] == "fullsize" {
		f, err := os.Open(filepath.Join(Conf.FullsizeRoot, r.FormValue("file")))
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}
		_, err = io.Copy(w, f)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	if cache.RequestInCache(r) {
		cache.GetDataFromCache(w, r)
	} else {
		resp := cache.SendRequestToService(r)
		defer resp.Body.Close()
		var buf bytes.Buffer
		rBody := io.TeeReader(resp.Body, &buf)
		_, err := io.Copy(w, rBody)
		if err != nil {
			log.Fatalln(err)
		}
		cache.StoreRespond(r, &buf)
	}
}
