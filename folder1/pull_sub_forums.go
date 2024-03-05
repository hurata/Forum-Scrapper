package folder1

import (
	"fmt"
	"github.com/gocolly/colly"
	DB "github.com/x/y/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"math"
	"strings"
	"sync"
)

func scrapSubForums(row int64, multiplier int64,wg1 *sync.WaitGroup){
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	//	colly.Debugger(&debug.LogDebugger{}),
	//	colly.CacheDir("./forum_url_cache"),
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





	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.
	collection := DB.DBForumUrl()
	//collectionSelector := DB.DBSiteInfo()


	// ################################ vv  collection
	opts := options.Find()
	// Limit by 10 documents only
	opts.SetSkip(0+row*multiplier)
	opts.SetLimit(multiplier)
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	//	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	// get a list of all returned documents and print them out
	//see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	// ################################ ^^


	SiteInfo := DB.DbSiteInfoMap()

	// selecting selectors in here
	dbSiteLink := SiteInfo[0]["site"].(string)
	dbMainDiv := SiteInfo[0]["site_main_div"].(string)
	dbForumLink := SiteInfo[0]["forum_link"].(string)
	dbForumLinkAttr := SiteInfo[0]["forum_link_attr"].(string)
	dbForumName := SiteInfo[0]["forum_name"].(string)
	dbForumDesc := SiteInfo[0]["forum_desc"].(string)


	var idMainForum interface{}
	// Scrap HTML Settings
	c.OnHTML(dbMainDiv, func(e *colly.HTMLElement) {
		link :=  e.ChildAttr(dbForumLink,dbForumLinkAttr)
		forumName := e.ChildText(dbForumName)
		forumDesc := e.ChildText(dbForumDesc)
		urlMainForum := e.Request.URL

		if link == ""{
			a := fmt.Sprintf("%s",urlMainForum)
			fmt.Println("||||||||||:",a)
			c.Visit(a)
		}


		// link, ?s 'den arındırılır
		base := dbSiteLink
		base2 := strings.Split(base, "://")
		base3 := base2[1]  			//   www.example.com
		//fmt.Println("----------------->",base3)

		i := strings.Index(link, base3)		// kaçıncı karakterden sonra söylenen karakter geliyor?
		if i > -1 {
			cleanLinkForum := strings.Split(link, base3)
			linkTwo := base + cleanLinkForum[1]

			cleanLink := strings.Split(linkTwo, "?")
			cleanLink2 := cleanLink[0]

			//**********************
			//**********************
			//Main ID çekilmesi için database'ye gidilir
			dbQuery, err := collection.Find(context.TODO(),bson.M {"url" : urlMainForum})
			if err != nil {
				log.Fatal(err)
			}
			// The line in the database brought from the site is transferred to the array
			var newDbArray []bson.M
			if err = dbQuery.All(context.TODO(), &newDbArray); err != nil {
				log.Fatal(err)
			}
			if newDbArray != nil {
				//	fmt.Println("Link located on DataBase............", newDbArray[0]["_id"])
				idMainForum  = newDbArray[0]["_id"]

			}else {
				//	fmt.Println("Link not located on DataBase.........XXX", urlMain)
			}
			//**********************
			//**********************



			// siteden getirilen link Database'de var olan
			// linkler arasında sorgulanır.
			selectResult, err := collection.Find(context.TODO(),bson.M {"url" : cleanLink2})
			if err != nil {
				log.Fatal(err)
			}

			// siteden getirilen linkin database'deki satırı diziye aktarılır
			var urlFiltered []bson.M
			if err = selectResult.All(context.TODO(), &urlFiltered); err != nil {
				log.Fatal(err)
			}
			if urlFiltered != nil {
				fmt.Println("xxxxx Sub-Forum located on DataBase==>",cleanLink2)
			}
			if urlFiltered == nil {

				// Adding datas to database
				insertResult, err := collection.InsertOne(context.TODO(), bson.D{
					{Key: "url", Value: cleanLink2},
					{Key: "forum_name", Value: forumName},
					{Key: "forum_desc", Value: forumDesc},
					{Key: "sub_forum", Value: 1},
					{Key: "main_id", Value: idMainForum},
					//	{Key: "last_page", Value: ""},
				})
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("----- Added Sub-Forum ID ==>",insertResult.InsertedID)
				//e.Request.Visit(link)
				// Adding datas to database End
			}
		} else {
		fmt.Println("Index not supported")
		fmt.Println(link)
		}
	})

	//Set error handler
	c.OnError(func(r *colly.Response, err error) {
		a := fmt.Sprintf("%s",r.Request.URL)
		b := fmt.Sprintf("%s", err.Error())
		fmt.Println("Error:",b)

		if b != "Forbidden"{
		c.Visit(a)
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
	defer wg1.Done()
}

func PullSubForums() {
	var i int32
	var wg1 sync.WaitGroup

	// 126 main Forum Link will searching and pull the sub-forums

	//collection := DB.DBForumUrl()
	//collSelector := DB.DBSiteInfo()
	SiteInfo := DB.DbSiteInfoMap()
	results := DB.DbForumUrlMap()


	// 126 main Forum Link --> goRoutine
	dbTotalGoRoutine := SiteInfo[0]["number_of_go_routines"].(int32)

	totalLinkLen := len(results)
	multiplier := int64(math.Ceil(float64(totalLinkLen) /float64(dbTotalGoRoutine)))



	for i=0;i<dbTotalGoRoutine;i++ {
		wg1.Add(1)
		go scrapSubForums(int64(i),multiplier, &wg1)
	//	time.Sleep(100 * time.Millisecond)
	}

	wg1.Wait()
	fmt.Println("#######################\n#### PullSubForums ####\n####### THE END #######")
}