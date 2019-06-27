[![Build Status](https://github.drone.protegear.io/api/badges/ulrichSchreiner/mockserver/status.svg)](https://github.drone.protegear.io/ulrichSchreiner/mockserver)

Mock your external REST API's with a local mock server.

You have to write a short `yaml` file and start this server. It will respond to the
given endpoints with the specified content. You can use static ouput or render header,
request or path variables from the request into the output.

Example:
~~~yaml
- name: the user2 endpoint
  path: "^/user2$"
  method: GET
  output:
    contentType: "text/xml"
    response: |
      <mydata>
        <name>max</name>
      </mydata>

- name: get the user endpoint
  path: "^/user/(?P<id>.*)$"
  method: GET
  output:
    contentType: "application/json"
    response: "{'name':'{{ .RQ.name }}', 'id':'{{ .PATH.id }}'}"
~~~

Now start the mockserver:
~~~shell
$ mockserver testfiles/usermocks.yaml
using mocks ...
- header: {}
  output:
    contentType: text/xml
    response: |
      <mydata>
        <name>max</name>
      </mydata>
  method: GET
  path: ^/user2$
- header: {}
  output:
    contentType: application/json
    response: '{''name'':''max''}'
  method: GET
  path: ^/user$
start service localhost:9099 ...
~~~

And now test your mocks:
~~~shell
$ curl http://localhost:9099/user23
$ curl http://localhost:9099/user2
<mydata>
  <name>max</name>
</mydata>
$ curl http://localhost:9099/user/123?name=max
{'name':'max', 'id':'123'}
~~~

Please make sure, your entries in the yaml file have the correct order. They are processed in order and
the first match will win.

You can also take data from posted data (json or xml) and use the data in the output:
~~~yaml
- name: post a new user
  path: "^/user$"
  method: POST
  output:
    contentType: "application/json"
    response: |
      {{ jsonpath .BODY "$.name"}}

- name: post a new user as xml
  path: "^/user/xml$"
  method: POST
  output:
    contentType: "text/xml"
    response: |
      {{ xmlpath .BODY "/user/name"}}
~~~
