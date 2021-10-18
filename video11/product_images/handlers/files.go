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

//creating a processing function for gorillamux, as a servehttp function we saw before
func (f *Files) RetardFunction(rw http.ResponseWriter, r *http.Request) {
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

//UploadMultipart something
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024) //128kB
	if err != nil {
		f.log.Error("bad request", err)
		http.Error(rw, "expected multipart data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	f.log.Info("process form for id,", id)
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
