package folder1

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	DB "github.com/x/y/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"strings"
	"sync"
)

func scrapPosts(row int64 ,multiplier int64,wg1 *sync.WaitGroup){

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		//		colly.CacheDir("./forum_url_cache"),

	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	/*
	// Rotate proxies
	rp, err := proxy.RoundRobinProxySwitcher("http://gate.smartproxy.com:7000")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)
	*/

	// yeni Object'e Scrap database'inin  forum_url(collection) ve articles(collection) tablosunu tanımladık.
	collThread := DB.DBArticles()
	collPosts := DB.DBPosts()

	// Limit by 10 documents only
	opts := options.Find()
	opts.SetSkip(0+row*multiplier)
	opts.SetLimit(multiplier)
	cursor, err := collThread.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	// get a list of all returned documents and print them out
	//see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	var idMain interface{}
	var postTitle interface{}
	// Scrap HTML Settings
	c.OnHTML("li.postbitim", func(e *colly.HTMLElement) {
		// postTitle := e.ChildText("div.postrow h2.title")
		postDate := e.ChildText("span.date")
		postUserName := e.ChildText("div.username_container a.username")
		postHTML := e.ChildText("div.content")
		urlMain := e.Request.URL

		//******************************************************************
		//******************************************************************
		//Main ID çekilmesi için database'ye gidilir
		//Databasede linkimiz sorulur ve ordan mainID çekilir
		dbQuery, err := collThread.Find(context.TODO(),bson.M {"url" : urlMain})
		if err != nil { log.Fatal(err) }

		// The line in the database brought from the site is transferred to the array
		var newDbArray []bson.M
		if err = dbQuery.All(context.TODO(), &newDbArray); err != nil { log.Fatal(err) }
		if newDbArray != nil {
			//fmt.Println("Link located on DataBase............", newDbArray[0]["_id"])
			idMain  = newDbArray[0]["_id"]
		}else { /*fmt.Println("Link not located on DataBase.........XXX", urlMain)*/}

		//pulling thread_name #####################
		dbQuery2, err := collThread.Find(context.TODO(),bson.M {"url" : urlMain})
		if err != nil { log.Fatal(err) }

		// The line in the database brought from the site is transferred to the array
		var newDbArray2 []bson.M
		if err = dbQuery2.All(context.TODO(), &newDbArray2); err != nil { log.Fatal(err) }
		if newDbArray2 != nil {
			//fmt.Println("Link located on DataBase............", newDbArray[0]["_id"])
			postTitle  = newDbArray2[0]["thread_name"]
		}else { /*fmt.Println("Link not located on DataBase.........XXX", postTitle)*/}

		//******************************************************************
		//******************************************************************


		// Adding datas to database
		collPosts.InsertOne(context.TODO(), bson.D{
			{Key: "post_title", Value: postTitle},
			{Key: "date", Value: postDate},
			{Key: "user_name", Value: postUserName},
			{Key: "content", Value: postHTML},
			{Key: "main_id", Value: idMain},
		})

		fmt.Println("Post added on DataBase+++++++++++", postTitle)
		//---> e.Request.Visit(prevNext)
		// Adding datas to database End

	})
	// Next-Page-Button link
	c.OnHTML("div.below_postlist", func(h *colly.HTMLElement) {
		prevNext := h.ChildAttr("span.prev_next a[rel=next]", "href")
		//-->the last page-->  prevNext := h.ChildAttr("span.first_last a", "href")
		cleanLinkNext := strings.Split(prevNext, "?")  // Next Page
		if cleanLinkNext[0] == "" {
			fmt.Println("x===x This Thread Finished x===x")
		}else{
			fmt.Println("===> NextPage Visiting ===>",cleanLinkNext[0])
		}
		c.Visit(cleanLinkNext[0])
	})


	for _, result := range results {
		mainForumUrl := result["url"].(string)
		fmt.Println("##### Looking Forum Thread Url #####: ",mainForumUrl)
		c.Visit(mainForumUrl)
		c.Wait()
	}
	defer wg1.Done()
}

func PullPostsMultiTask() {
	var i int32
	var wg1 sync.WaitGroup


	results := DB.DbForumUrlMap()
	SiteInfo := DB.DbSiteInfoMap()

	// 0000 main Forum Link --> goRoutine
	dbTotalGoRoutine := SiteInfo[0]["number_of_go_routines"].(int32)

	totalLinkLen := len(results)
	multiplier := int64(math.Ceil(float64(totalLinkLen) /float64(dbTotalGoRoutine)))



	for i=0;i<dbTotalGoRoutine;i++ {
		wg1.Add(1)
		go scrapPosts(int64(i),multiplier,&wg1)
	}

	wg1.Wait()
	fmt.Println("#######################\n# PullPostsMultiTask  #\n####### THE END #######")
}
