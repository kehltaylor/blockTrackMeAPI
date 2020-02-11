package main

import (
	_ "./contracts"
	"./pkg"
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"time"
)

var DB *gorm.DB
var err error

type ClientConnection struct {
	Status    string
}

type User struct {
	gorm.Model
	FistName  	string `json:"firstName"`
	LastName 	string `json:"lastName"`
	PubKey		string `json:"pubKey"`
}

func main() {
	_, err := dbConnect()
	if err != nil {
		log.Println("connection to DB failed, aborting...")
		log.Fatal(err)
	}

	log.Println("connected to DB")
	r := mux.NewRouter()
	r.HandleFunc("/", connection).Methods("POST")
	r.HandleFunc("/showUserHandler", showUserHandler).Methods("POST")
	r.HandleFunc("/createUserHeader", createUserHandler).Methods("POST")
	r.HandleFunc("/deployContract", features.DeployContract).Methods("POST")
	r.HandleFunc("/transaction", features.GenerateTransaction).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func dbConnect()(*gorm.DB, error){

	connectionParams := "dbname=goland user=postgres password=root sslmode=disable host=postgres"
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", connectionParams) // gorm checks Ping on Open
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return DB, err
	}

	DB.AutoMigrate(&User{})

	testPost := User{FistName: "Kojak", LastName: "Zoda", PubKey: "SexyBoy77"}
	DB.Create(&testPost)

	return DB, err
}

func connection(w http.ResponseWriter, r *http.Request) {
	client, err := ethclient.Dial("https://ropsten.infura.io/v3/511162a74a0c4a80a9fbab7b9d2718b8")
	if err != nil {
		log.Fatal(err)
	}
	_ = client // we'll use this in the upcoming sections
	clientConnection := ClientConnection{"Done"}

	js, err := json.Marshal(clientConnection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func handler(w http.ResponseWriter, r *http.Request) {
	title := "Test"

	from := ""
	if r.URL != nil {
		from = r.URL.String()
	}
	if from != "/favicon.ico" {
		log.Printf("title: %s\n", title)
	}

	client, err := ethclient.Dial("https://ropsten.infura.io/v3/511162a74a0c4a80a9fbab7b9d2718b8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	utils.SetClient(*client)
	_ = client // we'll use this in the upcoming sections
}






