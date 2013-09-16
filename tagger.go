package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
    "time"
)

type Patent struct {
	number string
	tags   []int
}

var tagset = make(map[string]int)

const NCPU = 2

var vectorfile, err = os.Create("tagcount.txt")
var w = bufio.NewWriter(vectorfile)

/** enumerates all tags in the taglist and inserts them into `tagset */
func extract_tags(taglist string) {
	tags := strings.Split(taglist, " ")
	for _, s := range tags {
		tagset[s] += 1
	}
}

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
}

/** takes in a taglist for a certain patent and returns a vector
 * projected into the large tag space */
func create_tag_vector(number, taglist string, cs chan *Patent) {
	fmt.Println(number)
	tags := strings.Split(taglist, " ")
	vector := []int{}
	index := 0
	for tag, _ := range tagset {
		for i, val := range tags {
			if tag == val {
				vector = append(vector, 1)
                tags = append(tags[:i], tags[i+1:]...)
				break
			}
			vector = append(vector, 0)
		}
		index += 1
	}
    p := new(Patent)
    p.number = number
    p.tags = vector
	cs <- p
}

func write_vector(cs chan *Patent) {
	p := <-cs
	number := p.number
	vector := p.tags
	fmt.Fprintln(w, number, vector)
}

func process_tags() {
	datafile, _ := os.Open("buzzx.csv")
	r := csv.NewReader(datafile)
	cs := make(chan *Patent)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		number := record[0]
		tags := record[3]
		go create_tag_vector(number, tags, cs)
		go write_vector(cs)
        time.Sleep(4 * 1e9)
	}
}

func main() {
	read_file()
    fmt.Println("done reading file")
	process_tags()
	w.Flush()
}
