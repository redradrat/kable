package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func PrintTable(headers []string, lines ...[]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader(headers)
	table.SetBorder(false)  // Set Border to false
	table.AppendBulk(lines) // Add Bulk Data
	table.Render()
}

func PrintError(format string, a ...interface{}) {
	fmt.Println(fmt.Errorf(color.RedString("! "+format, a...)))
	os.Exit(1)
}

func PrintSuccess(format string, a ...interface{}) {
	fmt.Println(color.GreenString("✔ "+format, a...))
}

func PrintMsg(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func PrintWarning(format string, a ...interface{}) {
	fmt.Println(color.YellowString("! "+format, a...))
}
