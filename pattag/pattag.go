package main

import (
  "math"
  "strings"
  "os"
  "encoding/csv"
  "io"
  "fmt"
)

type Patent struct {
	number string
	tags   map[string]int
}

/** 
    Computes euclidian distance between two patents
    by taking the square root of the number of tags
    they do not have in common
*/
func (p *Patent) EuclidianDistance(target *Patent) float64 {
    var count float64
    count = 0
    for tag, _ := range target.tags {
        if p.tags[tag] > 0 { // if tag in common, we count it
           count += 1 
        }
    }
    // distance is the square root of the total number of tags not in common
    return math.Sqrt(number_of_tags - count)
}

/**
    Computes euclidian distance as above but normalizes to 0 - 1 
*/
func (p *Patent) NormalizedEuclidianDistance(target *Patent) float64 {
    return p.EuclidianDistance(target) / sqrt_num_tags
}

func (p *Patent) JaccardDistance(target *Patent) float64 {
    var count, union float64
    count = 0
    union = float64(len(p.tags))
    for tag, _ := range target.tags {
        if p.tags[tag] > 0 {
            count += 1
        } else {
            union += 1
        }
    }
    return 1 - count / union
}

/**
  given a string representing a patent number and
  a string representing the space-delimited list of
  tags for a patent, returns a reference to a Patent
  object
*/
func makePatent(number, tagstring string) *Patent {
    p := new(Patent)
    p.tags = make(map[string]int)
    for _, tag := range strings.Split(tagstring, " ") {
        p.tags[tag] = 1
    }
    p.number = number
    return p
}

var tagset = make(map[string]int)
var patents = [](*Patent){}
var number_of_tags float64
var sqrt_num_tags float64

/** enumerates all tags in the taglist and inserts them into `tagset */
func extract_tags(taglist string) {
	tags := strings.Split(taglist, " ")
	for _, s := range tags {
		tagset[s] += 1
	}
}

/**
    reads buzzx.csv and accumulates all patent tags
*/
func read_file() {
	/* open buzzx.csv file and start counting tags */
	datafile, err := os.Open("buzzx.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer datafile.Close()
	reader := csv.NewReader(datafile)
	number_of_records := 0
	/* loop through file */
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		number_of_records += 1
		tags := record[3]
		extract_tags(tags)
	}
    number_of_tags = float64(len(tagset))
    sqrt_num_tags = math.Sqrt(number_of_tags)
}

/**
    loops through buzzx.csv and creates a patent instance
    for each row
*/
func make_patents() {
	/* open buzzx.csv file and start counting tags */
	datafile, err := os.Open("buzzx.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer datafile.Close()
	reader := csv.NewReader(datafile)
	/* loop through file */
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
        number := record[1]
		tags := record[3]
        p := makePatent(number, tags)
        patents = append(patents, p)
	}
}


func main() {
    fmt.Println("Creating tag set...")
    read_file()
    fmt.Println("Accumulated", number_of_tags, "tags")
    fmt.Println("Done creating tag set!")
    fmt.Println("Making patent instances...")
    make_patents()
    fmt.Println("Finished", len(patents),"patent instances")
    p1 := patents[1]
    p2 := patents[2]
    fmt.Println(p1.EuclidianDistance(p2))
    fmt.Println(p1.NormalizedEuclidianDistance(p2))
}
