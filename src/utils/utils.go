package utils

import "github.com/ethereum/go-ethereum/ethclient"

var client = ethclient.Client{}
var contractAddress = ""

func ReturnClient()(client1 ethclient.Client)  {
	client1 = client
	return
}

func SetClient(client1 ethclient.Client){
	client = client1
}

func SetContractAddress(_address string){
	contractAddress = _address
}

func GetContractAddress()(contractAddress string){
	return contractAddress
}