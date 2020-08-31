package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func PrintTable(headers []string, lines ...[]interface{}) {
	var tbl table.Table
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	// Compile Headings for writer
	expHeaders := make([]interface{}, len(headers))
	for i, header := range headers {
		expHeaders[i] = header
	}

	tbl = table.New(expHeaders...).WithPadding(2)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, line := range lines {
		expLine := make([]interface{}, len(line))
		// Compile interface slice
		for i, entry := range line {
			expLine[i] = entry
		}
		tbl.AddRow(expLine...)
	}

	tbl.Print()
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
