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
	"time"
)

func scrapThreads(routine int,jobs int,workCell int,forumUrl string,totalpages int, wg2 *sync.WaitGroup){
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
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


	// yeni Object'e Scrap database'inin  forum_url(collection) ve articles(collection) tablosunu tanımladık.
	collUrls := DB.DBForumUrl()
	collThread := DB.DBArticles()

	SiteInfo := DB.DbSiteInfoMap()
	// selecting selectors in here
	dbSiteLink := SiteInfo[0]["site"].(string)
	dbThMainDiv:= SiteInfo[0]["thread_main_div"].(string)
	dbThreadLink := SiteInfo[0]["thread_link"].(string)
	dbThreadLinkAttr := SiteInfo[0]["thread_link_attr"].(string)
	dbThreadName := SiteInfo[0]["thread_name"].(string)

	var idMain interface{}
	// Scrap HTML Settings
	c.OnHTML(dbThMainDiv, func(e *colly.HTMLElement) {
		link := e.ChildAttr(dbThreadLink, dbThreadLinkAttr)
		threadName := e.ChildText(dbThreadName)
		urlMain := e.Request.URL.String()

		cleanLinkMain := strings.Split(urlMain, "index")
		urlMain2 :=cleanLinkMain[0]

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

			// Adding datas to database
			collThread.InsertOne(context.TODO(), bson.D{
				{Key: "url", Value: cleanLink2},
				{Key: "thread_name", Value: threadName},
				{Key: "main_id", Value: idMain},
			})
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

		if b != "Forbidden" && b != "Not Found"{
			c.Visit(a)
		}
	})

	go func() {
		for k := 1; k <= jobs; k++ {
			page := (workCell)*(k-1)+routine
			s := fmt.Sprintf("%d", page) //int to string

			mainForumUrlPage :=forumUrl+"index"+s+".html"
			if page <= totalpages {
				fmt.Println("##### Looking Main or Sub Forum Url #####: ",mainForumUrlPage)
				c.Visit(mainForumUrlPage)
				c.Wait()
			}
			//if page > totalpages {
			//	c.Visit("")
			//	c.Wait()
			//}
			defer wg2.Done()
		}
	}()
}

func databasePull(row int64 ,multiplier int64,workCellForPages int, wg1 *sync.WaitGroup)  {

	// yeni Object'e Scrap database'inin  forum_url(collection) ve articles(collection) tablosunu tanımladık.
	collUrls := DB.DBForumUrl()

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
		if result["last_page"] == nil || result["last_page"] == ""{
			continue
		}

		fmt.Println("array inside: ",result["last_page"].(int32))
		//a := fmt.Sprintf("%s",result["last_page"])
		forumPageSize := result["last_page"]
		page := forumPageSize.(int32)

		fmt.Println("##### Looking Main or Sub Forum Url #####:",mainForumUrl)
		fmt.Println("##### forum PageSize #####: ",page)

		jobs := int(math.Ceil(float64(page)/float64(workCellForPages))) //int to float64 after Ceil, float64 to int
		var wg2 sync.WaitGroup
		var l int
		for l=1; l <= workCellForPages; l++ {
			wg2.Add(jobs)
			go scrapThreads(l,jobs,workCellForPages,mainForumUrl,int(page),&wg2)
		}
		time.Sleep(time.Millisecond)
		wg2.Wait()
	}

	defer wg1.Done()
}

func PullThreadsMultiTask() {
	var i int32
	var wg1 sync.WaitGroup

	//collection := DB.DBForumUrl()
	//collSelector := DB.DBSiteInfo()
	SiteInfo := DB.DbSiteInfoMap()
	results := DB.DbForumUrlMap()

	// 126 main Forum Link --> goRoutine
	dbTotalGoRoutine := SiteInfo[0]["number_of_go_routines"].(int32)

	totalLinkLen := len(results)

	goRoutineForLinks := int32(math.Ceil(float64(dbTotalGoRoutine) / float64(2)))
	goRoutineForPages := int(math.Ceil(float64(dbTotalGoRoutine) / float64(2)))

	multiplier := int64(math.Ceil(float64(totalLinkLen) /float64(dbTotalGoRoutine)))



	for i=0; i<goRoutineForLinks ; i++ {
		wg1.Add(1)
		go databasePull(int64(i),multiplier,goRoutineForPages,&wg1)
	}

	wg1.Wait()
	fmt.Println("#######################\n PullThreadsMultiTask  \n####### THE END #######")
}