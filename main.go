// csvTochart — convert a bandwidth CSV to an offline ECharts HTML chart.
//
// Build:  go build -o csvTochart ./cmd/csvTochart
// Usage:  ./csvTochart [flags] <input.csv> [output.html]
//
// CSV format:
//   - Column 0: timestamps
//   - Remaining columns: bandwidth values (unit specified by --unit, default bps)
//   - A column named "Total", "Bandwidth(bps)", or "Bandwidth" is drawn as the
//     reference total line; all other columns are per-domain series.
//
// The output HTML is fully self-contained (no network required).
// Y-axis and tooltip use 1000-based prefixes (Kbps=1000bps, Mbps=1000Kbps, …).
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stanhui/csvtochart/internal/chart"
)

// unitMultipliers maps --unit flag value to its multiplier in bps.
var unitMultipliers = map[string]float64{
	"bps":  1,
	"kbps": 1e3,
	"mbps": 1e6,
	"gbps": 1e9,
	"tbps": 1e12,
}

func main() {
	fs := flag.NewFlagSet("csvTochart", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `csvTochart — 将带宽 CSV 转换为离线 ECharts HTML 图表

用法:
  csvTochart [flags] <input.csv> [output.html]

参数:
  input.csv     输入 CSV 文件（必填）
  output.html   输出 HTML 路径（可选，默认与 CSV 同名换 .html 后缀）

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
说明:
  CSV 格式要求：第一列为时间戳，其余列为带宽数值。
  列名为 "Total"、"Bandwidth(bps)" 或 "Bandwidth" 的列作为汇总线，其余列为域名线。

  Y 轴和 Tooltip 使用 1000 进制换算（1 Kbps = 1000 bps，1 Mbps = 1000 Kbps，以此类推）。

示例:
  csvTochart result.csv
  csvTochart --unit Mbps result.csv output.html
  csvTochart --unit Gbps --title "CDN Bandwidth" result.csv
`)
	}

	unitFlag := fs.String("unit", "bps", "CSV 中带宽数据的单位，可选: bps, Kbps, Mbps, Gbps, Tbps")
	titleFlag := fs.String("title", "", "图表标题（默认使用文件名）")

	// support mixed order: flags after positional args
	var flagArgs, posArgs []string
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
			if !strings.Contains(a, "=") && i+1 < len(os.Args) {
				name := strings.TrimLeft(a, "-")
				if f := fs.Lookup(name); f != nil {
					i++
					flagArgs = append(flagArgs, os.Args[i])
				}
			}
		} else {
			posArgs = append(posArgs, a)
		}
	}
	fs.Parse(append(flagArgs, posArgs...))
	args := fs.Args()

	if len(args) < 1 {
		fs.Usage()
		os.Exit(1)
	}

	csvPath := args[0]
	htmlPath := strings.TrimSuffix(csvPath, filepath.Ext(csvPath)) + ".html"
	if len(args) > 1 {
		htmlPath = args[1]
	}

	multiplier, ok := unitMultipliers[strings.ToLower(*unitFlag)]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown unit %q, valid values: bps, Kbps, Mbps, Gbps, Tbps\n", *unitFlag)
		os.Exit(1)
	}

	title := *titleFlag
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(csvPath))
	}

	if err := run(csvPath, htmlPath, multiplier, title); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("saved:", htmlPath)
}

func run(csvPath, htmlPath string, unitMultiplier float64, title string) error {
	raw, err := os.ReadFile(csvPath)
	if err != nil {
		return err
	}
	if len(raw) >= 3 && raw[0] == 0xEF && raw[1] == 0xBB && raw[2] == 0xBF {
		raw = raw[3:]
	}

	r := csv.NewReader(strings.NewReader(string(raw)))
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	all, err := r.ReadAll()
	if err != nil {
		return err
	}
	if len(all) < 2 {
		return fmt.Errorf("csv must have header + data rows")
	}

	header, dataRows := all[0], all[1:]

	xLabels := make([]string, len(dataRows))
	for i, row := range dataRows {
		if len(row) > 0 {
			xLabels[i] = row[0]
		}
	}

	var series []chart.Series

	for ci, name := range header[1:] {
		colIdx := ci + 1
		vals := make([]float64, len(dataRows))
		peak := 0.0
		for ri, row := range dataRows {
			if colIdx < len(row) {
				v, _ := strconv.ParseFloat(strings.TrimSpace(row[colIdx]), 64)
				v *= unitMultiplier
				vals[ri] = v
				if v > peak {
					peak = v
				}
			}
		}
		series = append(series, chart.Series{Name: name, Values: vals, Peak: peak})
	}

	return chart.WriteHTML(htmlPath, chart.Config{
		Title:   title,
		XLabels: xLabels,
		Series:  series,
	})
}
