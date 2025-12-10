package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
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
	
		Long: "cut [input.mp4] [minimumDecibel (default=-35dB)] [removesilenceslongerthan (default=1s)] [silencepadding (default=0.5s)]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			
			// to check if the first argument is a video file
			cmdcheck := exec.Command("ffmpeg", "-v", "error", "-i", args[0], "-t", "1", "-f", "null", "-")
			output, err := cmdcheck.CombinedOutput()
			if err != nil {
				fmt.Println("First argument must be a video file.")
				return
			} 
			
			fmt.Println("man I have no idea if shit works or no")
			
			argsv2 := [4]string{"-35", "1", "0.5", "0.1"};

			if len(args) > 1 {
				for i := 1; i < len(args); i++ {
					_, err = strconv.ParseFloat(args[i], 64)
					
					if err != nil {
						fmt.Println("Optional arguments must be float.")
						return
					}

					argsv2[i-1] = args[i]
				}

			}
			
			fmt.Println("cutting video")

			ffmpegrun := exec.Command("ffmpeg", "-i",args[0], "-af", "silenceremove=stop_threshold="+argsv2[0]+"dB:stop_duration="+argsv2[1]+":stop_silence="+argsv2[2]+":start_periods=1:stop_periods=-1", "-c:v", "copy", "output.mp4")
			output, err = ffmpegrun.CombinedOutput()
			
			if err != nil {
				fmt.Println(string(output))
				fmt.Println("Error cutting video.")
				fmt.Println(err)
				return
			}
			fmt.Println(string(output))
			fmt.Println("Cut Successful!")

		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cutVideoCmd)
	rootCmd.Execute()
}