package main

import (
	"context"
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

	collection := client.Database("scrap").Collection("site_info")

	collection.InsertOne(context.TODO(), bson.D{
		{Key: "site", Value: "https://www.ownedcore.com/forums/"},
		{Key: "site_main_div", Value: "div.forumrow"},
		{Key: "forum_link", Value: "h2.forumtitle a"},
		{Key: "forum_link_attr", Value: "href"},
		{Key: "forum_name", Value: "h2.forumtitle a"},
		{Key: "forum_desc", Value: ".forumdescription"},
		{Key: "forum_lastpage_div", Value: "div.below_threadlist"},
		{Key: "forum_lastpage", Value: "span.first_last a"},
		{Key: "forum_lastpage_attr", Value: "href"},
		{Key: "thread_main_div", Value: "div.threadlist ol.threads div.threadinfo"},
		{Key: "thread_link", Value: "h3.threadtitle a"},
		{Key: "thread_link_attr", Value: "href"},
		{Key: "thread_name", Value: "h3.threadtitle a.title"},
		{Key: "number_of_go_routines", Value: 4},
		{Key: "sub_forum_status", Value: true},
	})
}