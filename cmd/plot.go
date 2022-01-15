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
	"bufio"
	"fmt"
	"image/color"
	"khsload/stats"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Call struct {
	Milli int64
	Api   string
	Time  int
	Bytes int
	User  int
}

var urls = []string{}

var calls = []Call{}

var outputFile = "khsplot.png"

var users int

var throughput float64

type CallSorter []Call

func (a CallSorter) Len() int           { return len(a) }
func (a CallSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CallSorter) Less(i, j int) bool { return a[i].Milli < a[j].Milli }

// plotCmd represents the plot command
var plotCmd = &cobra.Command{
	Use:   "plot <CSV filename>",
	Short: "Generated Graph Plots from 'saved' load file",
	Long: `When performing a Load Test output from the test can be saved to a file, this file 
	is ingested by this command to produce a Graph Plot PNG'

	`,
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

		//	max := calcMax()
		//	users := calcUsers()

		//	scatter(max, users)
		doGraph()

		fmt.Printf("Plot generated to %s ... \n", outputFile)

	},
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func calcUsers() int {

	users := 0
	for _, c := range calls {

		if c.User > users {
			users = c.User
		}

	}

	return users

}

func calcMax() int64 {

	/*

		time := 0
		for _, c := range calls {

			if c.Time > time {
				time = c.Time
			}

		}*/

	seconds := (calls[len(calls)-1].Milli - calls[0].Milli) / 1000

	return seconds

}

func createCallData(data []string) (string, []Call) {
	var results []Call

	var title string

	for i, line := range data {

		if i == 0 {

			title := strings.Split(line, ",")

			var e error
			users, e = strconv.Atoi(title[0])

			check(e)

			t, e := strconv.ParseFloat(title[1], 64)

			throughput = t

			check(e)

		} else if i == 1 {

			urls = strings.Split(line, ",")

		} else {

			items := strings.Split(line, ",")

			c := Call{}

			c.Milli, _ = strconv.ParseInt(items[0], 10, 64)
			c.User, _ = strconv.Atoi(items[2])
			c.Api = items[3]
			if len(items) == 7 {
				c.Time, _ = strconv.Atoi(items[5])
				c.Bytes, _ = strconv.Atoi(items[6])

			} else {
				c.Time, _ = strconv.Atoi(items[4])
				c.Bytes, _ = strconv.Atoi(items[5])
			}

			results = append(results, c)

		}

	}

	return title, results
}

func load(file string) (string, []Call) {

	// open file
	f, err := os.Open(file)
	check(err)

	// remember to close the file at the end of the program
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	title, results := createCallData(lines)

	return title, results

}

func sortCalls() {

	sort.Sort(CallSorter(calls))

}

func doGraph() {

	rand.Seed(int64(0))

	p := plot.New()

	s := fmt.Sprintf("%v Simulated Users made %v API requests with a throughput of \n %.2f TPS ", users, len(calls), throughput)
	p.Title.Text = s
	p.X.Label.Text = "Time (sec)"
	p.Y.Label.Text = "Response (sec)"

	p.Add(plotter.NewGrid())

	err := plotutil.AddScatters(p,
		createScatterData(samplings()))

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 4*vg.Inch, outputFile); err != nil {
		panic(err)
	}

}

func scatter(data []Call) {

	rand.Seed(int64(0))

	p := plot.New()
	p.Title.Text = stats.Title()
	p.X.Label.Text = ""
	p.Y.Label.Text = "Response (Sec)"
	p.X.Label.Position = draw.PosRight
	p.Y.Label.Position = draw.PosTop
	p.X.Min = 0
	p.X.Max = float64(len(data))
	p.Y.Min = 0
	//p.Y.Max = float64(max)

	p.Add(plotter.NewGrid())

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "axis_labels.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(randomPoints(15))
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	p.Add(s)
	p.Legend.Add("scatter", s)

	//plotutil.AddScatters(p, createScatterData())

}

func samplings() []Call {

	results := make([]Call, 0)
	count := 0

	//	sampleSize := 20
	for i := range calls {

		//	if count%sampleSize == 0 {
		results = append(results, calls[i])
		//	}
		count++

	}

	return results

}

func createScatterData(data []Call) plotter.XYs {
	pts := make(plotter.XYs, len(data))

	start := calls[0].Milli

	for i := range data {

		pts[i].Y = float64(data[i].Time / 100.00)

		interval := float64((data[i].Milli - start) / 1000.00 / 1000.00)

		pts[i].X = float64(interval)

	}

	return pts
}

func init() {
	rootCmd.AddCommand(plotCmd)

	plotCmd.PersistentFlags().String("file", "khsplot.png", "File to save plotted GRAPH to")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
