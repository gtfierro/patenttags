import sys
from matplotlib import pyplot as plt
from csv_reader import read_file
import numpy as np
from collections import Counter
from matplotlib.ticker import FormatStrFormatter

def makeHistogram(filename, outname):
    data = [float(x[1]) for x in read_file(filename)]
    fig, ax = plt.subplots()
    plt.hist(data, log=True)
    counts, bins, patches = ax.hist(data)
    ax.set_xticks(bins)
    ax.xaxis.set_major_formatter(FormatStrFormatter('%0.1f'))
    bin_centers = 0.5 * np.diff(bins) + bins[:-1]
    for count, x in zip(counts, bin_centers):
        # Label the raw counts
        ax.annotate(str(count), xy=(x, 0), xycoords=('data', 'axes fraction'),
            xytext=(0, -18), textcoords='offset points', va='top', ha='center')
        # Label the percentages
        percent = '%0.3f%%' % (100 * float(count) / counts.sum())
        ax.annotate(percent, xy=(x, 0), xycoords=('data', 'axes fraction'),
            xytext=(0, -32), textcoords='offset points', va='top', ha='center',
            rotation=45)
        # Give ourselves some more room at the bottom of the plot
        plt.subplots_adjust(bottom=0.30)
    plt.xlabel("Jaccard Distance (1.0 is no shared tags)", labelpad=60)
    plt.ylabel("Number of Patents")
    plt.savefig(outname+'.png')

if __name__=="__main__":
    beforefile, afterfile = sys.argv[1], sys.argv[2]
    makeHistogram(beforefile, 'before')
    makeHistogram(afterfile, 'after')
