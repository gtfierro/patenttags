package main

import (
    "os"
    "fmt"
	"github.com/gtfierro/patentcluster"
    "log"
    "bufio"
    "sync"
    "strconv"
)

var wg sync.WaitGroup

type Record struct {
    patent *patentcluster.Patent
    similarity float64
}

func compare(focal, patent *patentcluster.Patent, before, after chan *Record) {
    similarity := focal.JaccardDistance(patent)
    r := new(Record)
    r.patent = patent
    r.similarity = similarity
    if patent.AppDate.Before(focal.AppDate) {
        before <- r
    } else {
        after <- r
    }
    wg.Done()
}

func loop(focal *patentcluster.Patent, patents [](*patentcluster.Patent)) (int, int) {
    beforeRecords := [](*Record){}
    afterRecords := [](*Record){}
    before := make(chan *Record)
    after := make(chan *Record)
    done := make(chan bool)
    go func() {
        for {
            select {
            case r := <- before:
                beforeRecords = append(beforeRecords, r)
            case r := <- after:
                afterRecords = append(afterRecords, r)
            case <-done:
                fmt.Println(len(beforeRecords), len(afterRecords))
                break
            }
        }
    }()
    for _, patent := range patents {
        wg.Add(1)
        go compare(focal, patent, before, after)
    }
    wg.Wait()
    done <- true
    numbefore := writeToFile(beforeRecords, "before.csv")
    numafter := writeToFile(afterRecords, "after.csv")
    return numbefore, numafter
}

func writeToFile(toWrite [](*Record), outfilename string) int {
    outfile, err := os.Create(outfilename)
    if err != nil {
        log.Fatal("Could not output to file %s", outfilename)
    }
    defer outfile.Close()
    writer := bufio.NewWriter(outfile)
    for _, record := range toWrite {
        line := record.patent.Number + "," + strconv.FormatFloat(record.similarity, 'g', -1, 64) + "\n"
        writer.WriteString(line)
    }
    writer.Flush()
    return len(toWrite)
}

func search(focalnumber string, patents [](*patentcluster.Patent)) *patentcluster.Patent {
    for _, patent := range patents {
        if patent.Number == focalnumber {
            return patent
        }
    }
    return nil
}

func main () {
    patentfile := os.Args[1]
    focalnumber := os.Args[2]
	patents := patentcluster.Extract_file_contents(patentfile, true)
    focalpatent := search(focalnumber, patents)
    if focalpatent == nil {
        log.Fatal("Could not find patent number %s", focalnumber)
    }
    numbefore, numafter := loop(focalpatent, patents)
    fmt.Println(numbefore,"before")
    fmt.Println(numafter,"after")
}
