#!/usr/bin/env python

import logging
import sys
import argparse
import zmq
import os
from helpers import insecure_client, secure_client


def run(pull_address, pub_address, insecure):
    ctx = zmq.Context.instance()

    pub = ctx.socket(zmq.PUB)
    logging.debug("PUB bound to %s" % pub_address)
    pub.bind(pub_address)

    if insecure:
        logging.warn("Creating insecure PULL socket.")
        client = insecure_client
    else:
        logging.debug("Creating secure PULL socket.")
        client = secure_client

    with client(zmq.PULL) as pull:
        logging.debug("PULL connected to %s" % pull_address)
        pull.connect(pull_address)

        while(True):
            message = pull.recv()
            logging.debug("Message received: %s" % message)
            pub.send(message)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Listens to a PULL socket and publishes to a PUB socket.')
    parser.add_argument('--pull', dest='pull', help='Host to listen on (localhost, port 7000 by default)', default="tcp://127.0.0.1:7000")
    parser.add_argument('--pub', dest='pub', help='Address to publish on (localhost, port 7001 by default)', default="tcp://127.0.0.1:7001")
    parser.add_argument('--insecure', dest='insecure', help='Run without security (useful for testing)', action='store_const', const=True, default=False)
    parser.add_argument('-v', dest='verbose', help='Verbose mode', action="store_const", const=True)
    args = parser.parse_args()
    
    if zmq.zmq_version_info() < (4,0):
        raise RuntimeError("Security is not supported in libzmq version < 4.0. libzmq version {0}".format(zmq.zmq_version()))

    if args.verbose:
        level = logging.DEBUG
    else:
        level = logging.INFO

    logging.basicConfig(level=level, format="[%(levelname)s] %(message)s")
    run(args.pull, args.pub, args.insecure)
