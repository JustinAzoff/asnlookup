from .backend import ASNLookup, FIELDS, ASRecord

import json
import logging
import time
import time
import zmq

logger = logging.getLogger(__name__)

def main():
    logging.basicConfig(level=logging.DEBUG)
    context = zmq.Context()
    socket = context.socket(zmq.REP)
    socket.bind("tcp://*:5555")

    logger.info("Initializing...")
    l = ASNLookup()
    logger.info("Startup complete")

    last_reload_check = time.time()

    while True:
        #  Wait for next request from client
        msg = socket.recv_string()
        if msg == "fields":
            socket.send_string(json.dumps(FIELDS))
            continue

        ips = msg.split()
        response = [l.lookup(ip) for ip in ips]
        #  Send reply back to client
        socket.send_string(json.dumps(response))


        #TODO: better way to do this?
        if time.time() - last_reload_check > 300:
            l.reload_if_neaded()
            last_reload_check = time.time()

if __name__ == "__main__":
    main()
