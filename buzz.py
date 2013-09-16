from csv_reader import read_file
from collections import Counter
from datetime import datetime
import numpy, scipy.io
import cPickle

# Here, we create the common set of tags and count how common
# these tags are among all our documents
tagset = set()
tagfrequencies = Counter()
lines = read_file('buzzx.csv')
lines.next() # skip CSV schema
for line in lines:
    number = line[0]
    tags = line[3].split(' ')
    tagfrequencies.update(tags)
    for tag in tags:
        tagset.add(tag)

# put all counts to 0
tagcounter = tagfrequencies
for item in tagcounter:
    tagcounter[item] = 0
print 'Done compiling tags'


def expand_tagspace(taglist):
    """
    Takes a list of tags and returns a vector in the
    space of the common taglist `tagset`:

    Args:
    @taglist: iterable containing strings of tags

    Returns:
    List of 1 or 0, where 1 if tag existed in the original
    taglist, and 0 else.
    """
    tc = tagcounter
    tc.update(taglist)
    return tc.values()


# Now, we map all of our patents into the higher dimensional space
lines = read_file('buzzx.csv')
lines.next() # skip CSV schema
patents = {}
i = 0
for line in lines:
    i += 1            
    if i % 10000 == 0:
        print 'finished ', i, str(datetime.now())
    number = line[0]
    tags = line[3].split(' ')
    patents[number] = expand_tagspace(tags)
tags = numpy.array(patents.values())
numbers = numpy.array(patents.keys())
scipy.io.savemat('patenttags', {'tags': tags, 'patents': numbers})
