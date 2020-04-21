package main

import (
	"encoding/json"
	"fileio/storage"
	"fmt"
	"io/ioutil"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid"
)

var indexHTML, _ = ioutil.ReadFile("./web/views/index.html")
var stg storage.Storage

func main() {

	stg = storage.Storage{
		Type: "redis",
	}

	// host, username, password, database
	stg.Connect("127.0.0.1:6379", "", "", "0")

	http.HandleFunc("/", Index)

	// serve static file
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// start listen
	fmt.Println(http.ListenAndServe(":8080", nil))
}

// Index for testing
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fileID := r.URL.Query().Get("id")
		if fileID == "" {
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexHTML)
			return
		} else {

			bytes, err := stg.Get(fileID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(404)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success": false, "error": 404, "message": "file not found"} `))
				return
			}

			// todo get mime-type by []byte file
			// w.Header().Set("Content-Type", "")
			w.Write(bytes)

			// Delete file from storage
			stg.Del(fileID)

			return
		}
	}
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20)
		// FormFile returns the first file for the given key `file`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		id, err := gonanoid.Generate("AiUeO69", 6)
		if err == nil {
			err = stg.Set(id, fileBytes, 1800)
			if err == nil {
				data := map[string]interface{}{
					"success": true,
					"key":     id,
					"link":    "",
					"expiry":  "30 minutes",
					"sec_exp": 1800,
				}
				resp, _ := json.Marshal(data)
				w.Write(resp)
				return
			}
		}
		w.Write([]byte(`{"success": false, "error": 402, "message": "failing"} `))
		return
	}
}
