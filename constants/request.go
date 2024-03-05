package constants

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func DBForumUrl() *mongo.Collection {
	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {	log.Fatal(err)	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {	log.Fatal(err)	}

	collection := client.Database("scrap").Collection("forum_url")
	return collection
}

func DBSiteInfo() *mongo.Collection {
	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {	log.Fatal(err)	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {	log.Fatal(err)	}

	collection := client.Database("scrap").Collection("site_info")
	return collection
}

func DBArticles() *mongo.Collection {
	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {	log.Fatal(err)	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {	log.Fatal(err)	}

	collection := client.Database("scrap").Collection("articles")
	return collection
}

func DBPosts() *mongo.Collection {
	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {	log.Fatal(err)	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {	log.Fatal(err)	}

	collection := client.Database("scrap").Collection("posts")
	return collection
}

func DbSiteInfoMap() []bson.M {

	collSelector := DBSiteInfo()

	// ################################ vv  collectionSelector
	cursor, err := collSelector.Find(context.TODO(), bson.D{{}})

	if err != nil {	log.Fatal(err)	}
	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	var resultsSelector []bson.M
	if err = cursor.All(context.TODO(), &resultsSelector); err != nil {
		log.Fatal(err)
	}
	// ################################ ^^
	return resultsSelector
}

func DbForumUrlMap() []bson.M {

	collection := DBForumUrl()

	// ################################ vv  collection
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	//	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	// get a list of all returned documents and print them out
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	// ################################ ^^

	return results
}
