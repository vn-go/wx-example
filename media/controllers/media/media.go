package media

import (
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/vn-go/wx"
	"github.com/wcharczuk/go-chart/v2"
)

type Media struct {
	wx.Handler
	RootDir string
}

func (media *Media) New() error {
	if media.RootDir == "" {
		media.RootDir = "./../uploads"
	}

	// Dùng MkdirAll để đảm bảo tạo cả cây thư mục nếu chưa có
	err := os.MkdirAll(filepath.Clean(media.RootDir), 0755)
	if err != nil {
		return err
	}

	return nil
}

// The function will return all files and folder
// in RootDir of Media
func (media *Media) ListAllFolderAndFiles() ([]string, error) {
	if media.RootDir == "" {
		return nil, nil
	}

	var results []string
	err := filepath.WalkDir(media.RootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // nếu có lỗi khi duyệt thì return ngay
		}
		results = append(results, path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return results, nil
}
func (media *Media) SaveFile(fileHeader multipart.FileHeader) (string, error) {
	files, err := media.ListAllFolderAndFiles()
	if err != nil {
		return "", err
	}
	if len(files) >= 2500 {
		return "Full", nil
	}

	// Mở file từ request
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open uploaded file: %w", err)
	}
	defer src.Close()

	// Lấy extension gốc của file
	ext := filepath.Ext(fileHeader.Filename)

	// Tạo tên file mới bằng UUID + extension
	newFileName := uuid.New().String() + ext

	// Tạo đường dẫn đầy đủ tới file trên server
	dstPath := filepath.Join(media.RootDir, newFileName)

	// Tạo file mới trên server
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file on server: %w", err)
	}
	defer dst.Close()

	// Copy nội dung từ file upload vào file server
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("cannot save file content: %w", err)
	}

	return newFileName, nil
}
func (media *Media) List(h *struct {
	wx.Handler `route:"method:get"`
}) (any, error) {
	return media.ListAllFolderAndFiles()
}

func (media *Media) Upload(ctx *wx.Handler, file multipart.FileHeader) error {
	newFileName, err := media.SaveFile(file)
	if err != nil {
		return err
	}
	(*ctx)().Res.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", newFileName)))
	return nil
}
func (medid *Media) LatencyChartHandler(w http.ResponseWriter) {
	chi := []float64{96.14, 89.22, 138.01, 163.33, 381.67}
	wx := []float64{85.43, 82.5, 116.32, 136.71, 405.12}
	labels := []string{"Avg", "Median(p50)", "p90", "p95", "Max"}

	graph := chart.BarChart{
		Title:      "Latency Comparison: WX vs CHI",
		TitleStyle: chart.Style{Hidden: false},
		Height:     400,
		BarWidth:   40,
		Bars: []chart.Value{
			{Value: chi[0], Label: labels[0], Style: chart.Style{FillColor: chart.ColorRed}},
			{Value: wx[0], Label: labels[0], Style: chart.Style{FillColor: chart.ColorBlue}},
			{Value: chi[1], Label: labels[1], Style: chart.Style{FillColor: chart.ColorRed}},
			{Value: wx[1], Label: labels[1], Style: chart.Style{FillColor: chart.ColorBlue}},
			{Value: chi[2], Label: labels[2], Style: chart.Style{FillColor: chart.ColorRed}},
			{Value: wx[2], Label: labels[2], Style: chart.Style{FillColor: chart.ColorBlue}},
			{Value: chi[3], Label: labels[3], Style: chart.Style{FillColor: chart.ColorRed}},
			{Value: wx[3], Label: labels[3], Style: chart.Style{FillColor: chart.ColorBlue}},
			{Value: chi[4], Label: labels[4], Style: chart.Style{FillColor: chart.ColorRed}},
			{Value: wx[4], Label: labels[4], Style: chart.Style{FillColor: chart.ColorBlue}},
		},
		Elements: []chart.Renderable{
			// Tạo legend thủ công
			chart.Legend(&chart.Chart{
				Series: []chart.Series{
					chart.ContinuousSeries{
						Name:  "CHI (Red)",
						Style: chart.Style{Hidden: false, FillColor: chart.ColorRed},
					},
					chart.ContinuousSeries{
						Name:  "WX (Blue)",
						Style: chart.Style{Hidden: false, FillColor: chart.ColorBlue},
					},
				},
			}),
		},
	}

	// Render trực tiếp ra HTTP response
	w.Header().Set("Content-Type", "image/png")
	_ = graph.Render(chart.PNG, w)
}
func (media *Media) LatencyLineChartHandler(w http.ResponseWriter) {
	// Dữ liệu
	labels := []string{"Avg", "Median(p50)", "p90", "p95", "Max"}
	chi := []float64{96.14, 89.22, 138.01, 163.33, 381.67}
	wx := []float64{85.43, 82.5, 116.32, 136.71, 405.12}

	// Tạo chart
	graph := chart.Chart{
		Title: "Latency Comparison (CHI vs WX)",
		XAxis: chart.XAxis{
			Name: "Metric",
			Style: chart.Style{
				Hidden: false,
			},
			ValueFormatter: func(v interface{}) string {
				i := int(v.(float64))
				if i >= 0 && i < len(labels) {
					return labels[i]
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name: "Latency (ms)",
			Style: chart.Style{
				Hidden: false,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "CHI (Red)",
				XValues: []float64{0, 1, 2, 3, 4},
				YValues: chi,
				Style: chart.Style{
					Hidden:      false,
					StrokeColor: chart.ColorRed,
					FillColor:   chart.ColorRed.WithAlpha(64),
				},
			},
			chart.ContinuousSeries{
				Name:    "WX (Blue)",
				XValues: []float64{0, 1, 2, 3, 4},
				YValues: wx,
				Style: chart.Style{
					Hidden:      false,
					StrokeColor: chart.ColorBlue,
					FillColor:   chart.ColorBlue.WithAlpha(64),
				},
			},
		},
	}

	// Thêm legend
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	// Render trực tiếp ra HTTP response
	w.Header().Set("Content-Type", "image/png")
	_ = graph.Render(chart.PNG, w)
}

func (media *Media) Chart(ctx *struct {
	wx.Handler `route:"@.png;method:get"`
}) {

	media.LatencyChartHandler(ctx.Handler().Res)
}
func (media *Media) LineChart(ctx *struct {
	wx.Handler `route:"@.png;method:get"`
}) {

	media.LatencyLineChartHandler(ctx.Handler().Res)
}
