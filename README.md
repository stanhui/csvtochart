# csvtochart

将带宽 CSV 数据转换为 Grafana 风格的离线交互式 HTML 图表。

Convert bandwidth CSV data into a Grafana-style offline interactive HTML chart.

## Features / 功能

- 生成完全离线的 HTML 图表（ECharts 内嵌，无需网络）
- 支持指定 CSV 中带宽数据的单位（bps / Kbps / Mbps / Gbps / Tbps）
- 深色/浅色主题一键切换
- 图例多选独显（Isolate）/ 隐藏（Hide）模式
- 时间轴缩放（滑块 + 鼠标滚轮）
- Y 轴自动换算显示单位（1000 进制）

---

- Generates fully offline HTML charts (ECharts embedded, no network needed)
- Supports specifying the bandwidth unit in CSV (bps / Kbps / Mbps / Gbps / Tbps)
- Dark/light theme toggle
- Legend multi-select isolate / hide modes
- Time-axis zoom (slider + mouse wheel)
- Y-axis auto-scales display unit (1000-based)

## Install / 安装

```bash
go install github.com/user/csvtochart@latest
```

Or build from source:

```bash
git clone https://github.com/user/csvtochart.git
cd csvtochart
go build -o csvtochart .
```

Cross-compile for macOS Apple Silicon:

```bash
GOOS=darwin GOARCH=arm64 go build -o csvtochart_darwin_arm64 .
```

## Usage / 用法

```bash
csvtochart [flags] <input.csv> [output.html]
```

### Examples / 示例

```bash
# 默认 bps 单位，输出同名 .html
csvtochart bandwidth.csv

# 指定 CSV 数据单位为 Mbps
csvtochart -unit Mbps bandwidth.csv

# 指定输出路径和标题
csvtochart -unit Gbps -title "CDN Bandwidth" data.csv report.html
```

### Flags / 参数

| Flag | Default | Description |
|------|---------|-------------|
| `-unit` | `bps` | CSV 中带宽数据的单位 / Unit of bandwidth values in CSV: `bps`, `Kbps`, `Mbps`, `Gbps`, `Tbps` |
| `-title` | 文件名 | 图表标题 / Chart title |

### Input CSV Format / 输入格式

第一列为时间戳，其余列为带宽数值，列名即为图表中的系列名称：

Column 0 = timestamps, remaining columns = bandwidth values. Column names become series names in the chart.

```csv
,domain1.com,domain2.com,domain3.com
2025-01-01 00:00:00,100.5,200.3,150.8
2025-01-01 00:05:00,110.2,190.1,160.5
```

> 第一列列名可以为空（如上例），数据行第一列为时间。

### Chart Interaction / 图表交互

- **☀ Light / 🌙 Dark** — 切换深色/浅色主题
- **● Isolate** — 单击右侧图例选中域名（可多选），选中线高亮，tooltip 只显示选中数据
- **✕ Hide** — 切换到隐藏模式，单击图例隐藏/显示指定线
- **Show All** — 恢复显示全部
- **底部滑块 / 鼠标滚轮** — 时间轴缩放

### Unit Conversion / 单位换算

Y 轴和 Tooltip 使用 **1000 进制**：

| Unit | = bps |
|------|-------|
| 1 Kbps | 1,000 bps |
| 1 Mbps | 1,000,000 bps |
| 1 Gbps | 1,000,000,000 bps |
| 1 Tbps | 1,000,000,000,000 bps |

## License

MIT
