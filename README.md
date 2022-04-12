zmq-soundtouch
---------------
A series of projects exploring secure ZeroMQ connections and SoundTouch event data collection. Result  is the live SoundTouch event data visualization seen below. 


![Screencast](./docs/demo.gif "Screencast")




Main Project List : 
-------------------

1) - ./go/src/github.com/redsofa/`soundtouch` - Project that connects to Bose SoundTouch over WebSocket and pushes event notification messages to a secure ZeroMQ TCP PULL socket

2) - ./python/`publisher`.py - Project that connects to secure TCP PUSH socket, receives SoundTouch notifications and broadcasts them to subscribers.

3) - ./python/`cache`.py - Project that connects to ZeroMQ TCP PUB socket, creates a cache containing a list of recent SoundTouch notification messages	and makes them available over a ZeroMQ Router socket.

4) - ./go/src/github.com/redsofa/`collector` - Project that connects to ZeroMQ TCP PUB socket, receives notification messages and passes them on to WebSocket clients. Project also connects to ZeroMQ Router socket to get a list of most recent messages. In addition, it serves static web content.

Install
-------
Project contains a Vagrantfile with centos7. Before proceeding install Vagrant on your local machine: https://www.vagrantup.com/downloads.html

You will now create a virtual machine and us it develop and run the code. To install all the dependencies and get the machine ready, do this:

1. `git clone git@github.com:redsofa/zmq-soundtouch.git`
2. Navigate to root directory of the project: `cd zmq-soundtouch`.
3. Create a vagrant vm: `vagrant up`
4. Ssh into vm: `vagrant ssh`.

You can find the `src` directory shared in the vagrant vm under the `/usr/local/src` path. You can find the root directory itself shared under `/vagrant'.



Project Topology :
-------------------
![Topology](./docs/topology.png "Topology")


License
-------
This project is licensed under AGPL V3. See COPYING file for a copy of the included license or visit http://www.gnu.org/licenses/.
