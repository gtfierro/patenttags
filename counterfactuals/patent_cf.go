package main

import (
    "time"
    "encoding/csv"
    "fmt"
    "os"
    "io"
	"github.com/gtfierro/patentcluster"
    "code.google.com/p/gcfg"
    "bufio"
)

type CFG struct {
    Config struct {
        WindowBefore string
        WindowAfter string
        FocalDate string
        DistanceThreshold float64
        DistanceMetric string
        Patentcorpus string
        DateLayout string
    }
}

func getDates(cfg CFG) (time.Time, time.Time, time.Time) {
    windowbefore, _ := time.Parse(cfg.Config.DateLayout, cfg.Config.WindowBefore)
    windowafter, _ := time.Parse(cfg.Config.DateLayout, cfg.Config.WindowAfter)
    focaldate, _ := time.Parse(cfg.Config.DateLayout, cfg.Config.FocalDate)
    return windowbefore, windowafter, focaldate

}

func readconfig(configfile string) CFG {
    fmt.Println("Reading configuration from", configfile)
    cfg := CFG{}
    err := gcfg.ReadFileInto(&cfg, configfile)
    if err != nil {
        fmt.Println("Failed to parse config file", configfile, ":", err)
    }
    return cfg
}

/**
Params:
earlydate, latedate: specify the lower and upper bounds on the window we want to search
threshold: specifies the maximum distance between two patents before they can be considered similar
search: patent we're comparing against
patents: the full list of patents we're searching in

Returns:
list of patents that match the above criteria
*/
func findPatents(earlydate, latedate time.Time, threshold float64, search *patentcluster.Patent, patents [](*patentcluster.Patent)) [](*patentcluster.Patent) {
    res := [](*patentcluster.Patent){}
    for _, patent := range patents {
        if patent == search {
            continue
        }
        /* check that the patent is within the requisite range */
        if patent.AppDate.After(latedate) || patent.AppDate.Before(earlydate) {
            continue
        }
        if search.JaccardDistance(patent) <= threshold {
            res = append(res, patent)
        }
    }
    return res
}

//windowbefore, windowafter, focaldate := getDates(config)
func FindCounterFactuals(windowbefore, windowafter, focaldate time.Time, threshold float64, searchlist, patents [](*patentcluster.Patent)) (map[string][]string, map[string][]string) {
    before := make(map[string][]string)
    after := make(map[string][]string)
    for _, search := range searchlist {
        // search before
        for _, patent := range findPatents(windowbefore, focaldate, threshold, search, patents) {
            before[search.Number] = append(before[search.Number], patent.Number) 
        }
        // search after
        for _, patent := range findPatents(focaldate, windowafter, threshold, search, patents) {
            after[search.Number] = append(after[search.Number], patent.Number) 
        }
    }
    return before, after
}

/**
Reads the file from `patentfile`, where each line should be a Patent ID
*/
func getPatentList(patentfile string, patents [](*patentcluster.Patent)) [](*patentcluster.Patent) {
    datafile, err := os.Open(patentfile)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }
    defer datafile.Close()
    reader := csv.NewReader(datafile)
    searchids := [](*patentcluster.Patent){}
    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return nil
        }
        for _, patent := range patents {
            if patent.Number == line[0] {
                searchids = append(searchids, patent)
                break
            }
        }
    }
    return searchids
}

func To_file(filename string, windowbefore, windowafter, focaldate time.Time, before, after map[string][]string) {
    outfile, err := os.Create(filename)
    if err != nil {
        fmt.Println("Could not output to file", filename, ":", err)
        return
    }
    defer outfile.Close()
    writer := bufio.NewWriter(outfile)
    header := "Between "+windowbefore.Format("Jan 02 2006")+" to "+ focaldate.Format("Jan 02 2006")+"\n"
    writer.WriteString(header)
    for number, list := range before {
        line := number + ","
        for _, patent := range list {
            line += " " + patent
        }
        line += "\n"
        writer.WriteString(line)
    }
    header = "Between "+focaldate.Format("Jan 02 2006")+" to "+ windowafter.Format("Jan 02 2006")+"\n"
    writer.WriteString(header)
    for number, list := range after {
        line := number + ","
        for _, patent := range list {
            line += " " + patent
        }
        line += "\n"
        writer.WriteString(line)
    }
    writer.Flush()
}

func main() {
	patentfile := os.Args[1]
    config := readconfig("config.ini")
    windowbefore, windowafter, focaldate := getDates(config)
	patents := patentcluster.Extract_file_contents(config.Config.Patentcorpus, true)
    searchlist := getPatentList(patentfile, patents)
    fmt.Println(len(searchlist), "patents to analyze")
	fmt.Println(len(patentcluster.Tagset), "unique tags")
	fmt.Println(len(patents), "unique patents")
    before, after := FindCounterFactuals(windowbefore, windowafter, focaldate, config.Config.DistanceThreshold, searchlist, patents)
    fmt.Println("Found",len(before), "before")
    fmt.Println("Found",len(after),"after")
    To_file("counterfactual.out", windowbefore, windowafter, focaldate, before, after)
}
