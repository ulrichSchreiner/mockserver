package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
	"net/http/httputil"
)

var (
	listen = flag.String("address", "localhost:9099", "sets the listen port")
	debug  = flag.Bool("debug", false, "use request logging")
)

func main() {
	flag.Parse()

	input := "mocks.yaml"
	if flag.Arg(0) != "" {
		input = flag.Arg(0)
	}

	infile, err := os.Open(input)
	if err != nil {
		log.Fatalf("cannot open file %q: %v", input, err)
	}
	defer infile.Close()

	mocks, err := readServices(infile)
	if err != nil {
		log.Fatalf("cannot parse the file %q: %v", input, err)
	}
	fmt.Println("using mocks ...")
	yaml.NewEncoder(os.Stdout).Encode(mocks)
	fmt.Printf("start service %s ...\n", *listen)
	log.Fatalf("cannot start server: %v", http.ListenAndServe(*listen, handleMocks(mocks)))
}

func handleMocks(srvs services) http.Handler {
	return http.HandlerFunc(func(rsp http.ResponseWriter, rq *http.Request) {
		if *debug {
			out, _ := httputil.DumpRequest(rq, true)
			log.Printf("REQUEST: %s\n", string(out))
		}
		for _, srv := range srvs {
			if matchService(srv, rq) {
				if srv.Output.ContentType != "" {
					rsp.Header().Add("Content-Type", srv.Output.ContentType)
				}
				rsp.WriteHeader(srv.Output.Code)
				fmt.Fprint(rsp, srv.Output.Response)
				log.Printf("[INFO] use handler '%s %s'", srv.Method, srv.Name)
				return
			}
		}
		rsp.WriteHeader(http.StatusNotFound)
	})
}

func matchService(s serviceEntry, rq *http.Request) bool {
	if rq.Method != s.Method {
		return false
	}
	for k, v := range s.Header {
		rqval := rq.Header.Get(k)
		if rqval != v {
			return false
		}
	}

	return s.pathre.MatchString(rq.RequestURI)
}
