package handlers

import (
	"net/http"
	"path/filepath"
	"product_img/product_images/files"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage // this is of type interface
}

func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// serveHTTP implements the http.Handler interface as studied before
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // reading the request
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn) // hclog package used for logging instead of log.logger here

	// no need to check for invalid id or filename as the mux router will not
	// sends requests here unless they have the correct parameters
	// if id == "" || fn == "" {
	// 	f.invalidURI(r.URL.String(), rw)
	// 	return
	// }

	f.saveFile(id, fn, rw, r)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("invalid path", "path", uri)
	http.Error(rw, "invaild file path should be in format /[id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) { // id and path is of typestring
	f.log.Info("save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)   //getting relative path
	err := f.store.Save(fp, r.Body) //calling the storing interface
	// assumption says if any Save method is written, it gets called here
	// the method is implemented in storage interface
	// hence we dont care what the internal components of save are
	// which can be either hard drive or ssd or cloud, anything
	// here we see implementation hiding

	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}

}
