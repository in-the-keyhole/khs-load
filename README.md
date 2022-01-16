# Keyhole Software Load Testing Utility

This command Line utility, written in [Go](https://go.dev),  load tests APIs against specified simulated users. 
It can generate throughput statistics and performance Graphs. 

## Running 

1. Clone Repository 

2. Open a command line in cloned directory and execute one of the following:

###### macOs Intel

```
    ./bin/macOs/khsload do https://keyholesoftware.com --users 2 --duration 30 
```

*Note: we currently supply only macOs Intel pre-built targets in this repository. 
Follow the __Installation and Running from Source__ section below to create an executable 
binary for other environemnts* 


3. The console displays results output 

## Command Line Flags 

Use the following command line flags to configure your load test. 
You can define options in a `YAML` config file as well. Reference: 

```
    --duration      Int      Number of seconds to run test
    --users         Int      Number of users to simulate 
    --ramp          Int      Number of seconds between starting users
    --wait          Int      Number of seconds to wait between requests  
    --config        string   YAML config file see YAML Config section below
    --save          string   Save results to CSV file
    --replace       string   Save and replace file if exists results to CSV file
    --contenttype   string   Type (e.g. application/json) for POSTING data
    --authtoken     string   Authorization token 
    --tokentemplate string   Expression to format authtoken in request header
```
## Installing and Running from Source Code

1. [Install Go](https://go.dev/doc/install) 

2. Clone Repo 

3. Open a command line terminal, navigate to the repo directory, and then enter these commands:

```
    go install
    go build  
```

That will create an executable named `khsload` in your directory. 
Execute it with this command: 

```
    ./khsload 
```
See previous sections for options.

## Plotting Results 

You can plot in a scatter graph saved to a `CSV` file.  

**Run a Load test saved to `test.csv` by using the following command:**

```
    ./khsload do http://keyholesoftware.com --users 4 --duration 20 --save test.csv 
```

**Generate a plot from the generated `test.csv`**

```
    ./khsload plot test.csv
```

That command creates a scatter-based graph in a file named `khsplot.png`. 

Example graph:

![](khsplot.png)

### Configuration YAML 

Instead of supplying command line flags, you can define test options in a `YAML` configuration file
by specifing file path in a `--config` flag. 

An example YAML file...
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
  - POST~<URL>~key=value&key2=value2
  - GET~<URL>

```
## POST Requests 
HTTP GET is the default request method. Prefix the URL with `POST~` to carry 
out an HTTP POST request, as shown here:

```
    ./khsload do "POST~http://<address>~key1=value&key2=value"
```
Supply key/value data after the address. 
Content type of `POST` data defaults to `application/json`. 
You may chnage it `application/x-www-form-urlencoded` using a `--contenttype` flag:

```
    ./khsload do "POST~https://<Your POST Address Here>~akey1=avalue&akey2=avalue" --contenttype application/x-www-form-urlencoded
```

## Token Based Authentication 
This utility supports load testing `TOKEN`-based authentication schemes. 
If an API has a persistent access token applicable to request headers,
you can specify it using the command line `--authtoken` flag (or define it in the `YAML` config). 

```
    ./khsload do <some url> --authtoken <auth token>
```

Tokens apply to request headers using the `tokentemplate` expression. 
This appies the token value to an `authorization` request `Header` field. 

``` 
    ./khsload do <some url> --authtoken <auth token> --tokentemplate "{{Bearer .}}"
```










