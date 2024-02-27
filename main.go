package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	authenticationapi "k8s.io/api/authentication/v1beta1"
)

func authN(tr *authenticationapi.TokenReview) {
	tr.Status.Authenticated = true
	tr.Status.User = authenticationapi.UserInfo{
		Username: "mock",
		UID:      "mock",
		Groups:   []string{"group-mock"},
		Extra:    nil,
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving %s", r.Method)

	if r.Method != "POST" {
		http.Error(w, "Only accept POST requests", http.StatusMethodNotAllowed)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Receiving Token: %s\n", string(payload))

	// Unmarshal JSON from POST request to TokenReview object
	tr := &authenticationapi.TokenReview{}
	err = json.Unmarshal(payload, tr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authN(tr)

	// Marshal the TokenReview to JSON and send it back
	result, err := json.Marshal(*tr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
	w.Header().Set("Content-Type", "application/json")
}

func main() {

	// Set up a /authenticate resource handler
	http.HandleFunc("/authenticate", helloHandler)

	// Listen to port 443 and wait
	log.Println("Listening on port 443 for requests...")
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}
