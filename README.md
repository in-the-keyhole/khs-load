# Keyhole Software Load Testing Utility

Command Line utility written in [Go](https://go.dev) that load tests API's with a specified simulated users. Reports throughput statistics and performance Graphs can be generated. 

## Running 

1. Clone Repository 

2. Open a command line in cloned directory and execute the following

```
    $ ./bin/macosx/khsload do https://keyholesoftware.com --users 2 --duration 30 
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

1. [Install Go](https://go.dev/doc/install) 

2. Clone Repo 

3. Open a command line terminal and navigate to the repo directory and enter the following commands 

```
    $ go install
    $ go build  
```

This will create an executable named `khsload` in your directory. that can be executed with this command. See previous sections for options.

```
    $./khsload 
```

## Plotting Results 

Results saved to a `CSV` file can be plotted to a scatter graph with the following command 

** Run a Load test and save to `test.csv` with the following command **

```
    $./khsload do http://keyholesoftware.com --users 4 --duration 20 --save test.csv 
```

** Generate a plot with the resulting `test.csv` **

```
    $./khsload plot test.csv
```

A scatter based Graph will be plotted to a file named `khsplot.png` 

Here's and example graph

![](khsplot.png)








### Configuration YAML 

Options can be specified 

By default 





