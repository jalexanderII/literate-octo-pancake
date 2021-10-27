package example_code_snippets

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func name() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Oops", http.StatusBadRequest)
			return
		}
		_, err = fmt.Fprintf(w, "Hello %s", d)
		if err != nil {
			return
		}
	})

	http.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request){
		log.Println("Goodbye!")
	})

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
