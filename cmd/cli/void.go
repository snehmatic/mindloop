package cli

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var voidCmd = &cobra.Command{
	Use:     "void [minutes]",
	Short:   "Enter the void",
	Long:    `Enter the void for a specified number of minutes to focus on your breathing. Default is 5 minutes.`,
	Example: `mindloop void 5`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		minutes := 5
		if len(args) > 0 {
			var err error
			minutes, err = strconv.Atoi(args[0])
			if err != nil || minutes <= 0 {
				utils.PrintErrorln("Minutes must be a positive integer.")
				return
			}
		}

		duration := time.Duration(minutes) * time.Minute
		utils.PrintRocketln(fmt.Sprintf("Entering the void for %d minutes...", minutes))
		time.Sleep(1 * time.Second)

		// Hide cursor
		fmt.Print("\033[?25l")
		// Ensure cursor is shown on exit
		defer fmt.Print("\033[?25h")

		// Handle interrupt signals (Ctrl+C) gracefully
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Print("\033[?25h")     // Show cursor
			fmt.Print("\033[H\033[2J") // Clear screen
			fmt.Println("\nThe void demands your attention.")
			os.Exit(0)
		}()

		end := time.Now().Add(duration)
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		start := time.Now()

		// Clear screen to enter the void
		fmt.Print("\033[H\033[2J")

		for {
			now := time.Now()
			if now.After(end) {
				break
			}
			remaining := end.Sub(now)
			elapsed := now.Sub(start)

			// 16-second Box Breathing cycle
			// Inhale: 4s, Hold: 4s, Exhale: 4s, Hold: 4s. Total: 16s
			cycle := math.Mod(elapsed.Seconds(), 16.0)

			var phase string
			var progress float64

			if cycle < 4.0 {
				phase = "INHALE"
				t := cycle / 4.0
				invT := 1.0 - t
				progress = 1.0 - (invT * invT * invT) // ease out
			} else if cycle < 8.0 {
				phase = "HOLD"
				progress = 1.0
			} else if cycle < 12.0 {
				phase = "EXHALE"
				t := (cycle - 8.0) / 4.0
				progress = 1.0 - (t * t * (3.0 - 2.0*t)) // smoothstep
				if progress < 0 {
					progress = 0
				}
			} else {
				phase = "HOLD"
				progress = 0.0
			}

			// Render the UI
			// Move cursor to top-left to redraw smoothly
			fmt.Print("\033[H")

			fmt.Print("\033[K\n\033[K\n\033[K\n")
			// 15 chars, centered at 24 = pad 17
			fmt.Print("                 \033[1;37mT H E   V O I D\033[0m\033[K\n")
			// 19 chars, centered at 24 = pad 15
			fmt.Print("               \033[38;5;238m───────────────────\033[0m\033[K\n\033[K\n\033[K\n")

			// Center the phase text around column 24
			paddingPhase := 24 - (len(phase) / 2)
			fmt.Printf("%s\033[1;36m%s\033[0m\033[K\n\033[K\n\033[K\n", strings.Repeat(" ", paddingPhase), phase)

			// Breathing visual
			minWidth := 3
			maxWidth := 12

			// Interpolate between minWidth and maxWidth
			currentWidth := minWidth + int(math.Round(progress*float64(maxWidth-minWidth)))

			// Each block is "█ ". We trim the trailing space to ensure perfect centering.
			// Visual width is currentWidth * 2 - 1
			visualWidth := currentWidth*2 - 1
			paddingBar := 24 - (visualWidth / 2)

			barStr := strings.TrimSpace(strings.Repeat("█ ", currentWidth))
			bar := strings.Repeat(" ", paddingBar) + barStr

			fmt.Printf("\033[1;36m%s\033[0m\033[K\n\033[K\n\033[K\n", bar)

			// Timer
			paddingTimer := 24 - 2 // 22 spaces
			fmt.Printf("%s\033[38;5;245m%02d:%02d\033[0m\033[K\n\033[K\n", strings.Repeat(" ", paddingTimer), int(remaining.Minutes()), int(remaining.Seconds())%60)

			// Exit text
			paddingExit := 24 - 12 // 12 spaces
			fmt.Printf("%s\033[38;5;238m[ Press Ctrl+C to exit ]\033[0m\033[K\n", strings.Repeat(" ", paddingExit))
			// Clear rest of screen
			fmt.Print("\033[J")

			<-ticker.C
		}

		// Clear screen on exit
		fmt.Print("\033[H\033[2J")
		utils.PrintSuccessln("You have successfully returned from the void.")
	},
}

func init() {
	rootCmd.AddCommand(voidCmd)
}
