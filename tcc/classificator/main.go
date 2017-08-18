package main

import (
	"github.com/julianogalgaro/classificator/api"
	"github.com/julianogalgaro/classificator/control"
)

func main() {

	go control.NewControl("http://localhost:8081/").StartPredict()

	api.NewApi("80").StartServer()

}
