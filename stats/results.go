// Package stats
/*
Copyright 2022 Keyhole Software (http://keyholesoftware.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
package stats

import (
	"fmt"
	"khsload/config"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Results struct {
	Success  int
	Failures int
	Duration int
	Users    int
}

type Call struct {
	Ts    int
	Api   string
	Time  int
	Bytes int
	User  int
}

//var calls = make(map[string]int)
var calls = []Call{}

var Singleton = &Results{}

var successMutex = sync.RWMutex{}

var failureMutex = sync.RWMutex{}

var callMutex = sync.RWMutex{}

func AddCall(id string, mil int, bytes int, user int) {

	callMutex.Lock()
	c := Call{Ts: int(time.Now().UnixNano() / 1000), Api: id, Time: mil, Bytes: bytes, User: user}
	calls = append(calls, c)
	callMutex.Unlock()
}

func throughputUri(uri string) (float32, float32) {

	total := 0
	count := 0
	bytes := 0

	for i := 0; i < len(calls); i++ {

		c := calls[i]
		if strings.EqualFold(c.Api, uri) {
			total += c.Time
			bytes += c.Bytes
			count = count + 1
		}
	}

	if count == 0 {
		return -1, -1
	}

	// convert to seconds
	ms := float32(total / count)

	return 1000 / ms, float32((1000.00 / ms) * float32(bytes/count))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Bottleneck(url string) int {

	//previous := 0
	count := 0

	for i := 0; i < len(calls); i++ {

		if calls[i].Api == url {

			count++

			if (calls[i].Time) > 1000 {

				return count

			} else {

				//previous = calls[i].Time

			}

		}

	}

	return -1

}

func Sla(url string) int {

	//previous := 0
	count := 0

	for i := 0; i < len(calls); i++ {

		if calls[i].Api == url {

			count++

			if (calls[i].Time) > config.Sla() {

				return count

			} else {

				//previous = calls[i].Time

			}

		}

	}

	return -1

}

func GetCalls() []Call {
	return calls
}

// Save in a CSV format
func Save(file string, urls []string) {

	f, err := os.Create(file)
	check(err)

	defer f.Close()

	t, _ := throughput()
	sla := config.Sla()
	s := fmt.Sprintf("%v,%f,%v\n", Singleton.Users, t, sla)

	f.WriteString(s)

	f.WriteString(strings.Join(urls, ",") + "\n")

	for i := 0; i < len(calls); i++ {

		s := strconv.Itoa(calls[i].Ts) + "," + strconv.Itoa(calls[i].User) + "," + calls[i].Api + "," + strconv.Itoa(calls[i].Time) + "," + strconv.Itoa(calls[i].Bytes) + "\n"
		f.WriteString(s)

	}

}

func Users(users int) {
	Singleton.Users = users
}

func Duration(duration int) {
	Singleton.Duration = duration
}

func Sucess() {

	successMutex.Lock()
	Singleton.Success++
	successMutex.Unlock()
}

func OutsideSla() int {

	outsideSla := 0

	for _, c := range calls {

		if c.Time > config.Sla() {
			outsideSla++
		}

	}

	return outsideSla

}

func throughput() (float32, float32) {

	total := 0
	count := 0
	bytes := 0
	for i := 0; i < len(calls); i++ {

		c := calls[i]
		total += c.Time
		bytes += c.Bytes
		count++
	}

	if count == 0 {
		return -1, -1
	}

	// convert to seconds
	ms := float32(total / count)

	return 1000 / ms, float32((1000.00 / ms) * float32(bytes/count))

}

func Failure() {

	failureMutex.Lock()
	Singleton.Failures++
	failureMutex.Unlock()

}

func Title() string {

	return strconv.Itoa(Singleton.Users) + " USERS will be RAMPED for Testing every " + strconv.Itoa(config.Ramp()) + " seconds and run for " + strconv.Itoa(Singleton.Duration) + " Seconds"

}

func GetResults() Results { return *Singleton }

func ReportBegin(urls []string) {

	wait := config.Wait()

	fmt.Println("---")
	fmt.Println(Title())
	fmt.Println("Each API call will WAIT for ", wait, " seconds between calls")
	fmt.Println("---")
	fmt.Println("Testing the Following URIs:")

	for i := 0; i < len(urls); i++ {

		fmt.Println("-", urls[i])

	}

	fmt.Println("---")

}

func ReportEnd(urls []string) {

	fmt.Println("-------- R E P O R T ----------")
	fmt.Println("Load Testing Complete ")
	fmt.Println("Duration:", Singleton.Duration, " Seconds")
	fmt.Println("Users:", Singleton.Users)
	fmt.Println("Success:", Singleton.Success)
	fmt.Println("Failures:", Singleton.Failures)

	fmt.Printf("SLA (%v ms): %v  out of %v transactions took over %v milliseconds to exeucte", config.Sla(), OutsideSla(), Singleton.Success, config.Sla())

	t, b := throughput()
	fmt.Printf("\nTotal Throughput ( %v ) %6.2f TPS %6.2f bytes per second", len(urls), t, b)

	for i := 0; i < len(urls); i++ {

		ms, bytes := throughputUri(urls[i])

		fmt.Printf("\nThroughput (%s) %6.2f TPS %6.2f bytes per second", urls[i], ms, bytes)

		bottleneck := Bottleneck((urls[i]))

		if bottleneck >= 0 {

			fmt.Printf("\nSlowdown Detected after %v requests", bottleneck)

		}

	}

	fmt.Printf("\n")

}
