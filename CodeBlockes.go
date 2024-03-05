package main

import (
	"context"
//	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {

	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.
	collection := client.Database("scrap").Collection("forum_url")

	// *********************************************************
	// Pulling data from DataBase with an array
	// dataBase'den url çekilir...
	// *********************************************************
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	// get a list of all returned documents and print them out
	//see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {

		mainForumUrl := result["url"]
		fmt.Println(mainForumUrl)
	}
	// *********************************************************
	// Pulling data END
	// *********************************************************
}