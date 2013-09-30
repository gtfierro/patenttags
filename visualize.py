import cPickle as pickle
import sys
import numpy as np
from scipy import sparse
from csv_reader import read_file
from collections import defaultdict
import matplotlib.pyplot as plt
import matplotlib.colors as colors
import matplotlib.cm as cmx

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
    cluster = line[1]
    for tag in mytags:
        index = tagset[tag]
        patent[index, 0] = 1
    x = (patent.transpose() * e1)[0][0]
    y = (patent.transpose() * e2)[0][0]
    if dimensions == 2:
        return cluster, x, y
    if dimensions == 3:
        z = (patent.transpose() * e3)[0][0]
        return cluster, x, y, z

def plot_patents(lines, tagset, e1, e2, e3=None, dimensions=2):
    clusters = []
    plots = defaultdict(lambda: defaultdict(list))
    for line in lines:
        cluster,x,y = get_coordinates(line, tagset, e1, e2)
        plots[cluster]['xcoords'].append(x)
        plots[cluster]['ycoords'].append(y)
        clusters.append(cluster)
    cm = plt.get_cmap('gist_rainbow')
    colors = [cm(1.*i/len(set(clusters))) for i in xrange(len(set(clusters)))]
    colormap = dict([(cluster, color) for cluster, color in zip(set(clusters), colors)])
    scatters = []
    for cluster in plots.iterkeys():
        p = plt.scatter(plots[cluster]['xcoords'], plots[cluster]['ycoords'], color=colormap[cluster])
        scatters.append(p)
    lgd = plt.legend(scatters, plots.iterkeys(), loc='upper center', bbox_to_anchor=(0.5, -0.1))
    plt.savefig('out.png', bbox_extra_artists=(lgd,), bbox_inches='tight')


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
