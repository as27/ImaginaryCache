package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// ImaginaryCache is the implementation of the Cache interface
type ImaginaryCache struct {
	rootPath          string
	imaginaryHostPort string
}

// RequestInCache is the implementation of the Cache interface
func (c *ImaginaryCache) RequestInCache(r *http.Request) bool {
	fp := filepath.Join(c.rootPath, c.makePath(r))
	_, err := os.Stat(fp)
	return os.IsExist(err)
}

// GetDataFromCache is the implementation of the Cache interface
func (c *ImaginaryCache) GetDataFromCache(w http.ResponseWriter, r *http.Request) {
	if c.RequestInCache(r) == false {
		return
	}
	f, err := os.Open(c.makePath(r))
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(w, f)
	if err != nil {
		log.Fatalln(err)
	}

}

// SendRequestToService is the implementation of the Cache interface
func (c *ImaginaryCache) SendRequestToService(r *http.Request) *http.Response {
	u := r.URL
	u.Host = c.imaginaryHostPort
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

// StoreRespond is the implementation of the Cache interface
func (c *ImaginaryCache) StoreRespond(r *http.Request, respBodyBuffer io.Reader) {
	fp := c.makePath(r)
	os.MkdirAll(filepath.Dir(fp), 0777)
	f, err := os.Create(fp)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(f, respBodyBuffer)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *ImaginaryCache) makePath(r *http.Request) string {
	strToHash := ""
	vars := mux.Vars(r)
	strToHash += vars["type"]
	strToHash += vars["file"]
	strToHash += vars["width"]
	strToHash += vars["height"]
	h := sha1.New()
	io.WriteString(h, strToHash)
	hstr := fmt.Sprintf("%x", h.Sum(nil))
	p := filepath.Join(string(hstr[:2]), hstr)
	return p
}
