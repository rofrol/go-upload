// http://sanatgersappa.blogspot.ie/2013/03/handling-multiple-file-uploads-in-go.html
package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		if false {
			hah, err := ioutil.ReadAll(r.Body);
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Printf("%v", string(hah))
			return
		}

		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			//fmt.Printf("%+v\n", part) // Error
			//fmt.Printf("%+v\n", part.Header) // Error
			fmt.Printf("%+v\n", part.FormName())
			//fmt.Println(part.FormName())

			if part.FileName() == "" { // normal text field
				b, err := ioutil.ReadAll(part)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Println(string(b))
			} else {
				dst, err := os.Create("assets/" + part.FileName())
				defer dst.Close()

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if _, err := io.Copy(dst, part); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		//display success message.
		display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	//static file handler.
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	fs := JustFilesFilesystem{http.Dir("assets/")}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(fs)))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
