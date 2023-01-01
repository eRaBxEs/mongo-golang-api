package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erabxes/mongo-golang/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{c}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !primitive.IsValidObjectID(id) {
		w.WriteHeader(http.StatusNotFound)
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
	}

	u := models.User{}

	ctx := context.TODO()
	collection := uc.client.Database("testdb").Collection("users")
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		w.WriteHeader(404)
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.Id = primitive.NewObjectID()

	ctx := context.TODO()
	collection := uc.client.Database("testdb").Collection("users")
	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		w.WriteHeader(404)
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !primitive.IsValidObjectID(id) {
		w.WriteHeader(404)
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
	}

	collection := uc.client.Database("tedstdb").Collection("users")
	ctx := context.TODO()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		w.WriteHeader(404)
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "Deleted User: ", oid)

}
