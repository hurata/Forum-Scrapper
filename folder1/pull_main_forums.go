package folder1

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	DB "github.com/x/y/constants"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
)

func PullMainForums() {

	c := colly.NewCollector(
//		colly.MaxDepth(1),
		colly.Async(true),
//		colly.CacheDir("./forum_url_cache"),

	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 16})

	// Rotate proxies

	//rp, err := proxy.RoundRobinProxySwitcher("http://5.79.73.131:13040")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c.SetProxyFunc(rp)


	// yeni Object'e Scrap database'inin  forum_url(collection) tablosunu tanımladık.
	collection := DB.DBForumUrl()
//	collectionSelector := DB.DBSiteInfo()


	SiteInfo := DB.DbSiteInfoMap()

	// selecting selectors in here
	dbSiteLink := SiteInfo[0]["site"].(string)
	dbMainDiv := SiteInfo[0]["site_main_div"].(string)
	dbForumLink := SiteInfo[0]["forum_link"].(string)
	dbForumLinkAttr := SiteInfo[0]["forum_link_attr"].(string)
	dbForumName := SiteInfo[0]["forum_name"].(string)
	dbForumDesc := SiteInfo[0]["forum_desc"].(string)


	// Scrap HTML Settings
	c.OnHTML(dbMainDiv, func(e *colly.HTMLElement) {

		link :=  e.ChildAttr(dbForumLink,dbForumLinkAttr)
		forumName := e.ChildText(dbForumName)
		forumDesc := e.ChildText(dbForumDesc)


		// link, ?s 'den arındırılır
		base := dbSiteLink
		base2 := strings.Split(base, "://")
		base3 := base2[1]  			//   www.example.com
		//fmt.Println("----------------->",base3)

		i := strings.Index(link, base3)		// kaçıncı karakterden sonra "belirtilen" karakter geliyor?
		if i > -1 {
			cleanLinkForum := strings.Split(link, base3)
			linkTwo := base + cleanLinkForum[1]

			cleanLink := strings.Split(linkTwo, "?")
			cleanLink2 := cleanLink[0]

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
				fmt.Println("Link located on DataBase==>", cleanLink2)
			}

			if urlFiltered == nil {

				// Adding datas to database
				insertResult, err := collection.InsertOne(context.TODO(), bson.D{
					{Key: "url", Value: cleanLink2},
					{Key: "forum_name", Value: forumName},
					{Key: "forum_desc", Value: forumDesc},
					{Key: "sub_forum", Value: ""},
					{Key: "main_id", Value: ""},
				})
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Added Forum ID ==>",insertResult.InsertedID)
				// e.Request.Visit(link)
				// Adding datas to database End
			}
		} else {
			fmt.Println("Index not supported")
			fmt.Println(link)
		}
	})

	c.Visit(dbSiteLink)
	fmt.Println("Visiting HomePage ==> ",dbSiteLink)
	c.Wait()
	fmt.Println("#######################\n### PullMainForums ####\n####### THE END #######")
}