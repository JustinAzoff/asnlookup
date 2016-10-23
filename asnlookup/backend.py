#!/usr/bin/env python
from .data_manager import update_asnnames, update_asndb

from collections import namedtuple
import logging
import pyasn
import json

logger = logging.getLogger(__name__)

def load_asnames(fn):
    with open(fn) as f:
        return json.load(f)

FIELDS = "ip", "asn", "prefix", "owner"
ASRecord = namedtuple("ASRecord", FIELDS)

class ASNLookup(object):
    namedb_filename = 'asnames.json'
    db_filename = 'asn.db'

    def __init__(self):
        update_asndb(self.db_filename, 24)
        update_asnnames(self.namedb_filename, 24)

        self.asndb = pyasn.pyasn(self.db_filename)
        self.asnames = load_asnames(self.namedb_filename)

    def lookup_asname(self, asn):
        return self.asnames.get(str(asn), "NA")

    def lookup(self, ip):
        try:
            rec =  self.asndb.lookup(ip)
        except:#FIXME
            logger.exception("Lookup failed for ip=%s", ip)
            rec = None
        if not rec:
            return ASRecord(ip, 'NA', 'NA', 'NA')
        asn, prefix = rec
        asn = asn if asn else 'NA'
        prefix = prefix if prefix else 'NA'

        owner = self.lookup_asname(asn)
        return ASRecord(ip, asn, prefix, owner)
