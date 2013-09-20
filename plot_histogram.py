import sys
from matplotlib import pyplot as plt

datafile = open(sys.argv[1])

def deliver():
  for line in datafile:
    val = float(line.split(' ')[0])
    yield val

data = list(deliver())

plt.hist(data, log=True)
plt.savefig('hist.png')
