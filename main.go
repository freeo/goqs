/**
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// [START all]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	gorillajson "github.com/gorilla/rpc/json"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
}

var users = []User{
	User{ID: 1, FirstName: "Max", LastName: "Mustermann"},
	User{ID: 2, FirstName: "Erika", LastName: "Mustermann"},
	User{ID: 3, FirstName: "Markus", LastName: "Mustermann"},
	User{ID: 4, FirstName: "Ralf", LastName: "Schmitz"},
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	json.NewEncoder(w).Encode(users)
}

// hello responds to the request with a plain-text "Hello, world" message.
func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	json.NewEncoder(w).Encode(users)
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	params := mux.Vars(r)
	var id, x = strconv.Atoi(params["id"])
	log.Print(x)
	for _, item := range users {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = len(users) + 1
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
	log.Print("Received:", user)
}

type Task struct {
	Description string
}

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "PATH_TO_THE_JSON_FILE")

	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv

	}

	// register hello function to handle all requests
	server := mux.NewRouter().StrictSlash(true)
	server.HandleFunc("/", getInfo).Methods("GET")
	server.HandleFunc("/v1/users", getUsers).Methods("GET")
	server.HandleFunc("/v1/users/{id}", getDetails).Methods("GET")
	server.HandleFunc("/v1/users", postUsers).Methods("POST")

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)

	// simple go.mod test, any package, so that I can test the pipeline
	s := rpc.NewServer()
	s.RegisterCodec(gorillajson.NewCodec(), "application/json")
	s.RegisterService(new(HelloService), "")
}

type HelloService struct{}

func datastoreAccess() {
	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "dmk-hacking"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the kind for the new entity.
	kind := "Task"
	// Sets the name/ID for the new entity.
	name := "sampletask1"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	// Creates a Task instance.
	task := Task{
		Description: "Buy milk",
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}

	fmt.Printf("Saved %v: %v\n", taskKey, task.Description)
}

// [END all]
