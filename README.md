Mock your external REST API's with a local mock server.

You have to write a short `yaml` file and start this server. It will respond to the
given endpoints with the specified content.

Example:
~~~yaml
- path: "^/user2$"
  method: GET
  output:
    contentType: "text/xml"
    response: |
      <mydata>
        <name>max</name>
      </mydata>

- path: "^/user$"
  method: GET
  output:
    contentType: "application/json"
    response: "{'name':'max'}"
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
$ curl http://localhost:9099/user
{'name':'max'}
~~~

