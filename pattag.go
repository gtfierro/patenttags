package main

import (
	"fmt"
    "strings"
    "os"
    "github.com/gtfierro/patentcluster"
    "runtime/pprof"
    "log"
    "flag"
)
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")


func practice_run(db *patentcluster.DBSCAN, patents [](*patentcluster.Patent)) {
    for _, pat := range patents {
        seeds := db.RegionQuery(pat)
        if len(seeds) > 0 {
            fmt.Println(len(seeds))
        }
    }
}

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    filename := os.Args[1]
    rootname := strings.Split(filename, ".")[0]
    fmt.Println("reading from file", filename)
    patentcluster.Read_file(filename, true)
    fmt.Println(len(patentcluster.Tagset), "unique tags")
    patents := patentcluster.Make_patents(filename, true)
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
    db.To_file(rootname+".out")
}
