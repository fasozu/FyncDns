package BackEnd

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
    ListenServer string `json:"listenServer"` //ip and port to listen server default :8080
	CacheDurationSeconds int `json:cachedDurationSeconds"` //Seconds to return the cached version instead of make a new request default 1 hour = (60*60)seconds
	ApiRootSslLabs string `json:"apiRootSslLabs"` // Root of external api to check grade default https://api.ssllabs.com/api/v3/analyze?host=
	UrlTimeoutSeconds int `json:"urlTimeoutSeconds"` // Number of seconds to wait for url's response default 6 seconds
	ServerDownHttpCodes []string `json:"serverDownHttpCodes"` // List of url http codes for determining site is down default: [000, 500]
	DBCockroachConnection string  `json:"dBCockroachConnection"` // default "postgresql://maxroach@localhost:26257/fucdb?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt"
}

var config Configuration
var isConfigLoaded bool = false

//Get the Configuration struct, if file not exist, it creates one with
//default values.
func GetConfiguration()(Configuration, error){
	if isConfigLoaded {
		return config,nil
	}
	
	configTemp, err := loadConfiguaration()
	if err != nil{
		return config,err 
	}
	
	config = configTemp
	isConfigLoaded = true
	return config,nil
}


// Try to load config.json, if it not exist, create one with default
// values
func loadConfiguaration() (Configuration, error){
	if _, err := os.Stat("config.json"); err == nil {
		b, errRead := ioutil.ReadFile("config.json") 
		if errRead != nil {
			return config,err
		}
		json.Unmarshal(b, &config)
		return config,nil
	} else {
		errorCodes := []string{"500","000"}
		configTemp := Configuration{
			ListenServer: ":8080",
			CacheDurationSeconds: (60*60), //1 hour
			ApiRootSslLabs: "https://api.ssllabs.com/api/v3/analyze?host=",
			UrlTimeoutSeconds: 6,
			ServerDownHttpCodes: errorCodes,
			DBCockroachConnection: "postgresql://maxroach@localhost:26257/fucdb?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt",
		}
		err := saveConfiguration(configTemp)
		config = configTemp
		
		if err != nil {			
			return config,err
		}		
		return config,nil
	}
}


//Save `config` to config.json
func saveConfiguration(configInput Configuration) (error){
	b, err := json.MarshalIndent(configInput,"","  ")		
	if err != nil {			
		return err
	}

	f, errFile := os.Create("config.json")
    if errFile != nil {        
        return errFile
    }
	
    _, errWrite := f.WriteString(string(b))
    if errWrite != nil {
		f.Close()
		return errWrite
    }
	    
    errClose := f.Close()
    if errClose != nil {        
        return errClose
    }

	return nil
}
