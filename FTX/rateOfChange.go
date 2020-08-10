package FTX

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/goroumaru/test-code/utils"
	"github.com/montanaflynn/stats"
	"golang.org/x/xerrors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func analysis() {

	fileName := "BTC_USD BitMEX 過去データ1.csv"
	data, columnName := readData(fileName, false)
	fmt.Println(columnName)

	df, err := separateByColumn(data)
	if err != nil {
		fmt.Println(err)
	}
	date := getDateArray(data)
	changePrices, changeRatio, err := getChangePrice(df[1], df[2], true)
	if err != nil {
		fmt.Println(err)
	}

	targetData := stats.LoadRawData(changePrices)
	fmt.Println(targetData.StandardDeviationSample())

	// グラフ描画
	makeLineGraph("Daily Change Price", "date", "price", date, changePrices)
	makeLineGraph("Daily Change Ratio", "date", "price", date, changeRatio)
	makeHistgram("Histgram ~ Daily Change Price", "ratio", "frequency", changePrices)
	makeHistgram("Histgram ~ Daily Change Ratio", "ratio", "frequency", changeRatio)
}

// optABS : true = abs表示
func getChangePrice(openData, closeData []string, optABS bool) (changePrices, changeRatio []float64, err error) {
	if len(openData) != len(closeData) {
		return []float64{}, []float64{}, xerrors.New("close length don't equal open length.")
	}

	changePrices = make([]float64, 0)
	changeRatio = make([]float64, 0)
	for idx := 0; idx < len(openData); idx++ {
		close := utils.ToFloat64(strings.Replace(openData[idx], ",", "", -1)) // 数値に","が混ざっているため、削除する
		open := utils.ToFloat64(strings.Replace(closeData[idx], ",", "", -1))
		change := close - open // close-open
		if optABS {
			change = math.Abs(close - open) // abs(close-open)
		}
		ratio := change / close // change/close
		changePrices = append(changePrices, change)
		changeRatio = append(changeRatio, ratio)
	}
	return changePrices, changeRatio, err
}

// optColumn : ture = コラム抽出する
func readData(fileName string, optColumn bool) (data [][]string, columnName []string) {
	data = make([][]string, 0)

	csvFile, err := os.Open("./" + fileName)
	defer csvFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(csvFile)
	var line []string
	fmt.Println("--- START ---")
	for {
		line, err = reader.Read()
		if err == io.EOF {
			fmt.Println("--- END ---")
			break
		} else if err != nil {
			fmt.Println("--- ERROR END ---")
			fmt.Println(err)
			break
		}
		fmt.Println(line)
		data = append(data, line)
	}
	if optColumn {
		columnName = data[0] // columnだけ抽出
		data = append(data[1:len(data)], data[len(data):]...)
	}
	return data, columnName
}

// 列ごと配列とする
func separateByColumn(dataFlame [][]string) ([][]string, error) {
	if len(dataFlame) == 0 {
		return [][]string{}, xerrors.New("dataFlame is empty.")
	}
	colNum := len(dataFlame[0])
	rowNum := len(dataFlame)

	df := make([][]string, 0)
	for i := 0; i < colNum; i++ {
		col := make([]string, 0)
		for j := 0; j < rowNum; j++ {
			col = append(col, dataFlame[j][i])
		}
		df = append(df, col)
	}
	return df, nil
}

// UnixTimestamp配列を生成する
func getDateArray(dataFlame [][]string) (date []float64) {
	date = make([]float64, 0)
	layout := "2006年01月02日" //layout
	for _, data := range dataFlame {
		t, _ := time.Parse(layout, data[0])
		date = append(date, float64(t.Unix()))
	}
	return date
}

func makeLineGraph(title, xLabel, yLabel string, xData, yData []float64) {
	strRange := 0
	endRange := len(xData) - 1

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel
	p.Add(plotter.NewGrid())
	if err := plotutil.AddLinePoints(p, setXY(xData, yData, strRange, endRange)); err != nil {
		fmt.Println(err)
	}

	fileName := title + ".svg"
	if err := p.Save(5*vg.Inch, 5*vg.Inch, fileName); err != nil {
		fmt.Println(err)
	}
}

// start : start position of plot range
// end   : end posiiton of plot range
func setXY(x, y []float64, start, end int) plotter.XYs {
	// xデータなければ、整数列を生成する
	if len(x) == 0 {
		for i := 0; i < len(y); i++ {
			x = append(x, float64(i))
		}
	}
	pts := make(plotter.XYs, end-start+1)
	j := 0
	for i := start; i <= end; i++ {
		pts[j].X = x[i]
		pts[j].Y = y[i]
		j++
	}
	return pts
}

func makeHistgram(title, xLabel, yLabel string, data []float64) {
	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
	}
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel
	p.Legend.Top = true
	p.Add(plotter.NewGrid())
	h, err := plotter.NewHist(plotter.Values(data), 100)
	if err != nil {
		fmt.Println(err)
	}
	xmin, xmax, ymin, ymax := h.DataRange()
	fmt.Printf("Histgram Range (xmin, xmax, ymin, ymax) = (%v, %v, %v, %v)\n", xmin, xmax, ymin, ymax)
	// h.Normalize(ymax)
	p.Add(h)

	// standard devitation
	targetData := stats.LoadRawData(data)
	sd, err := targetData.StandardDeviationSample()
	if err != nil {
		fmt.Println(err)
	}
	setSigma(p, 1*sd, ymax)
	setSigma(p, 2*sd, ymax)
	setSigma(p, 3*sd, ymax)
	p.Legend.Add(fmt.Sprintf("1sigma:%v", utils.Round(1*sd, 3)))
	p.Legend.Add(fmt.Sprintf("2sigma:%v", utils.Round(2*sd, 3)))
	p.Legend.Add(fmt.Sprintf("3sigma:%v", utils.Round(3*sd, 3)))

	// ファイル出力
	fileName := title + ".svg"
	if err := p.Save(5*vg.Inch, 5*vg.Inch, fileName); err != nil {
		fmt.Println(err)
	}
}

func setSigma(plot *plot.Plot, sd, ymax float64) {
	lineSigma, err := plotter.NewLine(setXY([]float64{sd, sd}, []float64{0, ymax}, 0, 1))
	if err != nil {
		fmt.Println(err)
	}
	lineSigma.LineStyle.Color = color.RGBA{R: 191, G: 63, B: 63, A: 255}
	lineSigma.LineStyle.Dashes = []vg.Length{1, 2}
	plot.Add(lineSigma)
}
