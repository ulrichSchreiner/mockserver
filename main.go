package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"net/http/httputil"

	"io/ioutil"

	"gopkg.in/yaml.v2"
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
		if rq.Header.Get("Content-Type") == "application/x-www-form-urlencoded" && rq.Method == http.MethodGet {
			rq.ParseForm()
		}
		templatedata := rqdata(rq)
		for _, srv := range srvs {
			match, pathvars := matchService(srv, rq, templatedata["BODY"].(string))
			if match {
				templatedata["PATH"] = pathvars
				if srv.Output.ContentType != "" {
					rsp.Header().Add("Content-Type", srv.Output.ContentType)
				}
				rsp.WriteHeader(srv.Output.Code)
				var buf bytes.Buffer
				err := srv.Output.response.Execute(&buf, templatedata)
				if err != nil {
					log.Printf("[ERROR] cannot render template: %v", err)
				}
				if *debug {
					log.Printf("RESPONSE: %s\n", buf.String())
				}
				rsp.Write(buf.Bytes())
				log.Printf("[INFO] use handler '%s %s'", srv.Method, srv.Name)
				return
			}
		}
		rsp.WriteHeader(http.StatusNotFound)
	})
}

func rqdata(rq *http.Request) map[string]interface{} {
	rqparams := make(map[string]string)
	headers := make(map[string]string)

	for k := range rq.Header {
		headers[k] = rq.Header.Get(k)
	}
	for k := range rq.Form {
		rqparams[k] = rq.FormValue(k)
	}
	body, _ := ioutil.ReadAll(rq.Body)

	return map[string]interface{}{
		"RQ":     rqparams,
		"HEADER": headers,
		"BODY":   string(body),
	}
}

func matchService(s serviceEntry, rq *http.Request, body string) (bool, map[string]string) {
	if rq.Method != s.Method {
		return false, nil
	}
	for k, v := range s.Header {
		rqval := rq.Header.Get(k)
		if rqval != v {
			return false, nil
		}
	}
	frm := rq.Form
	qry := rq.URL.Query()
	for k,v := range s.Params {
		pval := frm.Get(k)
		if pval == "" {
			pval = qry.Get(k)
		}
		if pval != v {
			fmt.Printf("wrong pval: %q != %q\n",pval, v)
			return false, nil
		}
	}
	m := s.pathre.MatchString(rq.URL.Path)
	if s.BodyMatch != "" {
		m = strings.Contains(body, s.BodyMatch)
	}
	pathvars := make(map[string]string)
	if m {
		if len(s.pathvars) > 0 {
			subm := s.pathre.FindStringSubmatch(rq.URL.Path)
			for i, p := range s.pathvars {
				subidx := i + 1
				if len(subm) > subidx {
					pathvars[p] = subm[subidx]
				}
			}
		}
	}
	return m, pathvars
}
