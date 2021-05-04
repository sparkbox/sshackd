package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func ExampleConn_WhoAmI() {
	conn, err := ldap.DialURL("ldap://ldap.forumsys.com:389")
	if err != nil {
		log.Fatalf("Failed to connect: %s\n", err)
	}

	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: "uid=euler,dc=example,dc=com",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Failed to bind: %s\n", err)
	}

	res, err := conn.WhoAmI(nil)
	if err != nil {
		log.Fatalf("Failed to call WhoAmI(): %s\n", err)
	}
	fmt.Printf("I am: %s\n", res.AuthzID)
}

func main() {
	ExampleConn_WhoAmI()
}
