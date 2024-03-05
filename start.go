package main

import L "github.com/x/y/folder1"

func main() {

//########################################
//########################################

//	Scraping start
//	Scraping https://www.ownedcore.com/forums/

//########################################
//########################################

//	The Proxy not enabled

//	All scrapping is making with 16 Thread.
//	So, "16" link request at the same time.

//########################################
//########################################

	//SiteInfo := DB.DbSiteInfoMap()
	//subForumStatus := SiteInfo[0]["sub_forum_status"]
	//
	//L.PullMainForums()  // in https://www.ownedcore.com/forums/  scraping main forums
	//
	//if subForumStatus == "1" {
	//	L.PullSubForums()	// in main Forums, scraping Sub forums
	//}
	//
	//L.PullTotalPage()	// in All Forums, pull total page amount

	L.PullThreadsMultiTask()	// in All Forums, scraping all Threads

	//L.Update()

}