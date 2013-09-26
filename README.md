# Patent Tags

This repository is a set of scripts meant to facilitate the processing of patent tags

* `pattag.go` 
    * takes in a tag file, as found in
      [`cleantech_data/raw_tags.csv`](https://github.com/gtfierro/patenttags/blob/master/cleantech_data/raw_tags.csv)
      and outputs a CSV file of format `patent_id, cluster_id, space separated
      list of tags`. This output file is used for the other helper scripts.
* `cross_validate.go`
    * takes in a tag file like `pattag.go` does, and outputs a CSV file of statistics for various
      runs of the DBSCAN algorithm over several different combinations of parameters.
* `csv_reader.py`
    * helper library for opening CSV files as iterators. Handles Unicode files.
* `get_metadata.py` 
    * Takes as input the output CSV file from `pattag` and accumulates
      mainclass/subclass for each of the patents (by scraping Google) as well
      as counts for how often each of the tags appears
* `plot_stats.py`
    * Run `python plot_stats.py -h` for command line arguments. Takes as input
      the output CSV file from `cross_validate` and creates several graphs
      showing various statistics of the found clusters.
* `patent_PCA.py`
    * takes as input the output CSV file from `pattag` and computes the top 3 eigenvalues and vectors
      of the tagged data. Saves relevant metadata and output in `eigen.pickle`
