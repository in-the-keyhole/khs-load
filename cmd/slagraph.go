// Package cmd
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

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "sla",
	Short: "Renders a PNG file from a CSV file of response times vs time offsets",
	Long: `Produces a PNG file rendition of a ling graph connecting points defined in a CSV file.
Each point contains a response time versus a time offset. A rendered title states
the total number of requests with the total that exceeded the SLA limit.`,

	Run: func(cmd *cobra.Command, args []string) {

		file, _ := cmd.Flags().GetString("file")

		if file != "" {
			outputFile = file
		}

		if len(args) == 0 {

			fmt.Println("Input CSV file required")
			return

		}

		//	var title string
		_, calls = load(args[0])

		sortCalls()

		doLineGraph()

		fmt.Printf("Line Graph generated to %s ... \n", outputFile)

	},
}

func doLineGraph() {

	p := plot.New()

	s := fmt.Sprintf("%v Simulated Users made %v API requests \n %v requests took longer than %v ms ", users, len(calls), SlaCount(), Sla())
	p.Title.Text = s
	p.X.Label.Text = "Time (sec)"
	p.Y.Label.Text = "Response (sec)"

	p.Add(plotter.NewGrid())

	err := plotutil.AddLinePoints(p,
		createScatterData(slaSamplings()))

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 4*vg.Inch, outputFile); err != nil {
		panic(err)
	}

}

func slaSamplings() []Call {

	results := make([]Call, 0)

	//	sampleSize := 20
	for i := range calls {

		if calls[i].Time > Sla() {
			results = append(results, calls[i])
		}

	}

	return results

}

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.PersistentFlags().String("file", "khssla.png", "File to save SLA GRAPH to")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
