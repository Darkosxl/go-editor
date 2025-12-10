package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/u2takey/ffmpeg-go"
	"os/exec"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "myapp",
		Short: "saying hi to mom",
		Long: "saying hi to mom for yall",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world!")
		},
	}

	var versionCmd = &cobra.Command{
		Use: "version",
		Short: "print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v1.0.0")
		},
	}

	var cutVideoCmd = &cobra.Command{
		Use: "cut",
		Short: "Cut silent parts from a video",
		Long: "cut [input.mp4] [minimumDecibel (default=-35dB)] [silencepaddingMS (default=100ms)] [removesilenceslongerthanMS (default=1000ms)] [removetalksshorterthanMS (default=100ms)]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			
			// to check if the first argument is a video file
			cmdcheck := exec.Command("ffmpeg", "-v", "error", "-i", args[0], "-f", "null", "-")
			output, err := cmdcheck.CombinedOutput()
			if err != nil {
				fmt.Println("First argument must be a video file.")
				return
			}
			
			cmdstring = ""
			if len(args) > 1 {
				for i := 1; i < len(args); i++ {
					if reflect.TypeOf(args[i]) != "int" {
						fmt.Println("Optional arguments must be integers.")
						return
					}
					cmdstring += " " + args[i]
				}
			
			

			}
			fmt.Println("cutting video")
			
			cmd := exec.Command("ffmpeg", "-i",args[0], "-af", "silencedetect=n="+args[1],)

		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cutVideoCmd)
	rootCmd.Execute()
}