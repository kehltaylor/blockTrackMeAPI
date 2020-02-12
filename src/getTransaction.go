package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type UserConnected struct {
	FirstName 		string
	LastName 		string
	PK 				string
	Count 			int
}

type ResponseAPI struct{
	Status 			string `json:"status"`
	Message 		string `json:"message"`
	Result 			[]TransactionResult `json:"result"`
}

type TransactionResult struct{
	BlockNumber 			string `json:"blockNumber"`
	TimeStamp 				string `json:"timeStamp"`
	Hash 					string `json:"hash"`
	Nonce 					string `json:"nonce"`
	BlockHash 				string `json:"blockHash"`
	TransactionIndex 		string `json:"transactionIndex"`
	From 					string `json:"from"`
	To 						string `json:"to"`
	Value 					string `json:"value"`
	Gas 					string `json:"gas"`
	GasPrice 				string `json:"gasPrice"`
	IsError 				string `json:"isError"`
	Txreceipt_status 		string `json:"txreceipt_status"`
	Input 					string `json:"input"`
	ContractAddress 		string `json:"contractAddress"`
	CumulativeGasUsed 		string `json:"cumulativeGasUsed"`
	GasUsed 				string `json:"gasUsed"`
	Confirmations 			string `json:"confirmations"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(data))

	return json.NewDecoder(r.Body).Decode(target)
}

func Find(a []UserConnected, x string) int {
	for _, n := range a {
		if x == n.PK {
			return n.Count
		}
	}
	return 0
}

func GetTransaction(w http.ResponseWriter, r *http.Request){

	var user User
	//var responseAPI  = ResponseAPI{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(pubKey.publicKey)
	var PK = user.PubKey
	DB.Where("pub_key = ?",user.PubKey).First(&user)

	/*res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}*/

	//getJson("http://api-ropsten.etherscan.io/api?module=account&action=txlist&startblock=0&endblock=99999999&sort=asc&address="+ user.PubKey + "&apikey=HMRPT34GHBB29RDT4YW44AIABZRH4FPPKK", &responseAPI)
	response, err := http.Get("http://api-ropsten.etherscan.io/api?module=account&action=txlist&startblock=0&endblock=99999999&sort=asc&address="+ user.PubKey + "&apikey=HMRPT34GHBB29RDT4YW44AIABZRH4FPPKK")

	//res, err := http.Get("https://www.citibikenyc.com/stations/json")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}
	myList := []UserConnected{}
	s, err := getTxs([]byte(body))
	fmt.Println(s.Result)
	fmt.Println(len(s.Result))

	for _,element := range s.Result{
		if element.From == PK {
			DB.Where("pub_key = ?",element.From).First(&user)
			u := UserConnected{}
			u.PK = element.To
			u.LastName = user.LastName
			u.FirstName = user.FistName
			u.Count = Find(myList, u.PK) + 1
			myList = append(myList, u)
		}else{
			DB.Where("pub_key = ?",element.From).First(&user)
			u := UserConnected{}
			u.PK = element.From
			u.LastName = user.LastName
			u.FirstName = user.FistName
			u.Count = Find(myList, u.PK) + 1
			myList = append(myList, u)
		}
	}

	res, err := json.Marshal(myList)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	/*defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&responseAPI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("{}", responseAPI)*/


	/*if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		json.NewDecoder(r.Body).Decode(&responseAPI)
		fmt.Println("{}", responseAPI.result)
	}*/

	//fmt.Println("{}",res)

}
func getTxs(body []byte) (*ResponseAPI, error) {
	var s = new(ResponseAPI)
	err := json.Unmarshal(body, &s)
	if(err != nil){
		fmt.Println("whoops:", err)
	}
	return s, err
}