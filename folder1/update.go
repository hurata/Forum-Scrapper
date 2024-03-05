package folder1

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
	"time"
)

func updateThreads(forumUrl string,wg2 *sync.WaitGroup){
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		//	colly.CacheDir("./forum_url_cache"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	stop := false

	// Rotate proxies
	rp, err := proxy.RoundRobinProxySwitcher("http://gate.smartproxy.com:7000")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	//  MongoDB Settings
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {	log.Fatal(err) }
	err = client.Ping(context.TODO(), nil)
	if err != nil {	log.Fatal(err) }

	// yeni Object'e Scrap database'inin  forum_url(collection) ve articles(collection) tablosunu tanımladık.
	collUrls := client.Database("scrap").Collection("forum_url")
	collThread := client.Database("scrap").Collection("articles")

	var idMain interface{}
	// Scrap HTML Settings
	c.OnHTML("div.threadlist ol.threads div.threadinfo", func(e *colly.HTMLElement) {
		link := e.ChildAttr("h3.threadtitle a", "href")
		threadName := e.ChildText("h3.threadtitle a.title")
		urlMain := e.Request.URL.String()
		cleanLinkMain := strings.Split(urlMain, "index")
		urlMain2 :=cleanLinkMain[0]

		// link, arındırılır
		base := "https://www.ownedcore.com/forums/"
		cleanLinkForum := strings.Split(link, "forums/")
		linkTwo := base + cleanLinkForum[1]
		cleanLink := strings.Split(linkTwo, "?")

		//**********************
		//**********************
		//Main ID çekilmesi için database'ye gidilir
		//Databasede linkimiz sorulur ve ordan mainID çekilir
		dbQuery, err := collUrls.Find(context.TODO(), bson.M{"url": urlMain2})
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
			idMain = newDbArray[0]["_id"]

		} else {
			//	fmt.Println("Link not located on DataBase.........XXX", urlMain)
		}
		//**********************
		//**********************

		if !stop{
			// siteden getirilen link Database'de var olan linkler arasında "sorgulanır".
			// The link brought from the site is queried among the existing links in the Database.
			selectResult, err := collThread.Find(context.TODO(), bson.M{"url": cleanLink[0]})
			if err != nil {
				log.Fatal(err)
			}

			// siteden getirilen linkin database'deki satırı "diziye aktarılır"
			// The line in the database brought from the site is transferred to the array
			var urlFiltered []bson.M
			if err = selectResult.All(context.TODO(), &urlFiltered); err != nil {
				log.Fatal(err)
			}
			if urlFiltered != nil {
				fmt.Println("This thread located on DataBase............", threadName)
				stop = true
				//fmt.Println("This thread located on DataBase............", cleanLink[0])
			}
			if urlFiltered == nil {
				// Adding datas to database
				collThread.InsertOne(context.TODO(), bson.D{
					{Key: "url", Value: cleanLink[0]},
					{Key: "thread_name", Value: threadName},
					{Key: "main_id", Value: idMain},
				})

				fmt.Println("Thread added on DataBase+++++++++++", threadName)
				//---> e.Request.Visit(prevNext)
				// Adding datas to database End
			}
		}

	})

	// Next-Page-Button link
	c.OnHTML("div.below_threadlist", func(h *colly.HTMLElement) {
		prevNext := h.ChildAttr("span.prev_next a[rel=next]", "href")
		//-->the last page-->  prevNext := h.ChildAttr("span.first_last a", "href")
		cleanLinkNext := strings.Split(prevNext, "?") // Next Page

		if !stop {
			if cleanLinkNext[0] == "" {
				fmt.Println("x===x This Forum Finished x===x")
			} else {
				fmt.Println("===> NextPage Visiting ===>", cleanLinkNext[0])
			}

			c.Visit(cleanLinkNext[0])
		}

	})

	//Set error handler
	c.OnError(func(r *colly.Response, err error) {
		a := fmt.Sprintf("%s",r.Request.URL)
		b := fmt.Sprintf("%s", err.Error())
		fmt.Println("Error:",b)

		if b != "Forbidden" && b != "Not Found"{
			c.Visit(a)
		}
	})

	c.Visit(forumUrl)
	c.Wait()

	defer wg2.Done()
}

func databasePull2(row int64 ,multiplier int64, wg1 *sync.WaitGroup)  {
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
	// yeni Object'e Scrap database'inin  forum_url(collection) ve articles(collection) tablosunu tanımladık.
	collUrls := client.Database("scrap").Collection("forum_url")

	opts := options.Find()
	// Limit by 10 documents only
	opts.SetSkip(0+row*multiplier)
	opts.SetLimit(multiplier)
	cursor, err := collUrls.Find(context.TODO(), bson.D{{}}, opts)

	//	cursor, err := collUrls.Find(context.TODO(), bson.D{{}})
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

		mainForumUrl := result["url"].(string)
		fmt.Println("##### Looking Main or Sub Forum Url #####:",mainForumUrl)

		var wg2 sync.WaitGroup
		wg2.Add(1)
		go updateThreads(mainForumUrl,&wg2)

		time.Sleep(time.Millisecond)
		wg2.Wait()

	}

	defer wg1.Done()
}

func Update() {
	var i,pullForumLink int64
	var wg1 sync.WaitGroup

	// totalLinks = 338
	pullForumLink = 338	// 338 thread for each forum link pull from DataBase
						// 338 Total Thread

	//multiplier := int64(math.Ceil(float64(totalLinks)/float64(pullForumLink)))
	multiplier := int64(1)
	for i=0; i<pullForumLink ; i++ {
		wg1.Add(1)
		go databasePull2(i,multiplier,&wg1)
	}

	wg1.Wait()
	fmt.Println("#######################\n######  UPDATE  #######\n####### THE END #######")
}