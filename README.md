# Keyhole Software Load Testing Utility

Command Line utility written in [Go](https://go.dev) that load tests API's with a specified simulated users. Reports throughput statistics and performance Graphs can be generated. 

## Running 

1. Clone Repository 

2. Open a command line in cloned directory and execute the following

```
    $ ./bin/macosx/load do https://keyholesoftware.com --users 2 --duration 30 
```
3. Results will be output to the console 

## Command Line Flags 

The following command line flags can be used to configure your load test

```
    --duration  Int      Number of seconds to run test
    --users     Int      Number of users to simulate 
    --ramp      Int      Number of seconds between starting users
    --wait      Int      Number of seconds to wait between requests  
    --config    string   YAML config file see YAML Config section below
    --save      string   Save results to CSV file
    --replace   string   Save and replace file if exists results to CSV file
```
## Installing and Running from Source Code

- Prerequisites: +[Install Go](https://go.dev/doc/install) 

1. Clone Repo 

2. Build 



### Configuration YAML 

By default 


