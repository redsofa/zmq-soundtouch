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
    context = create_context()
    sub = connect_sub_socket(context, sub_address)
    router = bind_router_socket(context, router_address)
    poller = create_poller([sub, router])

    cache = create_cache(n)
    loop(poller, sub, router, cache)

def create_context():
    return zmq.Context.instance()

def connect_sub_socket(context, address):
    socket = context.socket(zmq.SUB)
    socket.setsockopt(zmq.SUBSCRIBE, "")
    logging.debug("Connecting SUB socket to %s" % address)
    socket.connect(address)
    return socket

def bind_router_socket(context, address):
    socket = context.socket(zmq.ROUTER)
    logging.debug("Binding ROUTER socket to %s" % address)
    socket.bind(address)
    return socket

def create_poller(sockets):
    poller = zmq.Poller()
    for socket in sockets:
        poller.register(socket, zmq.POLLIN)
    return poller

def create_cache(n):
    # store last N values (newest items in the front, or first)
    assert(n > 0)
    return collections.deque([], n)

def loop(poller, sub, router, cache):
    while True:
        poll(sub, router, poller, cache)

def poll(sub, router, poller, cache):
    events = poll_with_interrupt(poller)

    if check_for_sub(sub, events):
        handle_sub(sub, events, cache)

    if check_for_router(router, events):
        handle_router(router, events, cache)
                
def poll_with_interrupt(poller):
    try:
        return dict(poller.poll(1000))
    except KeyboardInterrupt:
        print("Interrupted")
        sys.exit(-1)

def handle_sub(sub, events, cache):
    msg = sub.recv()
    logging.debug("Received '%s'" % msg)
    cache.appendleft(msg)

def check_for_sub(sub, events):
    return sub in events

def check_for_router(router, events):
    return router in events

def handle_router(router, events, cache):
    # Forward cached items
    logging.debug("Request received.")
    ident, msg = router.recv_multipart()

    if not icanhaz(msg):
        logging.warn("Invalid request: '%s'." % msg)
        sys.exit(-1)
    else:
        send_cached_items(ident, cache, router)

def send_cache(ident, cache, router):
    for item in cache:
        router.send_multipart([ident, item])

def send_bye(ident, router):
    router.send_multipart([ident,'KTHXBYE'])

def icanhaz(msg):
    return msg == "ICANHAZ?"

def send_cached_items(ident, cache, router):
    send_cache(ident, cache, router)
    send_bye(ident, router)

def populate_parser_args(parser):
    parser.add_argument('--sub-address', dest='sub_address', help='Address to subscribe to (tcp://127.0.0.1:7001 by default)', default="tcp://127.0.0.1:7001")
    parser.add_argument('--router-address', dest='router_address', help='Address to handle REQ requests on (tcp://127.0.0.1:8000 by default)', default="tcp://127.0.01:8000")
    parser.add_argument('-n', dest='n', help='Number of items to cache (must be greater than zero). Defaults to 10', default=10)
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    
def user_args():
    parser = argparse.ArgumentParser(description='Cache server.')
    populate_parser_args(parser)
    args = parser.parse_args()
    return (args.verbose, args.sub_address, args.router_address, args.n)

def log_level(verbose):
    if verbose:
        return logging.DEBUG
    else:
        return logging.INFO

def configure_logging(verbose):
    level = log_level(verbose)
    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")

if __name__ == '__main__':
    (verbose, sub_address, router_address, n) = user_args()
    
    configure_logging(verbose)

    run(sub_address, router_address, n)
