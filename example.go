package main

import (
	"fmt"
	"github.com/pims/gosnowizard/snowizard"
	"log"
	"time"
)

func main() {

	hosts := make([]string, 2)
	hosts[0] = "snowizard-1.dev:6776"
	hosts[1] = "snowizard-2.dev:6776"

	timeout := time.Duration(2 * time.Second)
	client := snowizard.NewSnowizardTextClient(hosts, timeout)

	id, err := client.Next()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)

	pb_client := snowizard.NewSnowizardProtobufClient(hosts, timeout)
	pb_id, pb_err := pb_client.Next()
	if pb_err != nil {
		log.Fatal(pb_err)
	}

	fmt.Println(pb_id)

	json_client := snowizard.NewSnowizardJsonClient(hosts, timeout)
	json_id, json_err := json_client.Next()
	if json_err != nil {
		log.Fatal(json_err)
	}
	fmt.Println(json_id)
}
