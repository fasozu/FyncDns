package BackEnd

import (
    "fmt"
	//"github.com/buaazp/fasthttprouter"
	"io/ioutil"
	"net/http"
    "github.com/valyala/fasthttp"	
	"encoding/json"
	"regexp"
	"os/exec"
)

func CheckServer(ctx *fasthttp.RequestCtx) {
	var serverUrl string = fmt.Sprint(ctx.UserValue("serverUrl"))	
	responseRaw, err := getDataRawApi(serverUrl)
			
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")

	var response OutputResponse
	if err != nil{
		response = OutputResponse{
			Success: false,			
			ErrorMessage: err.Error(),
		}
		fmt.Fprint(ctx, response)
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))        		
		return
	}
	
	
	fmt.Println("Received...")
	fmt.Println(responseRaw.Host)

	
	response = OutputResponse{
		Success: true,
	}
	
	//---------------------	
	var worstGrade string = ""
	for _, endpoint := range responseRaw.Endpoints {
		if worstGrade < endpoint.Grade {
			worstGrade = endpoint.Grade
		}
		owner, country, _ := getOwnerAndCountry(endpoint.IpAddress)
		var server = OutputServer{
			Address: endpoint.IpAddress,
			Ssl_grade: endpoint.Grade,
			Country: country,
			Owner: owner,
		}
		response.Servers = append(response.Servers, server)
	}
	response.Ssl_grade = worstGrade
	//---------------------
	b, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(ctx, "Error %s \n", err.Error())
		return
	}
	fmt.Fprint(ctx, string(b))            
}


func CheckServerHistory(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")	
	fmt.Fprint(ctx, "[]")            
}

func CheckServerOptions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")	
	fmt.Fprint(ctx, "")            
}

func CheckServerHistoryOptions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")	
	fmt.Fprint(ctx, "")            
}

func getDataRawApi(url string) (*InputResponse,error){
	fmt.Println("Testing Http Request...")

	resp, err := http.Get(fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", url))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var response InputResponse 

	json.Unmarshal(body, &response)
	
	return &response,nil
}


func getOwnerAndCountry(ip string) (string,string,error){

	var cmd = exec.Command("whois", ip)
	
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var salida string = string(output) 
	fmt.Printf("%v\n", salida)

	r := regexp.MustCompile(`OrgName:\s*([\w\W]+?)\n[\w\W]*Country:\s*([\w\W]+?)\n`)
	var orgName string = r.FindStringSubmatch(salida)[1]
	var country string =  r.FindStringSubmatch(salida)[2]

	return orgName,country,nil

}

func getOwnerAndCountry(url string) (string,error){
	//curl -sL -w "%{http_code}" -I "www.google.com" -o /dev/null
	var cmdCurl = exec.Command(
		"curl",
		"-sL",
		`-w "%{http_code}"`,
		"-I",
		url,
		"-o",
		"/dev/null")
	
	outputCurl, errCurl := cmdCurl.Output()
	if errCurl != nil {
		panic(errCurl)
	}

	var salidaCurl string = string(outputCurl) 
	fmt.Printf("%vCURLSTATUS\n", salidaCurl)

}
