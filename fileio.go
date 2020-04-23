package main

import (
	"encoding/json"
	"fileio/storage"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	gonanoid "github.com/matoous/go-nanoid"
)

var fs = http.FileServer(http.Dir("./web/static"))
var indexHTML, _ = ioutil.ReadFile("./web/views/index.html")
var stg storage.Storage

func main() {

	stg = storage.Storage{
		Type: "redis",
	}

	// start: fix this
	auth := map[string]string{
		"username": "",
		"password": "",
		"host":     "127.0.0.1:6379",
		"database": "0",
	}

	// username:password@127.0.0.1:6379/0
	if os.Getenv("DB_URI") != "" {
		connURI, err := url.Parse(os.Getenv("DB_URI"))
		if err != nil {
			panic(err)
		}

		stg.Type = connURI.Scheme

		auth["username"] = connURI.User.Username()
		auth["password"], _ = connURI.User.Password()
		auth["host"] = connURI.Host
		auth["database"] = strings.Replace(connURI.Path, "/", "", 1)
	}

	// host, username, password, database
	stg.Connect(auth["host"], auth["username"], auth["password"], auth["database"])
	// end: fix this

	// start listen
	fmt.Println(http.ListenAndServe(":8080", http.HandlerFunc(Index)))
}

// Index for testing
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		// Serve static file
		if len(r.URL.Path) > 7 {
			if r.URL.Path[:7] == "/static" {
				http.StripPrefix("/static", fs).ServeHTTP(w, r)
				return
			}
		}

		// show index.html
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexHTML)
			return
		}

		// has file id
		fileID := strings.ReplaceAll(r.URL.Path, "/", "")
		if len(fileID) == 6 {

			// get file from database
			bytes, err := stg.Get(fileID)
			if err == nil {
				// todo get mime-type by []byte file
				// w.Header().Set("Content-Type", "")
				w.Write(bytes)

				// Delete file from storage
				stg.Del(fileID)

				return
			}
		}

		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": false, "error": 404, "message": "file not found"} `))
		return
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
					"link":    "http://" + r.Host + "/" + id,
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
