#!/usr/bin/env python
"""
Given output of the DBSCAN algorithm, which is a CSV file with schema:

    patent_id, cluster_id, space separated list of tags for patent_id

we want to facilitate the visualization of the clusters by performing PCA
on the matrix of tags and then projecting each patent onto the 3 eigenvectors
associated with the top 3 eigenvalues to obtain (x,y,z) coordinates.
"""

import sys
from csv_reader import read_file
from scipy import sparse
from scipy.sparse import linalg
import numpy as np
import cPickle as pickle

def construct_tagset(lines, cluster_ids=None):
    """
    Constructs a Python set() of all the tags
    
    Params:
    lines: an iterable of the lines in the input file

    Returns:
    2-tuple of:
    - a python dict of key = tags, value = index. This will allow for fast searching
    - number of lines in the file
    """
    tagset = set()
    num_lines = 0
    for line in lines:
        # option to filter to only patents from a single cluster
        if cluster_ids and line[1] not in cluster_ids: continue
        num_lines += 1
        tags = line[3].split(' ')
        tagset.update(tags)
    tmp = enumerate(list(tagset))
    return dict((y,x) for x,y in tmp), num_lines

def insert_sparse_patent(line, tagset, patents, col_idx):
    """
    Inserts the patent contained in `line` into the sparse matrix
    of all patents

    Params:
    line: 3-tuple of patent_id, cluster_id, space separated list of tags
    tagset: output of construct_tagset
    patents: scipy.sparse.lil_matrix we want to insert our patent into
    col_idx: index of column we're inserting into, assuming 0-indexed

    Returns:
    scipy.sparse.lil_matrix containing the patent from `line`
    """
    mytags = line[3].split(' ')
    for tag in mytags:
        index = tagset[tag]
        patents[index, col_idx] = 1
    return patents

def construct_sparse_patent_matrix(lines, tagset, cluster_ids=None):
    patents = sparse.lil_matrix((num_tags, num_lines))
    if cluster_ids:
        lines = filter(lambda x: x[1] in cluster_ids, lines)
    for i, line in enumerate(lines):
        patents = insert_sparse_patent(line, tagset, patents, i)
    print patents.sum(1)
    return patents

def demean(matrix):
    """
    Returns matrix, but all columns are 0-mean
    """
    # compute column-wise mean
    col_mean = matrix.mean(axis=1)
    matrix = matrix - col_mean
    return matrix

def get_n_eigenstuffs(patents, n):
    """
    Return top N eigenvalues and corresponding eigenvectors

    Returns 
    """
    patentcsr = patents.tocsr()
    patentcsr_t = patentcsr.transpose()

    # compute covariance matrix
    print 'Computing covariance matrix...'
    cov = patentcsr_t * patentcsr
    print 'Converting to zero-mean...'
    cov = demean(cov)
    print 'Computing eigenthings...'
    num_eigs = cov.shape[1]-1 # cov is symmetric, so rank is height-1
    eigs = linalg.eigsh(cov, num_eigs) # compute all eigenvalues
    yield eigs[0].sum()
    eigenvalues = eigs[0][::-1] # reverse order
    eigenvectors = eigs[1][::-1] # reverse order, then get top n
    for i, eigenvalue in enumerate(eigenvalues):
        yield eigenvalue, patentcsr * eigenvectors[:, i]


if __name__ == '__main__':
    inputfile = sys.argv[1]
    cluster_ids = sys.argv[2].split(',') if len(sys.argv) > 2 else None
    lines = read_file(inputfile)
    print 'Constructing tagset...'
    tagset, num_lines = construct_tagset(lines, cluster_ids)
    num_tags = len(tagset.keys())
    lines = read_file(inputfile)
    print 'Constructing sparse matrix...'
    patents = construct_sparse_patent_matrix(lines, tagset, cluster_ids)
    print 'Doing PCA...'
    e = get_n_eigenstuffs(patents, 3)
    eigensum = e.next()
    e1 = e.next()
    e2 = e.next()
    e3 = e.next()
    d = {1: (e1[0], e1[1].reshape(e1[1].shape[0], 1)),
         2: (e2[0], e2[1].reshape(e2[1].shape[0], 1)),
         3: (e3[0], e3[1].reshape(e3[1].shape[0], 1)),
         'tagset': tagset,
         'top2_variance': (e1[0] + e2[0]) / float(eigensum),
         'top3_variance': (e1[0] + e2[0] + e3[0]) / float(eigensum)}
    pickle.dump(d, open('eigen.pickle','wa'))
