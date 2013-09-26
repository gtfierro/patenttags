#!/usr/bin/env python
"""
Takes the output of pattag.go, which is a CSV file following the structure:
Patent ID, Cluster ID
and outputs histograms on the Nth largest cluster found
"""

import sys
import csv
from BeautifulSoup import BeautifulSoup as bs
from collections import defaultdict, Counter
from matplotlib import pyplot as plt
import eventlet
from eventlet.green import urllib2  
import cPickle as pickle

BASE_URL = 'http://www.google.com/patents/'

def fix_patent_number(patent):
        if patent[:2].upper() != 'US':
            patent = 'US'+patent
        return patent

def get_class(patent_id, cluster_id):
    url = BASE_URL + fix_patent_number(patent_id)
    opener = urllib2.build_opener()
    opener.addheaders = [('User-agent', 'Mozilla/5.0')]
    try:
        html = opener.open(url).read()
        soup = bs(html)
        toptd = soup.find('td',text='U.S. Classification')
        classstring = toptd.findNext('td').findNext('span').text
        mainclass, subclass = classstring.split('/')
        return patent_id, cluster_id, mainclass, subclass
    except:
        print 'Couldnt get', patent_id
        return ('', '', '', '')

def initialize(filename):
    patents = []
    counter = Counter()
    for line in csv.reader(open(filename)):
        patent_id = line[0]
        cluster_id = line[1]
        patents.append((patent_id, cluster_id))
        counter.update(cluster_id)
    return patents, counter

def construct_metadata(filename):
    metadata = defaultdict(dict)
    i = 0
    patents, counter = initialize(filename)
    pool = eventlet.GreenPool(size=200)
    for patent_id, cluster_id, mainclass, _ in pool.imap(get_class, *zip(*patents)):
        i += 1
        if i % 1000 == 0:
            print i
        metadata[cluster_id][patent_id] = mainclass
    return metadata, counter

if __name__ == "__main__":
    filename = sys.argv[1]
    metadata, counter = construct_metadata(filename)
    pickle.dump(metadata, open('cleantech_metadata.pickle', 'w'))
    pickle.dump(counter, open('cleantech_counts.pickle', 'w'))
    biggest_clusterid = counter.most_common()[0]
    plt.hist(metadata[biggest_clusterid].itervalues())
    plt.savefig('class_histogram.png')
