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
import sys
import argparse
import zmq
import os

def run(dealer_address):
    ctx = zmq.Context.instance()

    dealer = ctx.socket(zmq.DEALER)
    logging.debug("Connecting DEALER socket to %s" % dealer_address)
    dealer.connect(dealer_address)

    dealer.send("ICANHAZ?")
    while True:
        msg = dealer.recv()
        if msg != 'KTHXBYE':
            logging.info("Received '%s'" % msg)
        else:
            logging.info("End of transmission: %s" % msg)
            break



if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Makes a test request on DEALER socket.')
    parser.add_argument('--dealer', dest='dealer', help='Host to make request on (localhost, port 8000 by default)', default="tcp://127.0.0.1:8000")
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if zmq.zmq_version_info() < (4,0):
        raise RuntimeError("Security is not supported in libzmq version < 4.0. libzmq version {0}".format(zmq.zmq_version()))

    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.dealer)
