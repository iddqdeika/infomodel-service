package main

import (
	"github.com/iddqdeika/rrr"
	"infomodel-service/root"
)

func main() {
	r := root.New()
	rrr.BasicEntry(r)
}
