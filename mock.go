package main

import (
	"fmt"
	"io"
	"regexp"

	"gopkg.in/yaml.v2"
	"net/http"
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
}
type serviceEntry struct {
	Header map[string]string `yaml:"header"`
	Output serviceoutput     `yaml:"output"`
	Method string            `yaml:"method"`
	Path   string            `yaml:"path"`
	Name   string            `yaml:"name"`
	pathre *regexp.Regexp
}

type services []serviceEntry

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
		res[p].pathre = pathre
		if se.Output.Code == 0 {
			res[p].Output.Code = http.StatusOK
		}
	}
	return res, nil
}
