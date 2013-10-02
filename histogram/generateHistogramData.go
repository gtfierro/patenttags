package main

import (
    "os"
    "fmt"
	"github.com/gtfierro/patentcluster"
    "log"
    "bufio"
    "sync"
    "strconv"
    "math/rand"
)

var wg sync.WaitGroup

type Record struct {
    patent *patentcluster.Patent
    similarity float64
}

func compare(focal, patent *patentcluster.Patent, before, after chan *Record) {
    similarity := focal.JaccardSimilarity(patent)
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

func getKMedian(records [](*Record), k int) float64 {
    a := rand.Intn(len(records))
    pivot := records[a]
    smaller := [](*Record){}
    bigger := [](*Record){}
    for _, record := range records {
        if record.similarity < pivot.similarity {
            smaller = append(smaller, record)
        } else {
            bigger = append(bigger, record)
        }
    }
    if len(bigger) == len(records) || len(smaller) == len(records) {
        return pivot.similarity
    }
    if len(smaller) == k-1 {
        return pivot.similarity
    } else {
        if len(smaller) < k-1 {
            return getKMedian(bigger, k-len(smaller)-1)
        } else {
            return getKMedian(smaller, k)
        }
    }
}

func getMedian(records [](*Record)) float64 {
    numrecords := len(records)
    medianidx := numrecords / 2
    return getKMedian(records, medianidx)
}

func getMean(records [](*Record)) float64 {
    sum := 0.0
    for _, record := range records {
        sum += record.similarity
    }
    return sum / float64(len(records))

}

func loop(focal *patentcluster.Patent, patents [](*patentcluster.Patent)) ([](*Record), [](*Record)) {
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
    writeTSToFile(beforeRecords, "beforeTS.csv")
    writeTSToFile(afterRecords, "afterTS.csv")
    fmt.Println(numbefore,"before")
    fmt.Println(numafter,"after")
    return beforeRecords, afterRecords
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

/**
    Writes records to file as timeseries data:
    application date, similarity score
*/
func writeTSToFile(toWrite [](*Record), outfilename string) {
    outfile, err := os.Create(outfilename)
    if err != nil {
        log.Fatal("Could not output to file %s", outfilename)
    }
    defer outfile.Close()
    writer := bufio.NewWriter(outfile)
    for _, record := range toWrite {
        line := record.patent.AppDate.Format("02-Jan-2006")+ "," + strconv.FormatFloat(record.similarity, 'g', -1, 64) + "\n"
        writer.WriteString(line)
    }
    writer.Flush()
}

func search(focalnumber string, patents [](*patentcluster.Patent)) *patentcluster.Patent {
    for _, patent := range patents {
        if patent.Number == focalnumber {
            return patent
        }
    }
    return nil
}

/**
Returns number of items in list that are not 1.0
*/
func numberSignificant(list [](*Record)) int {
    num := 0
    for _, record := range list {
        if record.similarity < 1.0 {
            num += 1
        }
    }
    return num
}

func main () {
    patentfile := os.Args[1]
    focalnumber := os.Args[2]
	patents := patentcluster.Extract_file_contents(patentfile, true)
    focalpatent := search(focalnumber, patents)
    if focalpatent == nil {
        log.Fatal("Could not find patent number %s", focalnumber)
    }
    before, after := loop(focalpatent, patents)
    fmt.Println("Median of 'before' records:", getMedian(before))
    fmt.Println("Median of 'after' records:", getMedian(after))
    fmt.Println("Mean of 'before' records:", getMean(before))
    fmt.Println("Mean of 'after' records:", getMean(after))
    fmt.Println("Non zero of 'before' records:", numberSignificant(before))
    fmt.Println("Non zero of 'after' records:", numberSignificant(after))
}
