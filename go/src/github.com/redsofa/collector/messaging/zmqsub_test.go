/*
Copyright 2016 Rene Richard

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
*/

package messaging

import (
	"github.com/redsofa/logger"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

//TODO : Write actual/useful test
func TestZmqSub(t *testing.T) {
	logger.Info.Println("***************************************************")
	zmqSub := NewZmqSub()
	zmqSub.Start(2) //Start, but quit after 2 seconds
	logger.Info.Println("***************************************************")
}
