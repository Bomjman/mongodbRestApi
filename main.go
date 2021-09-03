package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest-api/helper"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("Starting the application...")
	r := mux.NewRouter()

	r.HandleFunc("/profile/{id}", getProfile).Methods("GET")
	r.HandleFunc("/profiles", getProfiles).Methods("GET")

	log.Fatal(http.ListenAndServe(":8888", r))
}

type Profiles struct {
	Skill    string
	Profile  string
	Gender   string
	Location string
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	var profile bson.M
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := (params["id"])
	collection := helper.Connect()
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Profiles{Profile: id}).Decode(&profile)
	if err == mongo.ErrNoDocuments {
		fmt.Println("record does not exist")
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
	}
	json.NewEncoder(w).Encode(profile)
}

func getProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	//var profile Profiles
	var profiles []bson.M
	collection := helper.Connect()

	query := Profiles{
		Location: r.Header.Get("location"),
		Gender:   r.Header.Get("gender"),
		Skill:    r.Header.Get("skill"),
	}

	// filter for location
	var location bson.M

	if query.Location != "" {
		location = bson.M{"$eq": query.Location}
	} else {
		location = bson.M{"$ne": ""}
	}
	//filter for gender
	var gender bson.M

	if query.Gender != "" {
		gender = bson.M{"$eq": query.Gender}
	} else {
		gender = bson.M{"$ne": ""}
	}
	// filter for skill
	var skill bson.M

	if query.Skill != "" {
		skill = bson.M{"$in": bson.A{query.Skill}}
	} else {
		skill = bson.M{"$elemMatch": bson.M{"$ne": ""}}
	}

	filter := bson.M{"$and": bson.A{
		bson.M{"location_locality": location},
		bson.M{"gender": gender},
		bson.M{"skills": skill}}}
	fmt.Println(filter)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}
	for cur.Next(ctx) {
		var r bson.M
		cur.Decode(&r)
		profiles = append(profiles, r)
	}
	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(profiles) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Not Found"))
		return
	}

	json.NewEncoder(w).Encode(profiles)
}
