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

package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"khsload/config"
	"log"
	"net/http"
	"strings"
)

var Token string

func Auth() {

	// Create a key/value map of strings
	kvPairs := make(map[string]string)
	kvPairs["email"] = config.UserId()
	kvPairs["password"] = config.Password()

	// Make this JSON
	postJson, err := json.Marshal(kvPairs)
	if err != nil {
		panic(err)
	}
	//	fmt.Printf("Sending JSON string '%s'\n", string(postJson))

	// http.POST expects an io.Reader, which a byte buffer does
	postContent := bytes.NewBuffer(postJson)

	log.Println("Authenticating")
	// Send request to OP's web server
	resp, err := http.Post(config.AuthUrl(), "application/json", postContent)
	if err != nil {

		log.Fatalln(err)

	}

	if resp.StatusCode == 200 {

		fmt.Println("Authentication Successful -> ", resp.StatusCode)
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}

		type Token struct {
			Value string
			Del   string
		}

		/*	funcMap := template.FuncMap{
				// The name "title" is what the function will be called in the template text.
				"split": strings.Split,
			}

			const templateText = "Input: Output 0: {{split .Value .Del}}"

				tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
				if err != nil {
					log.Fatalf("parsing: %s", err)
				}

				// Run the template to verify the output.
				//err = tmpl.Execute(os.Stdout, Token{"Hello World", " "})
				if err != nil {
					log.Fatalf("execution: %s", err)
				}   */

		tokenize := config.TokenizeUsing()
		token := config.GetToken()
		splitWith := config.SplitWith()

		fmt.Println("Tokenizing Auth response body using -> ", tokenize)
		fmt.Println("Getting token using -> ", token)
		fmt.Println("Splitting Token to get Value with ->", splitWith)
		body := string(b)
		tokens := strings.Split(body, tokenize)

		for _, v := range tokens {

			if strings.Contains(strings.ToUpper(v), strings.ToUpper(token)) {

				t := strings.Split(v, splitWith)[1]
				token := strings.Trim(t, "\"")
				config.SetAuthToken(token)

			}
		}

	} else {

		log.Fatalln("Authentication Failed -> ", resp.StatusCode)

	}

}
