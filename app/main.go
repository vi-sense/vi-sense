package main

import (
	"fmt"
	. "github.com/vi-sense/vi-sense/app/api"
	_ "github.com/vi-sense/vi-sense/app/docs"
	. "github.com/vi-sense/vi-sense/app/model"
	"github.com/adrianosela/sslmgr"
	"github.com/adrianosela/certcache"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
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
	ss, err := sslmgr.NewServer(sslmgr.ServerConfig{
		Hostnames: []string{"visense.f4.htw-berlin.de:44344"},
		HTTPPort:  ":" + os.Getenv("PORT"),
		HTTPSPort: ":" + os.Getenv("HTTPS_PORT"),
		Handler:   r,
		ServeSSLFunc: func() bool {
			return strings.ToLower(os.Getenv("PROD")) == "true"
		},
		CertCache: certcache.NewLayered(
			certcache.NewLogger(),
			autocert.DirCache("/certcache"),
		),
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		IdleTimeout:         25 * time.Second,
		GracefulnessTimeout: 5 * time.Second,
		GracefulShutdownErrHandler: func(e error) {
			log.Fatal(e)
		},
	})
	ss.ListenAndServe()
	if err != nil {
		panic(r)
	}
}