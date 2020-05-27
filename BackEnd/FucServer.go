package BackEnd

import (
    "fmt"
	"io/ioutil"
	"net/http"
    "github.com/valyala/fasthttp"	
	"encoding/json"
	"regexp"
	"os/exec"
	"crypto/md5"
    "encoding/hex"
	"reflect"
	"strings"
	"strconv"
)

//Endpoint 1 - check url server
func CheckServer(ctx *fasthttp.RequestCtx) {
	
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")

	var response OutputResponse
	var responseCached OutputResponse

	config, err := GetConfiguration()
	
	if err != nil{
		response = OutputResponse{
			Success: false,			
			ErrorMessage: err.Error(),
		}		
		b, errJson := json.Marshal(response)
		if errJson != nil {
			fmt.Fprintf(ctx, "Error %s \n", errJson.Error())
			return
		}
		fmt.Fprint(ctx, string(b))
		return
	}
	
	var serverUrl string = fmt.Sprint(ctx.UserValue("serverUrl"))
	var remoteAddress string = fmt.Sprint(ctx.RemoteIP().String())	
	var hashIdentifier = GetMD5Hash(remoteAddress)


	//Opening DB--------------------------------------------
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
	//------------------------------------------------------

	//Save the url in client's  history --------------------
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
	//------------------------------------------------------

	//Load Cache--------------------------------------------
	lastTimeChecked, jsonResponseCached, currentTime,_ := GetUrlCache(db, serverUrl)
	if jsonResponseCached != ""{
		json.Unmarshal([]byte(jsonResponseCached), &responseCached)
	}
	//------------------------------------------------------
	
	if (currentTime - int64(config.CacheDurationSeconds)) > lastTimeChecked || lastTimeChecked == 0 {//Cached invalid		
		responseRaw, err := getDataRawApi(serverUrl)		
		if err != nil{			
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
		
		//Check grade ------------------------------------------
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

		//Status server:
		response.Is_down = getServerIsDown(serverUrl)
		
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
		//use cached
		
		responseCached.Is_cached = true
		b, err := json.Marshal(responseCached)
		if err != nil {
			fmt.Fprintf(ctx, "Error %s \n", err.Error())
			return
		}
		fmt.Fprint(ctx, string(b))            
	}
	
}

//End point 2 - Check history
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

	// For each url in check history return 
	for _, url := range urls {
		var urlInfo OutputResponse
		_, jsonResponseCached, _, _ := GetUrlCache(db, url)
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

//Return the response from external api
func getDataRawApi(url string) (*InputResponse,error){
	config, _ := GetConfiguration()
	
	resp, err := http.Get(fmt.Sprintf("%s%s", config.ApiRootSslLabs, url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response InputResponse 

	json.Unmarshal(body, &response)
	
	return &response,nil
}


// Return the Owner and Country of `ip`
func getOwnerAndCountry(ip string) (string,string,error){
	var orgName string 
	var country string

	
	var cmd = exec.Command("whois", ip)
	
	output, _ := cmd.Output()
	
	var outputWhois string = string(output) 	

	r := regexp.MustCompile(`OrgName:\s*([\w\W]+?)\n[\w\W]*Country:\s*([\w\W]+?)\n`)
	infoMatched := r.FindStringSubmatch(outputWhois)
	if len(infoMatched) > 2{
		orgName = infoMatched[1]
		country = infoMatched[2]
	}
	return orgName,country,nil
}

// Return http code from url, 000 if timeout
func getUrlStatusCode(url string) (string,error){
	config, _ := GetConfiguration()
	//curl  --max-time 6 -sL -w "%{http_code}" -I "www.google.com" -o /dev/null
	var cmdCurl = exec.Command(
		"curl",
		"--max-time",
		strconv.Itoa(config.UrlTimeoutSeconds),
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

	var responseCurl string = string(outputCurl) 
	return responseCurl,nil
}

//Return title and icon from url's html 
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
			if( !strings.HasPrefix(icon, "http://") && !strings.HasPrefix(icon, "https://") ){
				icon = "http://"+url+"/"+icon
			}
		}		
	}	
	return title,icon,nil
}

// Return true if server is down
func getServerIsDown(serverUrl string) (bool){
	config, _ := GetConfiguration()
	statusCode, _  := getUrlStatusCode(serverUrl)
	for _, downHttpCode := range config.ServerDownHttpCodes {
        if statusCode == downHttpCode {
            return true
        }
    }	
	return false
} 

//Return the md5sum for input string
func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

