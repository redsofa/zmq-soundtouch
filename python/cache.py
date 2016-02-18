#!/usr/bin/env python

import logging
import os
import sys
import argparse
import zmq

def run(sub_address, rep_address, n=1):
    """
    @param n: How many last values to cache.
    """
    assert(n > 0)

    ctx = zmq.Context.instance()
    sub = ctx.socket(zmq.SUB)
    sub.setsockopt(zmq.SUBSCRIBE, "")
    logging.debug("Connecting SUB socket to %s" % sub_address)
    sub.connect(sub_address)

    rep = ctx.socket(zmq.REP)
    logging.debug("Binding REP socket to %s" % rep_address)
    rep.bind(rep_address)

    # store last N values
    cache = "None"

    poller = zmq.Poller()
    poller.register(sub, zmq.POLLIN)
    poller.register(rep, zmq.POLLIN)

    while True:
        try:
            events = dict(poller.poll(1000))
        except KeyboardInterrupt:
            print("Interrupted")
            break

        # If new data we cache
        if sub in events:
            msg = sub.recv()
            logging.debug("Received '%s'" % msg)
            cache = msg

        # If request - forward cache
        if rep in events:
            logging.debug("REQ received.")
            event = rep.recv()
            rep.send(cache)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Cache server.')
    parser.add_argument('--sub-address', dest='sub_address', help='Address to subscribe to (tcp://127.0.0.1:7001 by default)', default="tcp://127.0.0.1:7001")
    parser.add_argument('--rep-address', dest='rep_address', help='Address to handle REQ requests on (tcp://127.0.0.1:8000 by default)', default="tcp://127.0.01:8000")
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.sub_address, args.rep_address)
