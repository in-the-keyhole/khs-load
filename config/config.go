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
package config

import (
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type yml struct {
	name string
}

func check(e error) {
	if e != nil {
		log.Fatalln("Error reading Config file ", e)
	}
}

type ConfigData struct {
	Users    int
	Duration int
	Ramp     int
	Wait     int
	Auth     struct {
		Url           string
		Userid        string
		Password      string
		GetToken      string
		SplitWith     string
		TokenizeUsing string
	}
	Url           []string
	AuthToken     string
	TokenTemplate string
}

var yamlConfig ConfigData = ConfigData{}

var templateToken string

func IsAuth() bool {

	return yamlConfig.Auth.Url != ""

}

// SetName receives a pointer to Foo so it can modify it.
func SetDuration(duration int) {
	yamlConfig.Duration = duration
}

func AddUrl(url string) {
	m := "GET"
	u := url
	t := strings.Split(url, ",")
	if len(t) == 2 {
		m = t[0]
		u = t[1]
	}

	yamlConfig.Url = []string{m + "," + u}

}

func Duration() int {

	if yamlConfig.Duration == 0 {
		yamlConfig.Duration = 10
	}

	return yamlConfig.Duration
}

func Url() []string {
	return yamlConfig.Url
}

func Wait() int {

	return yamlConfig.Wait
}

func Ramp() int {

	if yamlConfig.Ramp == 0 {
		yamlConfig.Ramp = 1
	}

	return yamlConfig.Ramp
}

func Users() int {

	if yamlConfig.Users == 0 {
		yamlConfig.Users = 1
	}

	return yamlConfig.Users
}

// SetName receives a pointer to Foo so it can modify it.
func SetUsers(users int) {
	yamlConfig.Users = users
}

func TokenTemplate() string {

	return yamlConfig.TokenTemplate

}

func SetAuthToken(token string) {
	yamlConfig.AuthToken = token
}

func AuthToken() string {

	// only need to apply token template once
	if templateToken == "" {

		tmpl, err := template.New("").Parse(yamlConfig.TokenTemplate)
		if err != nil {
			log.Fatalf("Error Parsing Auth Token template: %s", err)
		}

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, yamlConfig.AuthToken)
		if err != nil {
			log.Fatalf("Error Executing token template: %s", err)
		}

		templateToken = tpl.String()

	}

	//return "Bearer " + yamlConfig.AuthToken
	return templateToken
}

func AuthUrl() string {
	return yamlConfig.Auth.Url
}

func UserId() string {
	return yamlConfig.Auth.Userid
}

func Password() string {
	return yamlConfig.Auth.Password
}

func SplitWith() string {

	if yamlConfig.Auth.SplitWith == "" {
		yamlConfig.Auth.SplitWith = ":"
	}

	return yamlConfig.Auth.SplitWith
}

func GetToken() string {

	if yamlConfig.Auth.GetToken == "" {
		yamlConfig.Auth.GetToken = "TOKEN"
	}

	return yamlConfig.Auth.GetToken
}

func TokenizeUsing() string {

	if yamlConfig.Auth.TokenizeUsing == "" {
		yamlConfig.Auth.TokenizeUsing = ","
	}

	return yamlConfig.Auth.TokenizeUsing
}

func Load(file string) {

	configFile := "config.yaml"

	if file != "" {

		configFile = file

	}

	//Config := ConfigData{}
	//yamlConfig = ConfigData{}

	data, err := os.ReadFile(configFile)
	check(err)
	//fmt.Print(string(dat))

	err = yaml.Unmarshal([]byte(data), &yamlConfig)
	if err != nil {
		log.Fatalf("Error UnMarshalling Config file: %v", err)
	}

}
