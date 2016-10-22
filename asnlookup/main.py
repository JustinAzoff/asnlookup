#!/usr/bin/env python

from .backend import ASNLookup
import sys

def main():
    import logging
    logging.basicConfig(level=logging.DEBUG)
    l = ASNLookup()
    for line in sys.stdin:
        ip = line.strip()
        rec = l.lookup(ip)
        print("\t".join(str(s) for s in rec))


if __name__ == '__main__':
    main()
