# FyncUrlChecker

## Dependencies

- Cockcroach installed and configured with a database and user created. 
- Curl
- whois
- Go packages installed:
  - github.com/buaazp/fasthttprouter
  - github.com/lib/pq

## Configuration

The file is `config.json` in the same directory of `fyncUrlChecker.go`, if the file not exists, a file will be created with default values:

*listenServer:* ip and port to listen server default :8080

*CacheDurationSeconds:* Seconds to return the cached version instead of make a new request default 1 hour = (60*60)seconds

*ApiRootSslLabs:* Root of external api to check grade default https://api.ssllabs.com/api/v3/analyze?host=

*UrlTimeoutSeconds:* Number of seconds to wait for url's response default 6 seconds

*ServerDownHttpCodes:* List of url http codes for determining site is down default: [000, 500]

*DBCockroachConnection:* String to connect to database, default "postgresql://maxroach@localhost:26257/fucdb?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt"
    
```{js}
{
  "listenServer": ":8080",
  "CacheDurationSeconds": 6,
  "apiRootSslLabs": "https://api.ssllabs.com/api/v3/analyze?host=",
  "urlTimeoutSeconds": 6,
  "serverDownHttpCodes": [
    "500",
    "000"
  ],
  "dBCockroachConnection": "postgresql://maxroach@localhost:26257/fucdb?ssl=true\u0026sslmode=require\u0026sslrootcert=certs/ca.crt\u0026sslkey=certs/client.maxroach.key\u0026sslcert=certs/client.maxroach.crt"
}
```


## Running

`$ go run fyncUrlChecker.go`

## Notes

### servers_changed

The `servers_changed` interpretation consists in use a cached version for 1 hour of the time specified before try again with SslLabs's webservice. 

### SslLabs's webservice 

The response of the url takes time and first returns the progress of the analysis and not the complete analysis result. An improved version for this implementation is take this behavior in account and paint the progress.

An unexpected behavior is the data if the analysis is ephemeral and the full analysis result is available for short time.

### Interpretation of site down

The strategy applied was to use the http response code to determine the site's status. By default if no answer or `500` then the site is marked down. This behavior can be changed inside `config.json` to add more response codes of interest.

 

### Gui implementation

For the gui vue and bootstrap was used. This was developed using [Nuxt](https://nuxtjs.org/) framework. 

From Nuxt, a static site is generated and put on project's  `FrontEnd` directory. This way all files are served from the Go implementation. 



