from . import asnlookup_pb2

import grpc

import logging
import sys

logger = logging.getLogger(__name__)

def chunk(it, slice=50):
    """Generate sublists from an iterator
    >>> list(chunk(iter(range(10)),11))
    [[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]]
    >>> list(chunk(iter(range(10)),9))
    [[0, 1, 2, 3, 4, 5, 6, 7, 8], [9]]
    >>> list(chunk(iter(range(10)),5))
    [[0, 1, 2, 3, 4], [5, 6, 7, 8, 9]]
    >>> list(chunk(iter(range(10)),3))
    [[0, 1, 2], [3, 4, 5], [6, 7, 8], [9]]
    >>> list(chunk(iter(range(10)),1))
    [[0], [1], [2], [3], [4], [5], [6], [7], [8], [9]]
    """

    assert(slice > 0)
    a=[]

    for x in it:
        if len(a) >= slice :
            yield a
            a=[]
        a.append(x)

    if a:
        yield a

class ASNClient:
    def __init__(self, endpoint='localhost:5555'):
        channel = grpc.insecure_channel(endpoint)
        self.stub = asnlookup_pb2.AsnlookupStub(channel)

        logger.debug("Connecting to asn lookup server")
        logger.debug(self.stub.Hello(asnlookup_pb2.Empty()))

    def lookup(self, ip):
        return self.stub.Lookup(asnlookup_pb2.LookupRequest(address=ip))

    def lookup_many(self, ips):
        for batch in chunk(ips, 100):
            requests = [asnlookup_pb2.LookupRequest(address=ip) for ip in batch]
            req = asnlookup_pb2.LookupRequestBatch(requests=requests)
            response = self.stub.LookupBatch(req)
            yield from response.replies

def main():
    logging.basicConfig(level=logging.DEBUG)

    endpoint = 'localhost:5555'
    if len(sys.argv) > 1:
        endpoint = sys.argv[1]
    c = ASNClient(endpoint)
    ips = (line.rstrip() for line in sys.stdin)
    for rec in c.lookup_many(ips):
        print("\t".join(str(f) for f in (rec.address, rec.asn, rec.prefix, rec.owner, rec.cc)))

if __name__ == "__main__":
    main()
