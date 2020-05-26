package main

import (
	"fmt"

	"os/exec"
	"regexp"
)

func main() {
	fmt.Println("Testing system call...")
	var cmd = exec.Command("whois", "54.239.132.139")
	
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var salida string = string(output) 
	fmt.Printf("%v\n", salida)


	// r := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
    // fmt.Printf("%#v\n", r.FindStringSubmatch(`2015-05-27`)[1])
    // fmt.Printf("%#v\n", r.SubexpNames())

	r := regexp.MustCompile(`OrgName:\s*([\w\W]+?)\n[\w\W]*Country:\s*([\w\W]+?)\n`)
    fmt.Printf("%#v\n", r.FindStringSubmatch(salida)[1])
	fmt.Printf("%#v\n", r.FindStringSubmatch(salida)[2])

	//Status
	//curl -sL -w "%{http_code}" -I "www.google.com" -o /dev/null
	var cmdCurl = exec.Command(
		"curl",
		"-sL",
		`-w "%{http_code}"`,
		"-I",
		"www.google.com",
		"-o",
		"/dev/null")
	
	outputCurl, errCurl := cmdCurl.Output()
	if errCurl != nil {
		panic(errCurl)
	}

	var salidaCurl string = string(outputCurl) 
	fmt.Printf("%vCURLSTATUS\n", salidaCurl)
	
	
	fmt.Println("Finish")
}



