package main

import (
	"fmt"

	"os/exec"
)

func main() {
	fmt.Println("Testing system call...")
	var cmd = exec.Command("whois", "54.239.132.139")
	
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(output))
	
	fmt.Println("Finish")
}



