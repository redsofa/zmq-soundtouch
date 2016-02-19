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

package handlers

import (
	"github.com/redsofa/logger"
	"net/http"
	"time"
)

/*
 * Basic Wrapper for handlers which logs all requests.
 * Handlers get passed to this function which then
 * call the serveHTTP method on the handler and logs
 * the call
 * Source :
 * - http://thenewstack.io/make-a-restful-json-api-go/
 * - https://groups.google.com/forum/#!topic/golang-nuts/s7Xk1q0LSU0
 *
 * Usage Example :
 *
 * func main() {
 *	 router := mux.NewRouter()
 *
 *   router.Handle("/api/someresource", &handlers.SomeHandler{})
 *   http.ListenAndServe(":3000",  handlers.HttpLog(router))
 * }
 *
 * Notice the handlers.HttpLog(router) in the last line of the code above
 */

func HttpLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info.Printf(
			"Method : %s\t URI: %s\t URL: %s\t Time: %s",
			r.Method,
			r.RequestURI,
			r.URL,
			time.Since(start),
		)

		handler.ServeHTTP(w, r)

	})
}
