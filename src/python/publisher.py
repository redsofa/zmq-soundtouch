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
from helpers import insecure_client, secure_client
from contextlib import contextmanager

def run(pull_address, pub_address, insecure):
    pub_socket = bind_pub_socket(pub_address)

    with connect_pull_socket(pull_address, insecure) as pull_socket:
        logging.debug("PULL connected to %s" % pull_address)
        forward(pull_socket, pub_socket)

def bind_pub_socket(address):
    ctx = zmq.Context.instance()
    socket = ctx.socket(zmq.PUB)
    socket.bind(address)
    logging.debug("PUB bound to %s" % address)
    return socket

@contextmanager
def connect_pull_socket(address, insecure):
    client_creator = choose_client_creator(insecure)
    with client_creator(zmq.PULL) as socket:
        socket.connect(address)
        yield socket

def choose_client_creator(insecure):
    if insecure:
        logging.warn("Using insecure PULL socket.")
        return insecure_client
    else:
        logging.debug("Using secure PULL socket.")
        return secure_client

def forward(in_socket, out_socket):
    while(True):
        message = in_socket.recv()
        logging.debug("Message received: %s" % message)
        out_socket.send(message)

def ensure_security_supported():
    if zmq.zmq_version_info() < (4,0):
        raise RuntimeError("Security is not supported in libzmq version < 4.0. libzmq version {0}".format(zmq.zmq_version()))

def log_level(verbose):
    if verbose:
        return logging.DEBUG
    else:
        return logging.INFO

def configure_logging(verbose):
    level = log_level(verbose)
    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")

def populate_parser_args(parser):
    parser.add_argument('--pull', dest='pull', help='Host to listen on (localhost, port 7000 by default)', default="tcp://127.0.0.1:7000")
    parser.add_argument('--pub', dest='pub', help='Address to publish on (localhost, port 7001 by default)', default="tcp://127.0.0.1:7001")
    parser.add_argument('--insecure', dest='insecure', help='Run without security (useful for testing)', action='store_const', const=True, default=False)
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)

def user_args():
    parser = argparse.ArgumentParser(description='Listens to a PULL socket and publishes to a PUB socket.')
    populate_parser_args(parser)
    args = parser.parse_args()
    return (args.verbose, args.pull, args.pub, args.insecure)

if __name__ == '__main__':
    (verbose, pull, pub, insecure) = user_args()
    
    ensure_security_supported()

    configure_logging(verbose)

    run(pull, pub, insecure)
