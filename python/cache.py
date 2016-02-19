#!/usr/bin/env python

"""
Copyright 2016 Andriy Drozdyuk

This file is part of zmq-soundtouch.

zmq-soundtouch is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

zmq-soundtouch is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with zmq-soundtouch.  If not, see <http://www.gnu.org/licenses/>.
"""

import logging
import os
import sys
import argparse
import zmq
import collections 
def run(sub_address, router_address, n):
    """
    @param n: How many last values to cache.
    """
    assert(n > 0)

    ctx = zmq.Context.instance()
    sub = ctx.socket(zmq.SUB)
    sub.setsockopt(zmq.SUBSCRIBE, "")
    logging.debug("Connecting SUB socket to %s" % sub_address)
    sub.connect(sub_address)

    router = ctx.socket(zmq.ROUTER)
    logging.debug("Binding ROUTER socket to %s" % router_address)
    router.bind(router_address)

    # store last N values (newest items in the front, or first)
    cache = collections.deque([], n)

    poller = zmq.Poller()
    poller.register(sub, zmq.POLLIN)
    poller.register(router, zmq.POLLIN)

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
            cache.appendleft(msg)

        # If request - forward cached items
        if router in events:
            logging.debug("Request received.")
            ident, msg = router.recv_multipart()
            if msg != "ICANHAZ?":
                logging.warn("Invalid request: '%s'." % msg)
                break

            for item in cache:
                router.send_multipart([ident, item])

            logging.info("Sending bye.")
            router.send_multipart([ident,'KTHXBYE'])

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Cache server.')
    parser.add_argument('--sub-address', dest='sub_address', help='Address to subscribe to (tcp://127.0.0.1:7001 by default)', default="tcp://127.0.0.1:7001")
    parser.add_argument('--router-address', dest='router_address', help='Address to handle REQ requests on (tcp://127.0.0.1:8000 by default)', default="tcp://127.0.01:8000")
    parser.add_argument('-n', dest='n', help='Number of items to cache (must be greater than zero). Defaults to 10', default=10)
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.sub_address, args.router_address, args.n)
