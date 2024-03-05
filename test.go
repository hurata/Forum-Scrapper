package main

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	DB "github.com/x/y/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

func main() {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		//	colly.Debugger(&debug.LogDebugger{}),
		//	colly.CacheDir("./forum_url_cache"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})


	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("scrap").Collection("forum_url")

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

	SiteInfo := DB.DbSiteInfoMap()

	//dbLastPageDiv := SiteInfo[0]["forum_lastpage_div"].(string)
	//dbLastPage := SiteInfo[0]["forum_lastpage"].(string)
	dbLastPageAttr := SiteInfo[0]["forum_lastpage_attr"].(string)
	//dbLastPageNumber := SiteInfo[0]["forum_lastpage_number"].(string)



	counterExtra :=0
	c.OnHTML("div.pagination div.right",func(l *colly.HTMLElement) {
		urlMain := l.Request.URL

		prevNextNumber2 := l.ChildAttr("ul li a.", dbLastPageAttr)
		//fmt.Println("#############",prevNextNumber2,"#############",prevNextNumber2)
		counterExtra++
		cleanLinkForum2 := strings.Split(prevNextNumber2, "/index")
		fmt.Println("=============",cleanLinkForum2[0],"=============",counterExtra)
		if cleanLinkForum2[1] != "" {
			lastPageString := cleanLinkForum2[1]
			cleanLinkLast := strings.Split(lastPageString, ".html")
			lastPage := cleanLinkLast[0]


			filterUpdate := bson.M{"url": bson.M{"$eq": urlMain}}

			update := bson.M{
				"$set": bson.M{
					"last_page": lastPage,
				},
			}

			collection.UpdateOne(context.TODO(),filterUpdate, update)
			fmt.Println("----- Added Sub-Forum LastPage ==>",lastPage)
			//e.Request.Visit(link)
			// Adding datas to database End
		}

	})


	// dataBase'den alınıp MainForum linkine gidilecek
	counter := 0
	for _, result := range results {
		counter++
		mainForumUrl := result["url"].(string)
		fmt.Println("##### Looking Main Forum Url #####",counter,":",mainForumUrl)
		c.Visit(mainForumUrl)
		c.Wait()
	}

}