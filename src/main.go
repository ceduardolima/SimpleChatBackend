package main

import (
	"SimpleChat/src/router"
	"fmt"
)

func main() {

	router := router.CreateRouter()
	fmt.Println("Starting server...")

	err := router.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

}
