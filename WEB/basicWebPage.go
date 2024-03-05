package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Welcome struct {
	Name string
	Time string
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	//Bir "Welcome struct nesnesi" örneği oluşturun ve bazı rastgele bilgileri iletin.
	//Kullanıcının adını URL'den sorgu parametresi olarak alacağız.
	welcome := Welcome{"Server V", time.Now().Format(time.Stamp)}

	//Go'ya html dosyanızı tam olarak nerede bulabileceğimizi söyleriz.
	//Go'dan html dosyasını ayrıştırmasını istiyoruz (Göreli yola dikkat edin).
	//Herhangi bir hatayı işleyen ve ölümcül hatalar varsa durduran bir "template.Must ()" çağrısına sarıyoruz.

	templates := template.Must(template.ParseFiles("welcome.html"))

	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println("maps! :", r.Form) // print information on server side.
	fmt.Println("path :", r.URL.Path)
	fmt.Println("scheme :", r.URL.Scheme)
	fmt.Println("r.Forum[\"site\"] --> :",r.Form["site"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	// fmt.Fprintf(w, "Hello astaxie!") // write data to response


	//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
	if name := r.FormValue("name"); name != "" {
		welcome.Name = name
	}
	//If errors show an internal server error message
	//I also pass the welcome struct to the welcome-template.html file.
	if err := templates.ExecuteTemplate(w, "welcome.html", welcome); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("welcome.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("\nmaps of FORM : ", r.Form)

		goRoutineNumber, _ := strconv.Atoi(r.FormValue("go_routines"))


		fmt.Println("...\nSite Link:", r.FormValue("site"))
		fmt.Println("How Much GoRoutine:", goRoutineNumber)
		fmt.Println("Sub Forum Status:", r.FormValue("sub_forum_status"))

		fmt.Println("Forum Main Div:", r.FormValue("site_main_div"))
		fmt.Println("Forum Name div:", r.FormValue("forum_name"))
		fmt.Println("Forum Description div:", r.FormValue("forum_desc"))
		fmt.Println("Forum Link div:", r.FormValue("forum_link"))
		fmt.Println("Forum Link Attribute:", r.FormValue("forum_link_attr"))

		fmt.Println("Forum Lastpage Main div:", r.FormValue("forum_lastpage_div"))
		fmt.Println("Forum Lastpage div:", r.FormValue("forum_lastpage"))
		fmt.Println("Forum Lastpage Attribute:", r.FormValue("forum_lastpage_attr"))
		fmt.Println("Forum Next Page:", r.FormValue("forum_lastpage_number"))

		fmt.Println("Forum Threads Main div:", r.FormValue("thread_main_div"))
		fmt.Println("Forum Threads Name div:", r.FormValue("thread_name"))
		fmt.Println("Forum Threads Link div:", r.FormValue("thread_link"))
		fmt.Println("Forum Threads Link Attribute:", r.FormValue("thread_link_attr"))
		fmt.Println("")

		// Add database

		//  MongoDB Settings
		clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		collection := client.Database("scrap").Collection("site_info")

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

		// Site Infos Update
		xx := results[0]["site"]

		filterUpdate := bson.M{"site": bson.M{"$eq": xx}}
		update := bson.M{
			"$set": bson.D{
				{Key: "site", Value: r.FormValue("site")},
				{Key: "number_of_go_routines", Value: goRoutineNumber},
				{Key: "sub_forum_status", Value: r.FormValue("sub_forum_status")},

				{Key: "site_main_div", Value: r.FormValue("site_main_div")},
				{Key: "forum_name", Value: r.FormValue("forum_name")},
				{Key: "forum_desc", Value: r.FormValue("forum_desc")},
				{Key: "forum_link", Value: r.FormValue("forum_link")},
				{Key: "forum_link_attr", Value: r.FormValue("forum_link_attr")},

				{Key: "forum_lastpage_div", Value: r.FormValue("forum_lastpage_div")},
				{Key: "forum_lastpage", Value: r.FormValue("forum_lastpage")},
				{Key: "forum_lastpage_attr", Value: r.FormValue("forum_lastpage_attr")},
				{Key: "forum_lastpage_number", Value: r.FormValue("forum_lastpage_number")},

				{Key: "thread_main_div", Value: r.FormValue("thread_main_div")},
				{Key: "thread_name", Value: r.FormValue("thread_name")},
				{Key: "thread_link", Value: r.FormValue("thread_link")},
				{Key: "thread_link_attr", Value: r.FormValue("thread_link_attr")},
			},
		}
		collection.UpdateOne(context.TODO(),filterUpdate, update)


	// just add site infos
 	/*	collection.InsertOne(context.TODO(), bson.D{
			{Key: "site", Value: r.FormValue("site")},
			{Key: "number_of_go_routines", Value: goRoutineNumber},
			{Key: "sub_forum_status", Value: r.FormValue("sub_forum_status")},

			{Key: "site_main_div", Value: r.FormValue("site_main_div")},
			{Key: "forum_name", Value: r.FormValue("forum_name")},
			{Key: "forum_desc", Value: r.FormValue("forum_desc")},
			{Key: "forum_link", Value: r.FormValue("forum_link")},
			{Key: "forum_link_attr", Value: r.FormValue("forum_link_attr")},

			{Key: "forum_lastpage_div", Value: r.FormValue("forum_lastpage_div")},
			{Key: "forum_lastpage", Value: r.FormValue("forum_lastpage")},
			{Key: "forum_lastpage_attr", Value: r.FormValue("forum_lastpage_attr")},
			{Key: "forum_lastpage_number", Value: r.FormValue("forum_lastpage_number")},

			{Key: "thread_main_div", Value: r.FormValue("thread_main_div")},
			{Key: "thread_name", Value: r.FormValue("thread_name")},
			{Key: "thread_link", Value: r.FormValue("thread_link")},
			{Key: "thread_link_attr", Value: r.FormValue("thread_link_attr")},
		})*/

		fmt.Println("database updated \n")
	}
}

//Go application entrypoint
func main() {

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	http.Handle("/templates/", //final url can be anything
		http.StripPrefix("/templates/",
			http.FileServer(http.Dir("templates")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
	http.HandleFunc("/" , sayHello)
	http.HandleFunc("/login", login)

	//Start the web server, set the port to listen to 1453. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening...")
	fmt.Println(http.ListenAndServe(":1453", nil))
}