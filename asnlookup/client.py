import json
import logging
import sys
import zmq

logger = logging.getLogger(__name__)

class ASNClient:
    def __init__(self, endpoint='tcp://localhost:5555'):
        context = zmq.Context()
        logger.debug("Connecting to asn lookup server")
        socket = context.socket(zmq.REQ)
        socket.connect(endpoint)
        self.socket = socket

    def lookup(self, ip):
        self.socket.send_string(ip)
        message = self.socket.recv_string()
        return json.loads(message)

def main():
    c = ASNClient()
    for line in sys.stdin:
        print(c.lookup(line.rstrip()))

if __name__ == "__main__":
    main()
