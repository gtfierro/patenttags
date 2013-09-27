import cPickle as pickle
import sys
import numpy as np
from scipy import sparse
from csv_reader import read_file
import matplotlib.pyplot as plt

def get_coordinates(line, tagset, e1, e2, e3=None, dimensions=2):
    """
    Given a line representing a patent, projects it onto the eigenvectors
    to obtain coordinates.

    Params:
    line: 3-tuple containing patent_id, cluster_id, space separated list of tags
    tagset: dictionary of tags to index. This is for maintaining tag order in vector
            representations of patents
    e1, e2, e3: numpy.ndarrays representing the top3 eigenvectors
    dimensions: number of dimensions we want to project our line into
    """
    patent = sparse.lil_matrix((len(tagset.keys()), 1))
    mytags = line[2].split(' ')
    for tag in mytags:
        index = tagset[tag]
        patent[index, 0] = 1
    x = (patent.transpose() * e1)[0][0]
    y = (patent.transpose() * e2)[0][0]
    if dimensions == 2:
        return x, y
    if dimensions == 3:
        z = (patent.transpose() * e3)[0][0]
        return x,y,z

def plot_patents(lines, tagset, e1, e2, e3=None, dimensions=2):
    xcoords = []
    ycoords = []
    labels = []
    for line in lines:
        x,y = get_coordinates(line, tagset, e1, e2)
        xcoords.append(x)
        ycoords.append(y)
        labels.append(line[0])
    plt.scatter(xcoords, ycoords)
    for i, label in enumerate(labels):
        plt.text(xcoords[i], ycoords[i], label)
    plt.savefig('out.png')


if __name__=='__main__':
    inputfile = sys.argv[1]
    lines = read_file(inputfile)
    eigenfile = sys.argv[2]
    eigenstuffs = pickle.load(open(eigenfile))
    cluster_ids = sys.argv[3].split(',') if len(sys.argv) > 3 else None
    if cluster_ids:
        lines = filter(lambda x: x[1] in cluster_ids, lines)
    tagset = eigenstuffs['tagset']
    e1 = eigenstuffs[1][1]
    e2 = eigenstuffs[2][1]
    e3 = eigenstuffs[3][1]
    print 'Top2 Variance', eigenstuffs['top2_variance']
    print 'Top3 Variance', eigenstuffs['top3_variance']
    plot_patents(lines, tagset, e1, e2)
