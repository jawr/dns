package main

import (
	"flag"
	"fmt"
	"github.com/jawr/dns/database/models/users"
)

func main() {
	user := flag.String("user", "", "username")
	pass := flag.String("pass", "", "password")
	flag.Parse()
	if len(*user) == 0 || len(*pass) == 0 {
		fmt.Println("Need to pass a username and password")
		return
	}
	_, err := users.New(*user, *pass)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User %s created.\n", *user)
}
