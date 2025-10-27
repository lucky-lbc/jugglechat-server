package utils

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

// 头像结果结构体
type avatarResult struct {
	index int
	image image.Image
	err   error
}

// 生成群组头像（微信九宫格布局）
func GenerateGroupAvatar(avatarURLs []string, outputPath string) error {
	log.Printf("GenerateGroupAvatar: 接收到 %d 个头像URL", len(avatarURLs))

	// 重新配置参数
	const (
		canvasSize     = 300 // 画布尺寸
		gridSpacing    = 4   // 格子间距
		borderSize     = 8   // 整体边距
		maxConcurrency = 5   // 最大并发数
	)

	// 根据头像数量确定头像尺寸
	var avatarSize int
	switch len(avatarURLs) {
	case 1:
		avatarSize = 200 // 单头像大尺寸
	case 2:
		avatarSize = 140 // 两个头像中等尺寸
	case 3:
		// 3个头像需要计算合适的尺寸来铺满一行
		// 可用宽度 = 画布宽度 - 2*边距 - 间距
		availableWidth := canvasSize - 2*borderSize - gridSpacing
		avatarSize = availableWidth / 2 // 两个头像平分可用宽度
	case 4:
		availableWidth := canvasSize - 2*borderSize - gridSpacing
		avatarSize = availableWidth / 2 // 两个头像平分可用宽度
	default:
		avatarSize = 100 // 其他数量使用标准尺寸
	}

	// 创建白色背景
	canvas := imaging.New(canvasSize, canvasSize, color.RGBA{R: 240, G: 240, B: 240, A: 255})

	// 计算九宫格布局
	positions := calculateLayout(len(avatarURLs), canvasSize, avatarSize, gridSpacing, borderSize)

	// 限制最大处理数量
	maxAvatars := len(avatarURLs)
	if maxAvatars > 9 {
		maxAvatars = 9
	}
	log.Printf("GenerateGroupAvatar: 将处理 %d 个头像 (尺寸: %dpx)", maxAvatars, avatarSize)

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 并发下载并处理头像
	semaphore := make(chan struct{}, maxConcurrency) // 控制并发数
	resultsChan := make(chan avatarResult, maxAvatars)
	var wg sync.WaitGroup

	// 启动并发下载
	for i := 0; i < maxAvatars && i < len(positions); i++ {
		wg.Add(1)
		i := i // 捕获循环变量
		go func() {
			defer wg.Done()

			// 控制并发数
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			avatar, err := downloadAndProcessAvatar(avatarURLs[i], avatarSize, client)
			resultsChan <- avatarResult{
				index: i,
				image: avatar,
				err:   err,
			}
		}()
	}

	// 等待所有下载完成并关闭结果通道
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果并绘制到画布
	successCount := 0
	for avatarRes := range resultsChan {
		if avatarRes.err != nil {
			log.Printf("GenerateGroupAvatar: 处理头像 %d 失败: %v，跳过", avatarRes.index+1, avatarRes.err)
			continue
		}

		// 绘制到画布
		if avatarRes.image != nil && avatarRes.index < len(positions) {
			x, y := positions[avatarRes.index][0], positions[avatarRes.index][1]
			canvas = imaging.Paste(canvas, avatarRes.image, image.Pt(x, y))
			successCount++
		}
	}

	if successCount == 0 {
		log.Printf("GenerateGroupAvatar: 没有成功处理任何头像")
		return nil
	}

	// 保存结果
	return saveImage(canvas, outputPath)
}

