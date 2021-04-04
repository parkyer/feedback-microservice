package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Queja struct { //datos
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Queja_user   string             `json:"queja_user,omitempty" bson:"queja_user,omitempty"`
	Calificacion float64            `json:"calificacion,omitempty" bson:"calificacion,omitempty" `
	ID_Parkyer   int                `json:"id_parkyer,omitempty" bson:"id_parkyer,omitempty" `
}

func main() {
	fmt.Println("Iniciando ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:rootmaster@clusterprueba.0qn4w.mongodb.net/quejasdb?retryWrites=true&w=majority") //URI
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/queja", CreateQuejaEndpoint).Methods("POST") //rutas
	router.HandleFunc("/quejas", GetQuejasEndpoint).Methods("GET")
	router.HandleFunc("/queja/{id}", GetQuejaEndpoint).Methods("GET")
	http.ListenAndServe(":4000", router) //puerto
}

func CreateQuejaEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var queja Queja
	_ = json.NewDecoder(request.Body).Decode(&queja) //apuntador
	collection := client.Database("quejasdb").Collection("queja")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, queja)
	json.NewEncoder(response).Encode(result)
	//fmt.Println(request.Body)
}

func GetQuejaEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var queja Queja
	collection := client.Database("quejasdb").Collection("queja")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Queja{ID: id}).Decode(&queja)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(queja)
}
func GetQuejasEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var quejas []Queja
	collection := client.Database("quejasdb").Collection("queja")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var quejasV Queja
		cursor.Decode(&quejasV)
		quejas = append(quejas, quejasV)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(quejas)
}
