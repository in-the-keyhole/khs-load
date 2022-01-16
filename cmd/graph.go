/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	Use:   "graph",
	Short: "Line graph of ",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	s := fmt.Sprintf("%v Simulated Users made %v API requests with a throughput of \n %.2f TPS ", users, len(calls), throughput)
	p.Title.Text = s
	p.X.Label.Text = "Time (sec)"
	p.Y.Label.Text = "Response (sec)"

	p.Add(plotter.NewGrid())

	err := plotutil.AddScatters(p,
		createScatterData(bottleneckSamplings()))

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 4*vg.Inch, outputFile); err != nil {
		panic(err)
	}

}

func bottleneckSamplings() []Call {

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

func init() {
	rootCmd.AddCommand(graphCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
