#!/usr/bin/env python

import logging
import sys
import argparse
import zmq
import os
import time

def run(push_address):
    ctx = zmq.Context.instance()

    push = ctx.socket(zmq.PUSH)
    logging.debug("Binding PUSH socket to %s" % push_address)
    push.bind(push_address)

    c = 0
    while True:
        logging.debug("Sending message...")
        push.send("hey %s" % c)
        c += 1
        c = c % 10
        time.sleep(1)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Makes a test request on REQ socket.')
    parser.add_argument('--push', dest='push', help='Address to push to (localhost, port 7000 by default)', default="tcp://127.0.0.1:7000")
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if zmq.zmq_version_info() < (4,0):
        raise RuntimeError("Security is not supported in libzmq version < 4.0. libzmq version {0}".format(zmq.zmq_version()))

    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.push)
