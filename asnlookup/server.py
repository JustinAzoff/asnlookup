from .backend import ASNLookup, FIELDS, ASRecord
from . import asnlookup_pb2

from concurrent import futures
import grpc
import logging
import time

logger = logging.getLogger(__name__)

class AsnlookupServicer(asnlookup_pb2.AsnlookupServicer):
    def __init__(self):
        logger.info("Initializing...")
        self.l = ASNLookup()
        logger.info("Startup complete")
        self.last_reload_check = time.time()

    def _get_response(self, address):
        response = self.l.lookup(address)
        return asnlookup_pb2.LookupReply(
            address=response.ip,
            asn=response.asn,
            prefix=response.prefix,
            owner=response.owner,
            cc=response.cc,
        )

    def Hello(self, request, context):
        if time.time() - self.last_reload_check > 300:
            self.l.reload_if_neaded()
            self.last_reload_check = time.time()
        return asnlookup_pb2.HelloReply(message="Hello!")

    def Lookup(self, request, context):
        return self._get_response(request.address)

    def LookupMany(self, request_iterator, context):
        for request in request_iterator:
            yield self._get_response(request.address)

def main():
    logging.basicConfig(level=logging.DEBUG)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    asnlookup_pb2.add_AsnlookupServicer_to_server(AsnlookupServicer(), server)
    server.add_insecure_port('[::]:5555')
    server.start()
    try:
        while True:
            time.sleep(1000)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    main()
