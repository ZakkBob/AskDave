package main

import (
	"github.com/ZakkBob/AskDave/backend/orm"
)

func main() {
	orm.Connect("")
	defer orm.Close()
}
