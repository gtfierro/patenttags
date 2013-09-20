package main

import (
	"fmt"
    "github.com/gtfierro/patentcluster"
)

func practice_run(db *patentcluster.DBSCAN, patents [](*patentcluster.Patent)) {
    for _, pat := range patents {
        seeds := db.RegionQuery(pat)
        if len(seeds) > 0 {
            fmt.Println(len(seeds))
        }
    }
}

func main() {
    fmt.Println("reading from file")
    patentcluster.Read_file("buzz500.csv")
    fmt.Println(len(patentcluster.Tagset), "unique tags")
    patents := patentcluster.Make_patents("buzz5000.csv")
    fmt.Println(len(patents), "unique patents")
    db := patentcluster.Init_DBSCAN(patents, .9, 3)
    db.Run()
    db.To_file("buzz5000.out")
}
