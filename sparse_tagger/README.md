#Sparse Tagger

We're given an input file that's the raw output of Guan-Cheng's tagger

```
1 A A B A C D ...
2 B A C F F E ...
...
```

where the number indicates the patent record, and the letters represent tags
extracted by the patent. This will form the basis of a tag matrix:

```
         ---------Tags--------
         ___A___B___C___D___E___F___...
         1| 3   1   1   1   0   0
         2| 1   1   1   0   1   1   
patents   |...
       ...|...
```

This matrix is *very* sparse, so we can represent it with 3-tuples:
`(i,j,occ)`, which means that tag `j` appears in patent `i` `occ` number of
times. For above, we'd have the following which can be written to a CSV file 'matrix.csv':

```
(1, 1, 3),
(1, 2, 1),
(1, 3, 1),
(1, 4, 1),
(2, 1, 1),
(2, 2, 1), ...
```

Another output file 'dict.csv' will be generated that will hold the mappings of the tag indexes:

```
1, A
2, B
3, C
...
```
