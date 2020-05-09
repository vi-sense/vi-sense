package main

import (
	"fmt"
	. "github.com/vi-sense/vi-sense/app/api"
	_ "github.com/vi-sense/vi-sense/app/docs"
	. "github.com/vi-sense/vi-sense/app/model"
	"io/ioutil"
)

func main() {
	SetupDatabase(true)
	CreateMockData("/sample-data", 1000)

	//check if bind mount is working
	dat, err := ioutil.ReadFile("/sample-data/info.txt")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(dat))

	r := SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err = r.Run(":8080")
	if err != nil {
		panic(r)
	}
}