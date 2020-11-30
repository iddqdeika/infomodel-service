package main

import (
	"github.com/iddqdeika/infomodel-service/root"
	"github.com/iddqdeika/rrr"
)

func main() {
	r := root.New()
	rrr.BasicEntry(r)
}
