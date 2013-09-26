import sys
import argparse
from matplotlib import pyplot as plt
from csv_reader import read_file
parser = argparse.ArgumentParser()
parser.add_argument("-e", "--epsilon", type=float,
                    help="Graph points with this epsilon")
parser.add_argument("-m", "--minpts", type=int,
                    help="Graph points with this minpts")
parser.add_argument("filename", help="Name of file containing stats")

def get_data(filename, epsilon=None, minpts=None):
    lines = read_file(filename)
    lines.next() # skip header
    means = []
    medians = []
    largest = []
    numclusters = []
    xpoints = []
    for line in lines:
        x = ''
        if epsilon:
            if float(line[0]) != epsilon: continue
            else: x = line[1]
        if minpts:
            if float(line[1]) != minpts: continue
            else: x = line[0]
        xpoints.append(x)
        numclusters.append(float(line[2]))
        means.append(float(line[3]))
        medians.append(float(line[4]))
        largest.append(float(line[5]))
    print xpoints, numclusters, means, medians, largest
    return xpoints, numclusters, means, medians, largest

def plot(fig, x, y, plotname, epsilon=None, minpts=None):
    fig.plot(x, y)
    fig.set_title(plotname)
    if epsilon:
        xlabel = 'Minpts, with epsilon={0}'.format(epsilon)
    if minpts:
        xlabel = 'Epsilon, with minpts={0}'.format(minpts)


if __name__=='__main__':
    args = parser.parse_args()
    epsilon = args.epsilon if args.epsilon else None
    minpts = args.minpts if args.minpts else None
    filename = args.filename
    if not epsilon and not minpts:
        print 'You must specify at least epsilon or minpts'
        sys.exit(0)
    f, axarr = plt.subplots(2, 2)
    xpoints, numclusters, means, medians, largest = get_data(filename, epsilon, minpts)
    plot(axarr[0, 0], xpoints, numclusters, 'Number of Clusters', epsilon, minpts)
    plot(axarr[0, 1], xpoints, means, 'Mean Cluster Size', epsilon, minpts)
    plot(axarr[1, 0], xpoints, medians, 'Median Cluster Size', epsilon, minpts)
    plot(axarr[1, 1], xpoints, largest, 'Largest Cluster Size', epsilon, minpts)
    f.set_figheight(10)
    f.set_figwidth(10)
    if epsilon:
        f.suptitle('Epsilon={0} X-axis:minpts'.format(epsilon))
        f.savefig('Epsilon_{0}.png'.format(epsilon))
    elif minpts:
        f.suptitle('Minpts={0} X-axis:epsilon'.format(minpts))
        f.savefig('Minpts_{0}.png'.format(minpts))
    
