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
	"strconv"
	"strings"
	"sync"
)

type outsideVariable struct {
	pageNumber string

}

func totalPage(row int64 ,multiplier int64,wg1 *sync.WaitGroup){
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

	collection := DB.DBForumUrl()
	//collSelector := DB.DBSiteInfo()

	// ################################ vv  collection
	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.

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

	//results := DB.DbForumUrlMap()
	SiteInfo := DB.DbSiteInfoMap()

	// selecting selectors in here
	//	dbSiteLink := SiteInfo[0]["site"].(string)
	dbLastPageDiv := SiteInfo[0]["forum_lastpage_div"].(string)
	dbLastPage := SiteInfo[0]["forum_lastpage"].(string)
	dbLastPageAttr := SiteInfo[0]["forum_lastpage_attr"].(string)
	dbLastPageNumber := SiteInfo[0]["forum_lastpage_number"].(string)


	// Last-Page-Button link
	c.OnHTML(dbLastPageDiv, func(h *colly.HTMLElement) {
		prevNext := h.ChildAttr(dbLastPage, dbLastPageAttr)
		prevNextNumber := h.ChildAttr(dbLastPageNumber, dbLastPageAttr)
//		prevNextNumber2 := h.ChildAttr("a.infiniteNext", dbLastPageAttr)
		//-->the last page-->  prevNext := h.ChildAttr("span.first_last a", "href")
		urlMain := h.Request.URL.String()
//		cleanLinkForum := strings.Split(prevNextNumber, "index")

		if prevNext =="" {

			if prevNextNumber == ""{

				lastPage := 1
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

			}else {

				TakeLastPageWithNextButton(urlMain)

			}
		}else {

			cleanLinkForum := strings.Split(prevNext, "index")
			lastPageString := cleanLinkForum[1]
			cleanLinkLast := strings.Split(lastPageString, ".html")
			lastPage := cleanLinkLast[0]

			lastPage2 ,err := strconv.Atoi(lastPage)
			if err != nil {
				log.Fatal(err)
			}

			filterUpdate := bson.M{"url": bson.M{"$eq": urlMain}}

			update := bson.M{
				"$set": bson.M{
					"last_page": lastPage2,
				},
			}

			collection.UpdateOne(context.TODO(),filterUpdate, update)
			fmt.Println("----- Added Sub/Forum LastPage ==>",lastPage2)
			//e.Request.Visit(link)
			// Adding datas to database End


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

func takePageNumber(startLink string) string {
	d := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	d.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	collection := DB.DBForumUrl()

	// ################################ vv  collection
	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.

	cursor, err := collection.Find(context.TODO(), bson.D{{}})
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
	dbLastPageDiv := SiteInfo[0]["forum_lastpage_div"].(string)
	dbLastPageAttr := SiteInfo[0]["forum_lastpage_attr"].(string)
	checkVarable := outsideVariable{}

	d.OnHTML(dbLastPageDiv, func(l *colly.HTMLElement) {

		prevNextNumber2 := l.ChildAttr("a.infiniteNext", dbLastPageAttr)
		checkVarable.pageNumber = prevNextNumber2

	})
	d.Visit(startLink)
	d.Wait()
	return checkVarable.pageNumber
}

func TakeLastPageWithNextButton(testLink string){
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})


	collection := DB.DBForumUrl()


	// ################################ vv  collection
	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.

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
	// ################################ ^^

	SiteInfo := DB.DbSiteInfoMap()

	// selecting selectors in here
	dbLastPageDiv := SiteInfo[0]["forum_lastpage_div"].(string)

	// Last-Page-Button Scrapping
	c.OnHTML(dbLastPageDiv, func(h *colly.HTMLElement) {

		urlMain := h.Request.URL.String()
		checkURL := urlMain
		for i:=0;i<100;i++{

			ss := takePageNumber(checkURL)
			checkURL =ss
			if ss ==""{

			//	cleanLinkForum2 := strings.Split(checkURL, "index")

				lastPage := i+1
				filterUpdate := bson.M{"url": bson.M{"$eq": urlMain}}

				update := bson.M{
					"$set": bson.M{
						"last_page": lastPage,
					},
				}

				collection.UpdateOne(context.TODO(),filterUpdate, update)
				fmt.Println("----- Added Sub/Forum LastPage ==>",i+1)

				// Adding datas to database End

				i=100
			}
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

	// dataBase'den alınıp MainForum linkine gidilecek
	fmt.Println("##### Looking Main Forum Url #####:",testLink)
	c.Visit(testLink)
	c.Wait()

	/*  // if you wanna take the last pages, you can use that codes
	counter := 0
	for _, result := range results {
		counter++
		mainForumUrl := result["url"].(string)
		if result["last_page"] == nil{
			fmt.Println("##### Looking Main Forum Url #####",counter,":",mainForumUrl)
			c.Visit(mainForumUrl)
			c.Wait()
		}
	}
	*/

}

func PullTotalPage() {
	var i int32
	var wg1 sync.WaitGroup
	// 315 Sub and Main Links

	//collection := DB.DBForumUrl()
	//collSelector := DB.DBSiteInfo()
	SiteInfo := DB.DbSiteInfoMap()
	results := DB.DbForumUrlMap()

	// 126 main Forum Link --> goRoutine
	dbTotalGoRoutine := SiteInfo[0]["number_of_go_routines"].(int32)

	totalLinkLen := len(results)
	multiplier := int64(math.Ceil(float64(totalLinkLen) /float64(dbTotalGoRoutine)))

	fmt.Println("multiplier:",multiplier)

	for i=0;i<dbTotalGoRoutine;i++ {
		wg1.Add(1)
		go totalPage(int64(i),multiplier, &wg1)
	}

	wg1.Wait()
	fmt.Println("#######################\n#### PullTotalPage ####\n####### THE END #######")
}