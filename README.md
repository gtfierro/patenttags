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
* `visualize.py`

## Sample Pipeline

Receive tagging file `buzzx.csv`. Run DBSCAN algorithm to obtain clusters:

```
go run pattag.go buzzx.csv
```

You will receive as output the file `buzzx.out`. Sort the file by clusters:

```
./sort_dbscan_output.sh buzzx.out
```

This will create a file `dbscansorted.csv`. Compute the eigenvalues/vectors for
a patent by running. The final argument is an optional specification of a
cluster identifier.  If a cluster identifier is specified, the PCA will only be
performed for patents in that cluster.

```
python patent_PCA.py dbscansorted.csv 5002618
```

The output is the file `eigen.pickle` which is a Python cPickle file of relevant
metadata about the eigenvectors/values of the PCA process. To visualize, run

```
python visualize.py dbscansorted.csv eigen.pickle 5002618
```

Again, the final argument is optional. If you specified it before, though, you
should include it again, or you will get weird results.
