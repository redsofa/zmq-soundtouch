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

from os.path import dirname
import shutil
import zmq
from zmq.auth import create_certificates

def create_keys(indir):
    """
    Generate client and server keys
    """
    create_certificates(indir, "client")
    create_certificates(indir, "server")

if __name__ == '__main__':
	version_info = zmq.zmq_version_info()
	if version_info < (4,0):
		raise ValueError("ZMQ version < 4 does not support curve crypto. Current version: %s" % zmq.zmq_version())

	create_keys(dirname(__file__))
