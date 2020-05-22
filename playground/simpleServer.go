package main

import (
    "fmt"
    "log"
    "encoding/json"
    
    "github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
)


type Server struct {
    Address string `json:"address"`
    Ssl_grade string `json:"ssl_grade"`
    Country string `json:"country"`
    Owner string `json:"owner"`     
}

type Response struct {
        Servers []Server `json:"servers"`
        Servers_changed bool `json:"servers_changed"`
        Ssl_grade string `json:"ssl_grade"`
        Logo string `json:"logo"`
        Is_down bool `json:"is_down"`    
}

func Index(ctx *fasthttp.RequestCtx) {
    fmt.Fprint(ctx, "Welcome!\n")
}

func SimpleJson(ctx *fasthttp.RequestCtx) {
    
	var s1 = Server{
		Address: "Server1",
		Ssl_grade: "B",
		Country: "US",
		Owner: "Amazon.com, Inc.",
	}

    var s2 = Server{
        Address: "Server2",
        Ssl_grade: "A",
        Country: "CO",
        Owner: "Godaddy",               
    }

    var r1 = Response{
        Servers: []Server{},
        Servers_changed: false,
        Ssl_grade: "B",
        Logo:  "https://server.com/icon.png",
        Is_down: false,
    }

    r1.Servers = append(r1.Servers, s1)
    r1.Servers = append(r1.Servers, s2)

    b, err := json.Marshal(r1)
    if err != nil {
        fmt.Fprintf(ctx, "Error %s \n", err.Error())        
    }else{
        fmt.Fprint(ctx, string(b))        
    }
}


func Hello(ctx *fasthttp.RequestCtx) {
    fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

func main() {

    router := fasthttprouter.New()
        
    router.GET("/api", Index)
    router.GET("/api/hello/:name", Hello)
    router.GET("/api/simplejson/:name", SimpleJson)

    router.NotFound = fasthttp.FSHandler("./static", 0)
        
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
