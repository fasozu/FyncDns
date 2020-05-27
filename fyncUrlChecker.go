package main

import (
	"log"
	"./BackEnd"	    
    "github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
)

func main() {
	config, err := BackEnd.GetConfiguration()
	if err != nil{
		log.Fatal(err.Error())
	}
	router := fasthttprouter.New()
        
    router.GET("/api/checkServer/:serverUrl", BackEnd.CheckServer)
    router.GET("/api/checkServerHistory", BackEnd.CheckServerHistory)

	router.OPTIONS("/api/checkServer/:serverUrl", BackEnd.CheckServerOptions)
    router.OPTIONS("/api/checkServerHistory", BackEnd.CheckServerHistoryOptions)
	
    router.NotFound = fasthttp.FSHandler("./FrontEnd", 0)
        
    log.Fatal(fasthttp.ListenAndServe(config.ListenServer, router.Handler))	
}
