package cmd

import (
	"fmt"
	"os"
	"weather/client"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the weather",
	Long:  `Get the weather. Args: specified city`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		city := args[0]

		apiKey := os.Getenv("API_KEY")
		isF, _ := cmd.Flags().GetBool("F")
		isKMH, _ := cmd.Flags().GetBool("kmh")

		cl := client.New(apiKey)
		res, err := cl.GetWeather(city, isF, isKMH)

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().Bool("F", false, "Get units: fahrenheit")
	getCmd.Flags().Bool("kmh", false, "Get units: km/h")
}
