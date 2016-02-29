Building The Go Project Artifacts
---------------------------------
- Run once to get go dependencies :
	- `make go-deps`

- Run to build all artifacts (ends up in a ./dist directory)
	- `make`


Notes About ZMQ and Curve :
--------------------------
In order to have curve support ensure that you have libsodium installed and configured.

- The Vagrantfile in the project's root automates the installation of libsodium and zmqlib. 


Additional Information : 
------------------------
- On Linux :
	- http://mythinkpond.com/2015/09/06/how-to-install-and-configure-zeromq-libsodium-on-centos-6-7/

- On Mac : 
	- See solved issue here : https://github.com/pebbe/zmq4/issues/72