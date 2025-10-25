package utils

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

// 生成群组头像（微信九宫格布局）
func GenerateGroupAvatar(avatarURLs []string, outputPath string) error {
	log.Printf("GenerateGroupAvatar: 接收到 %d 个头像URL", len(avatarURLs))

	// 重新配置参数 - 增大头像尺寸并重新计算所有布局
	const (
		canvasSize  = 300 // 画布尺寸
		avatarSize  = 100 // 增大每个头像尺寸 (从86增加到100)
		gridSpacing = 4   // 格子间距
		borderSize  = 8   // 整体边距
	)

	// 创建白色背景
	result := imaging.New(canvasSize, canvasSize, color.White)

	// 计算九宫格布局
	positions := calculateLayout(len(avatarURLs), canvasSize, avatarSize, gridSpacing, borderSize)

	// 限制最大处理数量
	maxAvatars := len(avatarURLs)
	if maxAvatars > 9 {
		maxAvatars = 9
	}
	log.Printf("GenerateGroupAvatar: 将处理 %d 个头像", maxAvatars)

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 下载并处理每个头像
	successCount := 0
	for i := 0; i < maxAvatars && i < len(positions); i++ {
		avatar, err := downloadAndProcessAvatar(avatarURLs[i], avatarSize, client)
		if err != nil {
			log.Printf("GenerateGroupAvatar: 处理头像 %d 失败: %v，跳过", i+1, err)
			continue
		}

		// 绘制到画布
		x, y := positions[i][0], positions[i][1]
		result = imaging.Paste(result, avatar, image.Pt(x, y))
		successCount++
	}

	if successCount == 0 {
		log.Printf("GenerateGroupAvatar: 没有成功处理任何头像")
		return nil
	}

	// 保存结果
	return saveImage(result, outputPath)
}

// 计算九宫格布局位置 - 完全重新计算以适应更大的头像尺寸
func calculateLayout(avatarCount, canvasSize, avatarSize, spacing, border int) [][]int {
	positions := [][]int{}

	// 可用区域（减去边距）
	//availableSize := canvasSize - 2*border

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
		// 上1下2布局
		// 计算总高度
		totalHeight := avatarSize*2 + spacing
		startY := (canvasSize - totalHeight) / 2

		// 上面单个头像：水平居中
		topX := (canvasSize - avatarSize) / 2

		// 下面两个头像：水平居中排列
		bottomTotalWidth := avatarSize*2 + spacing
		bottomStartX := (canvasSize - bottomTotalWidth) / 2

		positions = [][]int{
			{topX, startY}, // 上面1个
			{bottomStartX, startY + avatarSize + spacing},                        // 下面左
			{bottomStartX + avatarSize + spacing, startY + avatarSize + spacing}, // 下面右
		}

	case 4:
		// 2x2网格，整体居中
		totalSize := avatarSize*2 + spacing
		startXY := (canvasSize - totalSize) / 2
		positions = [][]int{
			{startXY, startXY},
			{startXY + avatarSize + spacing, startXY},
			{startXY, startXY + avatarSize + spacing},
			{startXY + avatarSize + spacing, startXY + avatarSize + spacing},
		}

	case 5:
		// 上2下3，整体居中
		// 计算总高度和宽度
		totalHeight := avatarSize*2 + spacing
		startY := (canvasSize - totalHeight) / 2

		// 上排2个居中
		topTotalWidth := avatarSize*2 + spacing
		topStartX := (canvasSize - topTotalWidth) / 2

		// 下排3个居中
		bottomTotalWidth := avatarSize*3 + spacing*2
		bottomStartX := (canvasSize - bottomTotalWidth) / 2

		positions = [][]int{
			// 上排2个
			{topStartX, startY},
			{topStartX + avatarSize + spacing, startY},
			// 下排3个
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
			// 上排1个（居中）
			{startX + avatarSize + spacing, startY},
			// 中排3个
			{startX, startY + avatarSize + spacing},
			{startX + avatarSize + spacing, startY + avatarSize + spacing},
			{startX + (avatarSize+spacing)*2, startY + avatarSize + spacing},
			// 下排3个
			{startX, startY + (avatarSize+spacing)*2},
			{startX + avatarSize + spacing, startY + (avatarSize+spacing)*2},
			{startX + (avatarSize+spacing)*2, startY + (avatarSize+spacing)*2},
		}

	case 8:
		// 上3中2下3，整体居中
		totalWidth := avatarSize*3 + spacing*2
		startX := (canvasSize - totalWidth) / 2
		totalHeight := avatarSize*3 + spacing*2
		startY := (canvasSize - totalHeight) / 2

		// 中排2个居中
		middleTotalWidth := avatarSize*2 + spacing
		middleStartX := (canvasSize - middleTotalWidth) / 2

		positions = [][]int{
			// 上排3个
			{startX, startY},
			{startX + avatarSize + spacing, startY},
			{startX + (avatarSize+spacing)*2, startY},
			// 中排2个
			{middleStartX, startY + avatarSize + spacing},
			{middleStartX + avatarSize + spacing, startY + avatarSize + spacing},
			// 下排3个
			{startX, startY + (avatarSize+spacing)*2},
			{startX + avatarSize + spacing, startY + (avatarSize+spacing)*2},
			{startX + (avatarSize+spacing)*2, startY + (avatarSize+spacing)*2},
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
		// 标准3x3九宫格，整体居中
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
	// 确保输出目录存在
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
