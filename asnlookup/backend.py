#!/usr/bin/env python
from .data_manager import update_asnnames, update_asndb

from collections import namedtuple
import json
import logging
import os
import pyasn
import time

logger = logging.getLogger(__name__)

def load_asnames(fn):
    data = {}
    with open(fn) as f:
        upstream_data = json.load(f)

    for k, v in upstream_data.items():
        if ',' not in v:
            owner = cc = ''
            if len(v) == 2:
                cc = v
            else:
                owner = v
        else:
            owner, cc = v.rsplit(",", 1)
            data[k] = (owner.strip(), cc.strip())
    return data

FIELDS = "ip", "asn", "prefix", "owner", "cc"
ASRecord = namedtuple("ASRecord", FIELDS)

class ASNLookup(object):
    namedb_filename = 'asnames.json'
    db_filename = 'asn.db'

    def __init__(self):
        update_asndb(self.db_filename, 24)
        update_asnnames(self.namedb_filename, 24)

        self.reload()

    def reload(self):
        start = time.time()
        logger.debug("reloading databases...")
        self.asndb = pyasn.pyasn(self.db_filename)
        self.asnames = load_asnames(self.namedb_filename)

        self.db_ino = os.stat(self.db_filename).st_ino
        self.namedb_ino = os.stat(self.namedb_filename).st_ino
        end = time.time()
        logger.debug("reloading databases complete seconds=%0.1f", end-start)

    def reload_neaded(self):
        db_ino = os.stat(self.db_filename).st_ino
        namedb_ino = os.stat(self.namedb_filename).st_ino

        if db_ino != self.db_ino or namedb_ino != self.namedb_ino:
            self.reload()

    def reload_if_neaded(self):
        if self.reload_neaded():
            self.reload()

    def lookup_asname(self, asn):
        rec = self.asnames.get(str(asn))
        if not rec:
            return "NA", "NA"
        return rec

    def lookup(self, ip):
        try:
            rec =  self.asndb.lookup(ip)
        except:#FIXME
            logger.exception("Lookup failed for ip=%s", ip)
            return ASRecord(ip, 'NA', 'NA', 'NA', 'NA')
        asn, prefix = rec
        asn = asn if asn else 'NA'
        prefix = prefix if prefix else 'NA'

        owner, cc = self.lookup_asname(asn)
        return ASRecord(ip, asn, prefix, owner, cc)
