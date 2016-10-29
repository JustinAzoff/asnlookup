from .backend import ASNLookup, FIELDS, ASRecord

import json
import logging
import time
import time
import zmq

logger = logging.getLogger(__name__)

fields_bytes = json.dumps(FIELDS).encode()

def main():
    logging.basicConfig(level=logging.DEBUG)
    context = zmq.Context()
    socket = context.socket(zmq.ROUTER)
    socket.bind("tcp://*:5555")

    logger.info("Initializing...")
    l = ASNLookup()
    logger.info("Startup complete")

    last_reload_check = time.time()

    while True:
        #  Wait for next request from client
        r = socket.recv_multipart()
        ident,  msg = r
        if msg == b"fields":
            socket.send_multipart([ident, fields_bytes])
            #TODO: better way to do this?
            if time.time() - last_reload_check > 300:
                l.reload_if_neaded()
                last_reload_check = time.time()
            continue

        ips = msg.decode().split()
        response = [l.lookup(ip) for ip in ips]
        #  Send reply back to client
        socket.send_multipart([ident, json.dumps(response).encode()])


if __name__ == "__main__":
    main()
