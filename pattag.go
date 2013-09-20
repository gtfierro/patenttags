package main

import (
	"fmt"
    "github.com/gtfierro/patentcluster"
)

func main() {
    fmt.Println("reading from file")
    patentcluster.Read_file("buzzx.csv")
}
