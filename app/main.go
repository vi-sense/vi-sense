package main

import (
	"fmt"
	. "github.com/vi-sense/vi-sense/app/api"
	_ "github.com/vi-sense/vi-sense/app/docs"
	. "github.com/vi-sense/vi-sense/app/model"
	"io/ioutil"
	"log"
)

func main() {
	SetupDatabase(true)
        CreateMockData("/sample-data", []string{"berlin", "cape-town", "puerto-natales", "sample-model"}, -1)

	//check if bind mount is working
	dat, err := ioutil.ReadFile("/sample-data/info.txt")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(dat))

	r := SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err = r.RunTLS(":44344", "/certs/live/visense.f4.htw-berlin.de/fullchain.pem", "/certs/live/visense.f4.htw-berlin.de/privkey.pem")
	if err != nil {
		log.Println(err)
		log.Fatal(r.Run(":8080"))
	}
}
