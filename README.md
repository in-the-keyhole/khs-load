# Keyhole Software Load Testing Utility

Command Line utility written in [Go](https://go.dev) that load tests API's with a specified simulated users. Reports throughput statistics and performance Graphs can be generated. 

## Running 

1. Clone Repository 

2. Open a command line in cloned directory and execute the following

```
    ./bin/macosx/khsload do https://keyholesoftware.com --users 2 --duration 30 
```
*Note: currently only Macosx build is avaiable, follow the __Installation and Running from Source__ section below to create an executable binary for other environemnts* 


3. Results will be output to the console 

## Command Line Flags 

The following command line flags can be used to configure your load test. These options can also be defined in a `YAML` config file, see section below.

```
    --duration      Int      Number of seconds to run test
    --users         Int      Number of users to simulate 
    --ramp          Int      Number of seconds between starting users
    --wait          Int      Number of seconds to wait between requests  
    --config        string   YAML config file see YAML Config section below
    --save          string   Save results to CSV file
    --replace       string   Save and replace file if exists results to CSV file
    --contenttype   string   Type (i.e. application/json) for POSTING data
    --authtoken     string   Authorization token 
    --tokentemplate string   Expression to format authtoken in request header
```
## Installing and Running from Source Code

1. [Install Go](https://go.dev/doc/install) 

2. Clone Repo 

3. Open a command line terminal and navigate to the repo directory and enter the following commands 

```
    go install
    go build  
```

This will create an executable named `khsload` in your directory. that can be executed with this command. See previous sections for options.

```
    ./khsload 
```

## Plotting Results 

Results saved to a `CSV` file can be plotted to a scatter graph with the following command 

**Run a Load test and save to `test.csv` with the following command**

```
    ./khsload do http://keyholesoftware.com --users 4 --duration 20 --save test.csv 
```

**Generate a plot with the resulting `test.csv`**

```
    ./khsload plot test.csv
```

A scatter based Graph will be plotted to a file named `khsplot.png` 

Here's and example graph

![](khsplot.png)

### Configuration YAML 

Instead of a command line flags test options can be defined in a `YAML` based configuration file. 

Here's an example YAML file...
```
#
# Number of Users to Simulate 
#
users: 25
#
# Seconds to wait while ramping us users
#
ramp: 2
#
# Seconds to run the API load test
#
duration: 120
#
# Secconds to wait inbetween API requests
#
wait: 1
#
# Template used to apply token to API request Headers
#
tokentemplate: "Bearer {{.}}"
#
# URL required to obtain an authorization Token
#
# 
auth:
  url: https://<authenticate URL>
  userid: xxxxx
  password: xxxxx
  tokenizeusing: ","
  gettoken: "token"
  splitwith: ":"
#
# URL's to load test
#
url:
  - POST,<URL>%201,key=value&key2=value2
  - GET,<URL>

```
## POST Requests 
GET requests are made by default, however POST requests can be made by prefixing the URL with `POST` as shown below.

```
    ./khsload do "post,http://<address>,key1=value&key2=value"
```

Key/Value data is supplied after the address. Content type is of `POSTED` data defaults to `application/json`. This can be changed to `application/x-www-form-urlencoded` using the following command line flag.


```
    ./khsload do "post,https://<Your POST Address Here>,akey1=avalue&akey2=avalue" -contenttype application/x-www-form-urlencoded
```


## Token Based Authentication 
This utilitty supports load testing `TOKEN` based authentication schemes. If an API has a persistent access token that can be applied to request headers,you can specify this using the command line `-token` flag, or define in the `YAML` config. 

```
    ./khsload do <some url> -token <auth token>
```

Tokens will be applied to request headers using the `tokentemplate` expression. This expression allows the token value to be applied to `authorization` request `Header`
field. 

``` 
    ./khsload do <some url> -authtoken <auth token> -tokentemplate "{{Bearer .}}"
```










