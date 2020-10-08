package main

import (
	"fmt"
	"os"

	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

// If someone should become a new administrator, then please add here
var admin = []string{
	"gslee@bighitcorp.com",
}

const region = "ap-northeast-2"

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("usage: go run hack/release/check_permission.go <email account> <bucket name>")
		os.Exit(1)
	}

	if !tools.IsStringInArray(args[1], admin) {
		fmt.Printf("you are not authorized: %s\n", args[1])
		os.Exit(1)
	}

	if err := checkBucket(args[2]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("You are an authorized administrator for release")
}

func checkBucket(bucket string) error {
	return nil
	//sess := client.GetAwsSession()
	//client := client.NewClient(sess, region, nil)
	//
	//return client.HeadS3Bucket(bucket)
}
