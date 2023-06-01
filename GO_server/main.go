package main

import (
	"fmt"
	"log"
	"net/http"
)

//Writing the hello handler function

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" { // checking if path is not hitting the "/hello" then we are throeing error as 404
		http.Error(w, "404 error not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" { //"GET" is nothing but "hello" we are posting
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, " Hello world!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { //passsing string layout as an argument
		fmt.Fprintf(w, " Parseform() error: %v", err)
	}

	fmt.Fprintf(w, " POST request successful!\n")

	name := r.FormValue("name")
	address := r.FormValue("address")
	contact := r.FormValue("contact")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	status := r.FormValue("status")

	fmt.Fprintf(w, "Name =%s\n", name)
	fmt.Fprintf(w, "Address =%s\n", address)
	fmt.Fprintf(w, "Age =%s\n", age)
	fmt.Fprintf(w, "Contact =%s\n", contact)
	fmt.Fprintf(w, "Gender =%s\n", gender)
	fmt.Fprintf(w, "Status =%s\n", status)

}

func main() {
	fileserver := http.FileServer(http.Dir("./static")) //going to read the index.html from static dir
	//("./static"))
	http.Handle("/", fileserver)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting theserver ion 8080 port\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("let's see what is the error ", err)
	}

}
