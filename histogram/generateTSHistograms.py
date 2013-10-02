import sys
from matplotlib import pyplot as plt
import numpy as np
from collections import Counter
from matplotlib.ticker import FormatStrFormatter
from matplotlib import dates as mdates

def reject_outliers(data, m=4):
    mask = (abs(data - np.mean(data)) < m * np.std(data))
    return mask

def makeHistogram(filename, outname):
    days, scores = np.loadtxt(filename, unpack=True, delimiter=",",
                        converters={ 0: mdates.strpdate2num('%d-%b-%Y')})
    mask = reject_outliers(days)
    days = days[mask]
    scores = scores[mask]
    plt.plot_date(x=days, y=scores, marker='.')
    plt.xlabel("Patent Application Date")
    plt.ylabel("Jaccard Similarity (1.0 is all shared tags) ")
    plt.savefig(outname+'.png')

if __name__=="__main__":
    beforefile, afterfile = sys.argv[1], sys.argv[2]
    makeHistogram(beforefile, 'before')
    makeHistogram(afterfile, 'after')
