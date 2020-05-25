package main

import (
	"log"

	"./BackEnd"
	//"encoding/json"
    
    "github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
)



func main() {

	router := fasthttprouter.New()
        
    router.GET("/api/checkServer/:serverUrl", BackEnd.CheckServer)
    router.GET("/api/checkServerHistory", BackEnd.CheckServerHistory)

	router.OPTIONS("/api/checkServer/:serverUrl", BackEnd.CheckServerOptions)
    router.OPTIONS("/api/checkServerHistory", BackEnd.CheckServerHistoryOptions)
	
    router.NotFound = fasthttp.FSHandler("./FrontEnd", 0)
        
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))	
}
