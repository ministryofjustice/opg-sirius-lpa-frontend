package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s%s", os.Getenv("PORT"), os.Getenv("HEALTHCHECK")))
	fmt.Println("Checking ", fmt.Sprintf("http://127.0.0.1:%s%s", os.Getenv("PORT"), os.Getenv("HEALTHCHECK")))
	if err != nil {
		fmt.Println("Healthcheck failed, status: ", res.StatusCode)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Println("Healthcheck failed, status: ", res.StatusCode)
		os.Exit(1)
	}
	fmt.Println("Healthcheck success, status: ", res.StatusCode)
}
