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
    fmt.Println("Initializing DBSCAN...")
    db := patentcluster.Init_DBSCAN(patents, .9, 3)
    fmt.Println("Running DBSCAN...")
    db.Run()
    fmt.Println("Finished running DBSCAN. Computing Stats...")
    num_clusters, mean_size, median_size, largest, _ := db.Compute_Stats()
    fmt.Println("Number of clusters found:", num_clusters)
    fmt.Println("Mean cluster size:", mean_size)
    fmt.Println("Median cluster size:", median_size)
    fmt.Println("Largest cluster size:", largest)
    db.To_file("buzz5000.out")
}
