package main

import (
	"github.com/ZakkBob/AskDave/crawler/daveapiclient"
)

func main() {
	c := daveapiclient.Create("http://localhost:8080")

	c.Run()
}
