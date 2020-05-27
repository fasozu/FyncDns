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
	//"./DBHelper"
	"crypto/md5"
    "encoding/hex"
	"reflect"
)

func CheckServer(ctx *fasthttp.RequestCtx) {
		
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")

	var serverUrl string = fmt.Sprint(ctx.UserValue("serverUrl"))
	var remoteAddress string = fmt.Sprint(ctx.RemoteIP().String())
	
	var hashIdentifier = GetMD5Hash(remoteAddress)
	fmt.Println("Direccion remota: ",remoteAddress)
	fmt.Println("Hash identifier: ",hashIdentifier)
	var response OutputResponse
	var responseCached OutputResponse
	
	db, errDb := GetDb()
	if errDb != nil{
		response = OutputResponse{
			Success: false,			
			ErrorMessage: errDb.Error(),
		}		
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))        				
		return
	}
	defer db.Close()

	errAddCheckHistory := AddCheckHistory(db, hashIdentifier, serverUrl)
	if errAddCheckHistory != nil{
		response = OutputResponse{
			Success: false,			
			ErrorMessage: errAddCheckHistory.Error(),
		}		
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))        				
		return
	}

	fmt.Println("After checksave")

	lastTimeChecked, jsonResponseCached, currentTime,_ := GetUrlCache(db, serverUrl)
	if jsonResponseCached != ""{
		json.Unmarshal([]byte(jsonResponseCached), &responseCached)
	}
	fmt.Println("lastTime")
	fmt.Println(lastTimeChecked)
	fmt.Println("currentTime")
	fmt.Println(currentTime)
	
	if (currentTime - 5) > lastTimeChecked || lastTimeChecked == 0 {//Cached invalid
		//Get new info for url
		responseRaw, err := getDataRawApi(serverUrl)
		
		if err != nil{
			fmt.Println("Cannot get data raw")
			response = OutputResponse{
				Success: false,			
				ErrorMessage: err.Error(),
			}		
			b, err := json.Marshal(response)
			if err != nil {
				fmt.Fprintf(ctx, "Error %s \n", err.Error())
				return
			}
			fmt.Fprint(ctx, string(b))        		
			return
		}
		

		
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

		//Status code:
		statusCode, _  := getUrlStatusCode(serverUrl)
		response.StatusCode = statusCode
		if statusCode == "200"{
			response.Is_down = false
		}else{
			response.Is_down = true
		}

		title, icon, _ := getTitleAndIcon(serverUrl)
		response.Title = title
		response.Logo = icon
		response.Url = serverUrl

		response.Previous_ssl_grade = responseCached.Ssl_grade

		if reflect.DeepEqual(response.Servers, responseCached.Servers){
			response.Servers_changed = false
		}else{
			response.Servers_changed = true
		}

		//---------------------
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		AddUrlCache(db, serverUrl, string(b))
		fmt.Fprint(ctx, string(b))            
		
	}else{
		//use the cached
		
		fmt.Println("Using cached")
		responseCached.Is_cached = true
		b, err := json.Marshal(responseCached)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))            
	}
	
}


func CheckServerHistory(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	var response OutputResponseHistory
	db, errDb := GetDb()
	if errDb != nil{
		response.Success = false
		response.ErrorMessage = errDb.Error()
		
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))        				
		return
	}
	defer db.Close()

	var remoteAddress string = fmt.Sprint(ctx.RemoteIP().String())
	var hashIdentifier = GetMD5Hash(remoteAddress)
	fmt.Println("Direccion remota: ",remoteAddress)
	fmt.Println("Hash identifier: ",hashIdentifier)

	urls, err := GetCheckHistory(db, hashIdentifier)
	if err != nil {
		response.Success = false
		response.ErrorMessage = err.Error()
		b, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))        				
		return
	}

	for _, url := range urls {
		var urlInfo OutputResponse
		_, jsonResponseCached, _, _ := GetUrlCache(db, url)
		fmt.Println("------")
		fmt.Println(url)
		fmt.Println(jsonResponseCached)
		if jsonResponseCached != "" {
			json.Unmarshal([]byte(jsonResponseCached), &urlInfo)
			response.Items = append(response.Items, urlInfo)
		}
	}
	
	response.Success = true
	
	b, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(ctx, "Error %s \n", err.Error())
		return
	}
	fmt.Fprint(ctx, string(b))		
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
		return "","",err
	}

	var salida string = string(output) 	

	r := regexp.MustCompile(`OrgName:\s*([\w\W]+?)\n[\w\W]*Country:\s*([\w\W]+?)\n`)
	var orgName string = r.FindStringSubmatch(salida)[1]
	var country string =  r.FindStringSubmatch(salida)[2]

	return orgName,country,nil
}

func getUrlStatusCode(url string) (string,error){
	//curl -sL -w "%{http_code}" -I "www.google.com" -o /dev/null
	var cmdCurl = exec.Command(
		"curl",
		"-sL",
		`-w%{http_code}`,
		"-I",
		url,
		"-o",
		"/dev/null")
	
	outputCurl, errCurl := cmdCurl.Output()
	if errCurl != nil {
		return "", errCurl
	}

	var salidaCurl string = string(outputCurl) 
	//fmt.Printf("%vCURLSTATUS\n", salidaCurl)
	return salidaCurl,nil
}

func getTitleAndIcon(url string) (string,string,error){
	

	var title = ""
	var icon = ""
	
	var cmd = exec.Command("curl", "-L", url)
	output, err := cmd.Output()
	if err != nil {
		return title,icon,err
	}
	
	var htmlPage string = string(output) 
	
	rTitle := regexp.MustCompile(`<title>([\w\W]*?)</title>`)
	if len(rTitle.FindStringSubmatch(htmlPage))>1{
		title = rTitle.FindStringSubmatch(htmlPage)[1]
	}
	

	rIconTag := regexp.MustCompile(`(<link[^>]*?rel\s*="[^"]*?icon[^"]*?".*?>)`)
	if len(rIconTag.FindStringSubmatch(htmlPage)) > 1{
		iconoTag := rIconTag.FindStringSubmatch(htmlPage)[1]		
		rIcon := regexp.MustCompile(`href\s*=\s*"([^"]*?)"`)
		if len(rIcon.FindStringSubmatch(iconoTag))> 1{
			icon = rIcon.FindStringSubmatch(iconoTag)[1]
		}		
	}
	

	
	
	return title,icon,nil
}


func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

