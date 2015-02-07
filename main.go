package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/braintree/manners"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/lestrrat/go-server-starter-listener"

	"github.com/mix3/plantumlor/controller"
	"github.com/mix3/plantumlor/middleware"
	"github.com/mix3/plantumlor/plantuml"
)

var (
	host         string
	port         uint
	plantumlPath string
)

func init() {
	const (
		defaultHost = ""
		defaultPort = 8080
		defaultPath = ""
	)
	flag.StringVar(&host, "host", defaultHost, "host")
	flag.StringVar(&host, "h", defaultHost, "host")
	flag.UintVar(&port, "port", defaultPort, "port")
	flag.UintVar(&port, "p", defaultPort, "port")
	flag.StringVar(&plantumlPath, "plantumlPath", defaultPath, "plantuml command path")
}

func main() {
	flag.Parse()

	log.Printf("host: %v", host)
	log.Printf("port: %v", port)
	log.Printf("path: %v", plantumlPath)

	plantUML, err := plantuml.NewPlantUML(plantumlPath)
	if err != nil {
		log.Fatal(err)
	}

	c := &controller.AppContext{plantUML}
	commonHandlers := alice.New(
		middleware.LoggerMiddleware,
		middleware.RecoverMiddleware,
	)
	router := httprouter.New()
	router.GET("/transfer/", wrapHandler(commonHandlers.ThenFunc(c.TransferHandler)))
	router.GET("/transfer/:data", wrapHandler(commonHandlers.ThenFunc(c.TransferHandler)))
	router.NotFound = commonHandlers.Then(
		http.FileServer(
			&assetfs.AssetFS{
				Asset:    Asset,
				AssetDir: AssetDir,
				Prefix:   "",
			},
		),
	).ServeHTTP
	serve(router)
}

func wrapHandler(h http.Handler) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
	return fn
}

func serve(mux http.Handler) {
	l, _ := ss.NewListener()
	if l == nil {
		var err error
		l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			log.Fatalf("Failed to listen to port %d", port)
		}
	}

	s := manners.NewServer()
	s.Serve(manners.NewListener(l, s), mux)
}
