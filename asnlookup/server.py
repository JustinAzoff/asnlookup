from .backend import ASNLookup

import time
import zmq
import json

def main():
    context = zmq.Context()
    socket = context.socket(zmq.REP)
    socket.bind("tcp://*:5555")

    l = ASNLookup()

    while True:
        #  Wait for next request from client
        ip = socket.recv_string()
        response = l.lookup(ip)
        #  Send reply back to client
        socket.send_string(json.dumps(response._asdict()))

if __name__ == "__main__":
    main()
