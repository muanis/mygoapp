package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"log"
)

type Cat struct {
	Name string `json:"name"`
	Breed string `json:"breed"`
}

func main() {

	http.HandleFunc("/weighcats", func(w http.ResponseWriter, r *http.Request) {
		var cat Cat

		err := json.NewDecoder(r.Body).Decode(&cat)
		if err != nil {
			http.Error(w, "bad cat", 400)
			return
		}

		log.Printf("we have a cat! %+v", cat)

		weight := len(fmt.Sprint(cat.Name, cat.Breed))

		fmt.Fprint(w , weight)
	})

	http.Handle("/", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8000",nil)
}
