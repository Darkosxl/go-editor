package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func cleanup() {
	os.Remove("cutlist.txt")
	os.Remove("segments.txt")
	matches, _ := filepath.Glob("segment*.mp4")
	for _, file := range matches {
		os.Remove(file)
	}
}
func check(e error) {
	if e != nil {
		fmt.Println(e)
		cleanup()
		panic(e)
	}
}
func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Simple Video Editor",
		Long:  "Simple Video Editor with high precision cuts in seconds",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world!")
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v1.1.2")
		},
	}

	var cutVideoCmd = &cobra.Command{
		Use:   "cut",
		Short: "Implement cuts in a video with high precision in 10 seconds",
		Long:  "cut [input.mp4] [minimumDecibel (default=-35dB, input positive numbers they will be flipped to negative)] [removesilenceslongerthan (default=1s)] [silencepadding (default=0.5s)] [removetalksshorterthan (default=0s)] [lengthofvideotocut (default=0s (entirevideo) )]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// to check if the first argument is a video file
			cmdcheck := exec.Command("ffmpeg", "-v", "error", "-i", args[0], "-t", "1", "-f", "null", "-")
			output, err := cmdcheck.CombinedOutput()
			check(err)
			cleanup()

			argsv2 := [5]string{"35", "1", "0.5", "0", "0"}

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

			silencepadding, _ := strconv.ParseFloat(argsv2[2], 64)
			mintalk, _ := strconv.ParseFloat(argsv2[3], 64)

			fmt.Println("Reading silences in the video")
			if argsv2[4] == "0" {
				ffmpegdetect := exec.Command("ffmpeg", "-i", args[0], "-af", "silencedetect=noise=-"+argsv2[0]+"dB:d="+argsv2[1], "-f", "null", "-")
				output, err = ffmpegdetect.CombinedOutput()
				check(err)
			} else {
				ffmpegdetect := exec.Command("ffmpeg", "-i", args[0], "-t", argsv2[4], "-af", "silencedetect=noise="+argsv2[0]+"dB:d="+argsv2[1], "-f", "null", "-")
				output, err = ffmpegdetect.CombinedOutput()
				check(err)
			}

			output_string := string(output)
			output_split := strings.Split(output_string, "silence_")

			start := "0"
			end := "0"
			f, err := os.Create("cutlist.txt")
			check(err)
			defer f.Close()

			for i := 0; i < len(output_split); i++ {
				if strings.Contains(output_split[i], "start:") {
					r, rerr := regexp.Compile("start: [0-9]+\\.?[0-9]*")
					check(rerr)
					end = r.FindString(output_split[i])
					end = strings.TrimPrefix(r.FindString(output_split[i]), "start: ")

					// Only write segment if there's actual content to keep (not 0-length)
					if start != end {
						n3, nerr := f.WriteString("cut_keep:" + start + ":" + end + "\n")
						check(nerr)
						fmt.Println(n3)
					}
				}

				if strings.Contains(output_split[i], "end:") {
					r, rerr := regexp.Compile("end: [0-9]+\\.?[0-9]*")
					check(rerr)
					start = r.FindString(output_split[i])
					start = strings.TrimPrefix(r.FindString(output_split[i]), "end: ")

				}
			}
			f.WriteString("cut_keep:" + start + ":99999\n")
			file, err := os.Open("cutlist.txt")
			check(err)
			defer file.Close()
			scanner := bufio.NewScanner(file)

			fmt.Println("Cutting video")
			counter := 0
			filesegments, err := os.Create("segments.txt")
			check(err)
			defer filesegments.Close()
			prevEnd := 0.0 // Track where the previous segment actually ended (with padding)
			for scanner.Scan() {

				cutkeep := strings.Split(scanner.Text(), ":")
				start, _ := strconv.ParseFloat(cutkeep[1], 64)
				end, _ := strconv.ParseFloat(cutkeep[2], 64)

				// Skip this segment if it overlaps with the previous one
				if start < prevEnd {
					// Adjust start to avoid overlap
					start = prevEnd
					cutkeep[1] = strconv.FormatFloat(start, 'f', -1, 64)
				}

				// Only cut if there's meaningful content after adjustment
				if end-start > mintalk {
					duration := end - start + silencepadding
					// Use re-encoding for precision (remove -c copy)
					ffmpegcut := exec.Command("ffmpeg", "-y", "-ss", cutkeep[1], "-i", args[0], "-t", strconv.FormatFloat(duration, 'f', -1, 64), "-c:v", "libx264", "-preset", "ultrafast", "-c:a", "aac", "segment"+strconv.Itoa(counter)+".mp4")
					cutOutput, err := ffmpegcut.CombinedOutput()
					if err != nil {
						fmt.Println("Segment", counter, "failed. cutkeep:", cutkeep)
						fmt.Println(string(cutOutput))
						check(err)
					}
					filesegments.WriteString("file" + " 'segment" + strconv.Itoa(counter) + ".mp4'\n")
					counter++
					prevEnd = start + duration // Update where this segment actually ends
				}
			}

			f.Sync()
			file.Sync()
			filesegments.Sync()

			outputName := strings.TrimSuffix(args[0], ".mp4") + "_cut.mp4"
			fmt.Println("Concatenating video")
			ffmpegrun := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", "segments.txt", "-c", "copy", outputName)

			ffmpegoutput, err := ffmpegrun.CombinedOutput()
			fmt.Println(string(ffmpegoutput))
			check(err)

			os.Remove("cutlist.txt")
			os.Remove("segments.txt")

			for i := 0; i < counter; i++ {
				os.Remove("segment" + strconv.Itoa(i) + ".mp4")
			}
			fmt.Println(string(output))

			fmt.Println("Cut Successful!")

		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cutVideoCmd)
	rootCmd.Execute()
}
