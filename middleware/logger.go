package middleware

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var ApacheFormatPattern = "%s - %s [%s] \"%s\" %d %d %d\n"

// https://gist.github.com/cespare/3985516
type ApacheLogRecord struct {
	http.ResponseWriter

	host                  string
	username              string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
}

func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)

	fmt.Fprintf(out,
		ApacheFormatPattern,
		r.host,
		r.username,
		timeFormatted,
		requestLine,
		r.status,
		r.responseBytes,
		int64(r.elapsedTime.Nanoseconds()/1000),
	)
}

func (r *ApacheLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

func (r *ApacheLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// https://github.com/gorilla/handlers/blob/master/handlers.go#L217
		url := *r.URL
		username := "-"
		if url.User != nil {
			if name := url.User.Username(); name != "" {
				username = name
			}
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		record := &ApacheLogRecord{
			ResponseWriter: w,
			host:           host,
			username:       username,
			time:           time.Time{},
			method:         r.Method,
			uri:            r.RequestURI,
			protocol:       r.Proto,
			status:         http.StatusOK,
			elapsedTime:    time.Duration(0),
		}

		startTime := time.Now()
		next.ServeHTTP(record, r)
		finishTime := time.Now()

		record.time = finishTime.UTC()
		record.elapsedTime = finishTime.Sub(startTime)

		record.Log(os.Stderr)
	}
	return http.HandlerFunc(fn)
}
