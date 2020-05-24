package main

import (
	"fmt"

	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Testing Http Request...")


	resp, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))

	fmt.Println("Finish")
}
