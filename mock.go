package main

import (
	"fmt"
	"io"
	"regexp"

	"encoding/json"
	"net/http"
	"text/template"

	"strings"

	"github.com/oliveagle/jsonpath"
	"gopkg.in/xmlpath.v2"
	"gopkg.in/yaml.v2"
)

/*
-   path: "/user":
    method: GET
	header:
		xxx: 123
	output:
		contentType: "application/json"
		response: "{'name':'max'}"
*/

type serviceoutput struct {
	ContentType string `yaml:"contentType"`
	Response    string `yaml:"response"`
	Code        int    `yaml:"code"`
	response    *template.Template
}
type serviceEntry struct {
	Header   map[string]string `yaml:"header"`
	Output   serviceoutput     `yaml:"output"`
	Method   string            `yaml:"method"`
	Path     string            `yaml:"path"`
	Name     string            `yaml:"name"`
	pathre   *regexp.Regexp
	pathvars []string
}

type services []serviceEntry

var (
	funcmap = template.FuncMap{
		"jsonpath": jsonpt,
		"xmlpath":  xmlpt,
	}
)

func readServices(in io.Reader) (services, error) {
	var res services
	err := yaml.NewDecoder(in).Decode(&res)
	if err != nil {
		return nil, err
	}
	for p, se := range res {
		pathre, err := regexp.Compile(se.Path)
		if err != nil {
			return nil, fmt.Errorf("cannot compile %q as a regeexp: %v", se.Path, err)
		}
		subexps := pathre.SubexpNames()
		if len(subexps) > 1 {
			res[p].pathvars = subexps[1:]
		}
		res[p].pathre = pathre
		if se.Output.Code == 0 {
			res[p].Output.Code = http.StatusOK
		}
		if se.Output.Response != "" {
			t, err := template.New(se.Name).Funcs(funcmap).Parse(se.Output.Response)
			if err != nil {
				return nil, fmt.Errorf("cannot compile %s as a template: %v", se.Output.Response, err)
			}
			res[p].Output.response = t
		}
	}
	return res, nil
}

// it is inefficient to parse the data on every invocation of the following functions
// but hey: this is a mocking-server, it should not be used for high performance
// data processing :-)

func jsonpt(data, key string) string {
	datamap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &datamap)
	if err != nil {
		return fmt.Sprintf("%s cannot be parsed as json", data)
	}
	res, err := jsonpath.JsonPathLookup(datamap, key)
	if err != nil {
		return fmt.Sprintf("%q not found in data: %v", key, err)
	}
	return fmt.Sprintf("%s", res)
}

func xmlpt(data, key string) string {
	path, err := xmlpath.Compile(key)
	if err != nil {
		return fmt.Sprintf("cannot compile the xpath %s: %v", key, err)
	}
	root, err := xmlpath.Parse(strings.NewReader(data))
	if err != nil {
		return fmt.Sprintf("cannot parse the input as xml: %v", err)
	}
	if value, ok := path.String(root); ok {
		return value
	}
	return fmt.Sprintf("%s not found", key)
}
