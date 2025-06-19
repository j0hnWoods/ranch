package main

import (
	"fmt"
	"ranch/ranch"
)

func main() {
	controller := ranch.RanchController{}
	middleware := ranch.Middleware{}
	service := ranch.RanchService{}
	render := ranch.RanchRender{}

	fmt.Printf("%v\n%v\n%v\n%v\n", controller, middleware, service, render)
}