// 计算九宫格布局位置
func calculateLayout(avatarCount, canvasSize, avatarSize, spacing, border int) [][]int {
	positions := [][]int{}

	// 根据头像数量计算布局
	switch avatarCount {
	case 1:
		// 单头像完全居中
		x := (canvasSize - avatarSize) / 2
		y := x
		positions = [][]int{{x, y}}

	case 2:
		// 水平排列，上下居中
		totalWidth := avatarSize*2 + spacing
		startX := (canvasSize - totalWidth) / 2
		y := (canvasSize - avatarSize) / 2
		positions = [][]int{
			{startX, y},
			{startX + avatarSize + spacing, y},
		}

	case 3:
		// 第一行：单个头像居中
		// 计算总高度：两行头像 + 行间距
		totalHeight := avatarSize*2 + spacing
		startY := (canvasSize - totalHeight) / 2

		// 第一行头像位置（居中）
		topX := (canvasSize - avatarSize) / 2
		positions = [][]int{{topX, startY}}

		// 第二行：两个头像整体居中，保持间距
		bottomY := startY + avatarSize + spacing

		// 计算第二行的总宽度（两个头像+间距）
		bottomTotalWidth := avatarSize*2 + spacing
		bottomStartX := (canvasSize - bottomTotalWidth) / 2

		positions = append(positions, []int{bottomStartX, bottomY})
		positions = append(positions, []int{bottomStartX + avatarSize + spacing, bottomY})

	case 4:
		startXY := border
		positions = [][]int{
			{startXY, startXY},
			{startXY + avatarSize + spacing, startXY},
			{startXY, startXY + avatarSize + spacing},
			{startXY + avatarSize + spacing, startXY + avatarSize + spacing},
		}
	case 5:
		// 上2下3，整体居中
		totalHeight := avatarSize*2 + spacing
		startY := (canvasSize - totalHeight) / 2

		topTotalWidth := avatarSize*2 + spacing
		topStartX := (canvasSize - topTotalWidth) / 2

		bottomTotalWidth := avatarSize*3 + spacing*2
		bottomStartX := (canvasSize - bottomTotalWidth) / 2

		positions = [][]int{
			{topStartX, startY},
			{topStartX + avatarSize + spacing, startY},
			{bottomStartX, startY + avatarSize + spacing},
			{bottomStartX + avatarSize + spacing, startY + avatarSize + spacing},
			{bottomStartX + (avatarSize+spacing)*2, startY + avatarSize + spacing},
		}

	case 6:
		// 上3下3，整体居中
		totalWidth := avatarSize*3 + spacing*2
		startX := (canvasSize - totalWidth) / 2
		totalHeight := avatarSize*2 + spacing
		startY := (canvasSize - totalHeight) / 2

		positions = [][]int{
			{startX, startY},
			{startX + avatarSize + spacing, startY},
			{startX + (avatarSize+spacing)*2, startY},
			{startX, startY + avatarSize + spacing},
			{startX + avatarSize + spacing, startY + avatarSize + spacing},
			{startX + (avatarSize+spacing)*2, startY + avatarSize + spacing},
		}

	case 7:
		// 上1中3下3，整体居中
		totalWidth := avatarSize*3 + spacing*2
		startX := (canvasSize - totalWidth) / 2
		totalHeight := avatarSize*3 + spacing*2
		startY := (canvasSize - totalHeight) / 2

		positions = [][]int{
			{startX + avatarSize + spacing, startY},
			{startX, startY + avatarSize + spacing},
			{startX + avatarSize + spacing, startY + avatarSize + spacing},
			{startX + (avatarSize+spacing)*2, startY + avatarSize + spacing},
			{startX, startY + (avatarSize+spacing)*2},
			{startX + avatarSize + spacing, startY + (avatarSize+spacing)*2},
			{startX + (avatarSize+spacing)*2, startY + (avatarSize+spacing)*2},
		}

	case 8:
		// 上2中3下3，整体居中
		// 上2居中、中3居中、下3居中
		totalHeight := avatarSize*3 + spacing*2
		startY := (canvasSize - totalHeight) / 2

		// 第一行：2个头像居中
		topTotalWidth := avatarSize*2 + spacing
		topStartX := (canvasSize - topTotalWidth) / 2

		// 第二行：3个头像居中
		middleTotalWidth := avatarSize*3 + spacing*2
		middleStartX := (canvasSize - middleTotalWidth) / 2

		// 第三行：3个头像居中
		bottomTotalWidth := avatarSize*3 + spacing*2
		bottomStartX := (canvasSize - bottomTotalWidth) / 2

		positions = [][]int{
			// 第一行（2个）
			{topStartX, startY},
			{topStartX + avatarSize + spacing, startY},
			// 第二行（3个）
			{middleStartX, startY + avatarSize + spacing},
			{middleStartX + avatarSize + spacing, startY + avatarSize + spacing},
			{middleStartX + (avatarSize+spacing)*2, startY + avatarSize + spacing},
			// 第三行（3个）
			{bottomStartX, startY + (avatarSize+spacing)*2},
			{bottomStartX + avatarSize + spacing, startY + (avatarSize+spacing)*2},
			{bottomStartX + (avatarSize+spacing)*2, startY + (avatarSize+spacing)*2},
		}

	case 9:
		// 标准3x3九宫格，整体居中
		totalSize := avatarSize*3 + spacing*2
		startXY := (canvasSize - totalSize) / 2

		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				x := startXY + col*(avatarSize+spacing)
				y := startXY + row*(avatarSize+spacing)
				positions = append(positions, []int{x, y})
			}
		}

	default: // 超过9个，按9个处理
		totalSize := avatarSize*3 + spacing*2
		startXY := (canvasSize - totalSize) / 2

		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if len(positions) >= 9 {
					break
				}
				x := startXY + col*(avatarSize+spacing)
				y := startXY + row*(avatarSize+spacing)
				positions = append(positions, []int{x, y})
			}
		}
	}

	return positions
}

// 下载并处理单个头像
func downloadAndProcessAvatar(url string, size int, client *http.Client) (image.Image, error) {
	if url == "" {
		return nil, nil
	}

	// 下载头像
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	// 解码头像
	srcImg, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// 调整大小并返回
	return imaging.Resize(srcImg, size, size, imaging.Lanczos), nil
}

// 保存图像
func saveImage(img image.Image, outputPath string) error {

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 以JPEG格式保存
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
}
