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

package invoke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"khsload/config"
	"khsload/stats"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Invoke(aurl string, count int, client http.Client, user int) {

	urlItems := strings.Split(aurl, ",")

	var keyValues map[string]string
	data := url.Values{}
	var headerContent io.Reader = nil
	if len(urlItems) > 2 {

		keyValues = ParseKeyValues(urlItems[2])

		for k, v := range keyValues {
			data.Set(k, v)
		}

		if strings.ToLower(config.ContentType()) == "application/x-www-form-urlencoded" {

			headerContent = strings.NewReader(data.Encode())

		} else if strings.ToLower(config.ContentType()) == "application/json" {

			b := new(bytes.Buffer)

			json.NewEncoder(b).Encode(data)

			headerContent = bytes.NewReader(b.Bytes())

		}
	}

	req, err := http.NewRequest(urlItems[0], urlItems[1], headerContent)
	if err != nil {
		stats.Failure()
		log.Println(err)
	}

	token := config.AuthToken()
	req.Header.Set("authorization", token)
	req.Header.Set("Content-Type", config.ContentType())

	start := time.Now().UnixNano() / 1000000

	resp, err := client.Do(req)
	if err != nil {
		stats.Failure()
		log.Println(err)
	} else {

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

			duration := (time.Now().UnixNano() / 1000000) - start

			b, err := io.ReadAll(resp.Body)

			if err != nil {
				log.Println("Error reading response ", err)
			}

			stats.AddCall(aurl, int(duration), len(b), user)
			stats.Sucess()
			fmt.Print(".")

		} else {

			log.Println("Failed ->", aurl, " Status Code = ", resp.StatusCode)
			stats.Failure()

		}

		defer resp.Body.Close()

	}

}

func ParseKeyValues(keyValues string) map[string]string {

	results := map[string]string{}

	pairs := strings.Split(keyValues, "&")

	for i := 0; i < len(pairs); i++ {

		p := strings.Split(pairs[i], "=")
		results[p[0]] = p[1]
	}

	return results

}
