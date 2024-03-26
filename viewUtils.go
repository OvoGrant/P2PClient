package main

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	NUMBER     = "NUMBER |"
	ADDRESS    = "ADDRESS |"
	THROUGHPUT = "THROUGHPUT (Mb/s) |"
	LATENCY    = "LATENCY(ms) |"
)

func printPeerTable(delay []peerDelay) {
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(NUMBER, ADDRESS, THROUGHPUT, LATENCY)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i, v := range delay {
		tbl.AddRow(i, v.addr, v.rate, v.latency)
	}

	tbl.Print()

}
