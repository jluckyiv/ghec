/*
Copyright © 2023 Jackson Lucky <jack@jacksonlucky.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/jluckyiv/ghec"
	"github.com/spf13/cobra"
)

// elemCmd represents the elem command
var elemCmd = &cobra.Command{
	Use:   "elem",
	Short: "Add specific element",
	Run: func(cmd *cobra.Command, _ []string) {
		any, _ := cmd.Flags().GetBool("any")
		enh := ghec.EnhanceSpecificElement
		desc := "specific"
		if any {
			enh = ghec.EnhanceAnyElement
			desc = "any"
		}
		e := ghec.
			NewEnhancement(enh).
			WithMultipleTarget(numTargets).
			WithLevel(ghec.Level(level)).
			WithPreviousEnhancements(ghec.PreviousEnhancements(previousEnhancements))
		fmt.Printf("Add %s element costs %d", desc, e.Cost())
	},
}

func init() {
	rootCmd.AddCommand(elemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// elemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	elemCmd.Flags().BoolP("any", "a", false, "Any element enhancement")
}
