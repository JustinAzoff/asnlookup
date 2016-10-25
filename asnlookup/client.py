from . import asnlookup_pb2

import grpc

import logging
import sys

logger = logging.getLogger(__name__)

class ASNClient:
    def __init__(self, endpoint='localhost:5555'):
        channel = grpc.insecure_channel(endpoint)
        self.stub = asnlookup_pb2.AsnlookupStub(channel)

        logger.debug("Connecting to asn lookup server")
        logger.debug(self.stub.Hello(asnlookup_pb2.Empty()))

    def lookup_many(self, ips):
        requests = (asnlookup_pb2.LookupRequest(address=ip) for ip in ips)
        return self.stub.LookupMany(requests)

    def lookup(self, ip):
        return self.stub.Lookup(asnlookup_pb2.LookupRequest(address=ip))

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
