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
package cmd

import (
	"fmt"
	"khsload/auth"
	"khsload/config"
	"khsload/invoke"
	"khsload/stats"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "do",
	Short: "API load testing by Keyhole Software (http://keyholesoftware.com)",
	Long: `Keyhole Software Load testing. This utility will load test single or multiple URI's. A URI, 
	can be specified as a command line argument, or multiple URI's can be listed in a YAML file. `,
	Run: func(cmd *cobra.Command, args []string) {

		applyFlags(cmd)

		saveFile, _ := cmd.Flags().GetString("save")

		versionFlag, _ := cmd.Flags().GetString("version")

		if versionFlag != "" {

			fmt.Println("Keyhole Software Api load tester Version - BETA")
			return

		}

		if saveFile != "" {

			_, err := os.Stat(saveFile)

			if err == nil {

				fmt.Println("File ", saveFile, " exists, remove or run use the --replace flag to overwrite ")
				return
			}

		}

		replaceFile, _ := cmd.Flags().GetString("replace")

		if replaceFile != "" {
			saveFile = replaceFile
		}

		// apply flags

		if config.IsAuth() {
			auth.Auth()
		}

		users := config.Users()
		ramp := config.Ramp()
		duration := config.Duration() + (users * ramp)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(duration))

		if len(config.Url()) == 0 && len(args) == 0 {
			fmt.Println("API URL required as argument, or define in config.yaml")

		} else {

			if len(config.Url()) == 0 {
				config.AddUrl(args[0])
			}

			stats.Users(users)
			stats.Duration(config.Duration())

			urls := config.Url()
			//var urls []string
			//urls = append(urls, "https://keyholesoftware.com")
			stats.ReportBegin(urls)

			log.Println("Begin")
			for i := 0; i < users; i++ {

				go doUser(ctx, urls, i)

				seconds := strconv.Itoa(ramp)
				dur, _ := time.ParseDuration(seconds + "s")
				fmt.Print("Ramping->User ", i, " ")
				time.Sleep(dur)

			}

			select {
			case <-ctx.Done():
				log.Println(".")
				log.Println("End")
				stats.ReportEnd(urls)
				if saveFile != "" {
					stats.Save(saveFile, urls)
				}
				return
			case <-time.After(time.Second * time.Duration(duration)):
				//Here I'm actually ending it earlier than the timeout with cancel().
				fmt.Println(".")
				log.Println("End")
				stats.ReportEnd(urls)
				if saveFile != "" {
					stats.Save(saveFile, urls)
				}
				cancel()
			}

		}

	},
}

func applyFlags(cmd *cobra.Command) {

	configFile, _ := cmd.Flags().GetString("config")

	usersFlag, _ := cmd.Flags().GetString("users")
	iusr, _ := strconv.Atoi(usersFlag)
	config.SetUsers(iusr)

	durationFlag, _ := cmd.Flags().GetString("duration")
	idur, _ := strconv.Atoi(durationFlag)
	config.SetDuration(idur)

	rampFlag, _ := cmd.Flags().GetString("ramp")
	iramp, _ := strconv.Atoi(rampFlag)
	config.SetRamp(iramp)

	waitFlag, _ := cmd.Flags().GetString("wait")
	iwait, _ := strconv.Atoi(waitFlag)
	config.SetWait(iwait)

	if configFile != "" {

		fmt.Println("Info->Configuration Found, values in config file will override command line args")
		config.Load(configFile)

	}

}

func doUser(ctx context.Context, urls []string, user int) {

	seconds := strconv.Itoa(config.Wait())
	wait, _ := time.ParseDuration(seconds + "s")
	client := http.Client{}
	for _, url := range urls {
		go doInvoke(ctx, url, client, user)
		time.Sleep(wait)
	}

}

func doInvoke(ctx context.Context, url string, client http.Client, user int) {
	count := 1
	for {
		select {
		case <-ctx.Done():
			log.Println("Done working")
			return
		default:
			invoke.Invoke(url, count, client, user)
			count++

		}
		time.Sleep(time.Millisecond)
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	testCmd.PersistentFlags().String("config", "", "YAML Configuration File")
	testCmd.PersistentFlags().String("save", "", "save API stats to file, in CSV format")
	testCmd.PersistentFlags().String("replace", "", "save API stats to file, replace file if exists, in CSV format")
	testCmd.PersistentFlags().String("users", "1", "Simulated Users")
	testCmd.PersistentFlags().String("wait", "1", "Seconds to wait between requests")
	testCmd.PersistentFlags().String("duration", "20", "Duration to Run Test")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//testCmd.Flags().Boolp("toggle", "t", false, "Help message for toggle")

}
