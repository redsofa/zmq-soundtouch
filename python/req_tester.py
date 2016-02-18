#!/usr/bin/env python

import logging
import sys
import argparse
import zmq
import os

def run(req_address):
    ctx = zmq.Context.instance()

    req = ctx.socket(zmq.REQ)
    logging.debug("Connecting REQ socket to %s" % req_address)
    req.connect(req_address)

    req.send("HELLO")
    msg = req.recv()
    logging.info("Received '%s'" % msg)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Makes a test request on REQ socket.')
    parser.add_argument('--req', dest='req', help='Host to listen on (localhost, port 8000 by default)', default="tcp://127.0.0.1:8000")
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if zmq.zmq_version_info() < (4,0):
        raise RuntimeError("Security is not supported in libzmq version < 4.0. libzmq version {0}".format(zmq.zmq_version()))

    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.req)
