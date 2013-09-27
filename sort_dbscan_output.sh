#!/bin/bash

filename=$1
sort -t',' -k +2n $filename > dbscansorted.csv
