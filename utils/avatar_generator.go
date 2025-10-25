package utils

import (
	_ "bytes"
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

// 根据多个头像URL生成群组头像
func GenerateGroupAvatar(avatarURLs []string, outputPath string) error {
	log.Printf("GenerateGroupAvatar: 接收到 %d 个头像URL", len(avatarURLs))

	// 设置输出图像的尺寸
	const (
		size       = 300
		avatarSize = 100
	)

	// 创建白色背景的输出图像
	//log.Println("GenerateGroupAvatar: 创建300x300白色背景")
	result := imaging.New(size, size, color.White)

	// 计算头像位置（根据头像数量动态调整）
	var positions [][]int
	n := len(avatarURLs)
	switch n {
	case 1:
		positions = [][]int{{(size - avatarSize) / 2, (size - avatarSize) / 2}}
	case 2:
		positions = [][]int{
			{(size - avatarSize*2) / 3, (size - avatarSize) / 2},
			{(size-avatarSize*2)/3*2 + avatarSize, (size - avatarSize) / 2},
		}
	case 3:
		positions = [][]int{
			{(size - avatarSize) / 2, 0},
			{0, size - avatarSize},
			{size - avatarSize, size - avatarSize},
		}
	case 4:
		positions = [][]int{
			{0, 0}, {avatarSize, 0},
			{0, avatarSize}, {avatarSize, avatarSize},
		}
	case 5: // 上 2 下 3
		positions = [][]int{
			{(size - avatarSize*2) / 3, 0},
			{(size-avatarSize*2)/3*2 + avatarSize, 0},
			{0, size - avatarSize},
			{size/2 - avatarSize/2, size - avatarSize},
			{size - avatarSize, size - avatarSize},
		}
	case 6: // 上下两排 每排 3 个 整体居中 上下留空
		marginTop := (size - avatarSize*2) / 3
		marginLR := (size - avatarSize*3) / 4
		positions = [][]int{
			{marginLR, marginTop},
			{marginLR + avatarSize, marginTop},
			{marginLR + avatarSize*2, marginTop},
			{marginLR, marginTop + avatarSize},
			{marginLR + avatarSize, marginTop + avatarSize},
			{marginLR + avatarSize*2, marginTop + avatarSize},
		}
	case 7: // 第一排 1 个居中
		positions = [][]int{
			{(size - avatarSize) / 2, 0},
			{(size - avatarSize*3) / 4, avatarSize},
			{(size-avatarSize*3)/4 + avatarSize, avatarSize},
			{(size-avatarSize*3)/4 + avatarSize*2, avatarSize},
			{(size - avatarSize*3) / 4, avatarSize * 2},
			{(size-avatarSize*3)/4 + avatarSize, avatarSize * 2},
			{(size-avatarSize*3)/4 + avatarSize*2, avatarSize * 2},
		}
	case 8: // 第一排 2 个居中
		top2Left := (size - avatarSize*2) / 3
		positions = [][]int{
			{top2Left, 0},
			{top2Left + avatarSize, 0},
			{(size - avatarSize*3) / 4, avatarSize},
			{(size-avatarSize*3)/4 + avatarSize, avatarSize},
			{(size-avatarSize*3)/4 + avatarSize*2, avatarSize},
			{(size - avatarSize*3) / 4, avatarSize * 2},
			{(size-avatarSize*3)/4 + avatarSize, avatarSize * 2},
			{(size-avatarSize*3)/4 + avatarSize*2, avatarSize * 2},
		}
	default: // 9 人 3×3
		positions = [][]int{
			{0, 0}, {avatarSize, 0}, {avatarSize * 2, 0},
			{0, avatarSize}, {avatarSize, avatarSize}, {avatarSize * 2, avatarSize},
			{0, avatarSize * 2}, {avatarSize, avatarSize * 2}, {avatarSize * 2, avatarSize * 2},
		}
	}
	// 限制最大处理的头像数量
	maxAvatars := len(avatarURLs)
	if maxAvatars > 9 {
		maxAvatars = 9
	}
	log.Printf("GenerateGroupAvatar: 将处理 %d 个头像", maxAvatars)

	// 创建HTTP客户端，设置超时
	client := &http.Client{Timeout: 10 * time.Second}

	// 下载并处理每个头像
	successCount := 0
	for i := 0; i < maxAvatars && i < len(positions); i++ {
		url := avatarURLs[i]
		if url == "" {
			log.Printf("GenerateGroupAvatar: 跳过空URL")
			continue
		}

		log.Printf("GenerateGroupAvatar: 处理头像URL: %s", url)

		// 下载头像
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("GenerateGroupAvatar: 创建请求失败: %v，跳过此头像", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("GenerateGroupAvatar: 下载头像失败: %v，跳过此头像", err)
			continue
		}

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Printf("GenerateGroupAvatar: 下载头像返回非200状态码: %d，跳过此头像", resp.StatusCode)
			continue
		}

		// 解码头像
		srcImg, _, err := image.Decode(resp.Body)
		resp.Body.Close() // 及时关闭响应体

		if err != nil {
			log.Printf("GenerateGroupAvatar: 解码图像失败: %v，跳过此头像", err)
			continue
		}

		// 调整头像大小
		log.Printf("GenerateGroupAvatar: 调整头像 %d 大小为100x100", i+1)
		avatar := imaging.Resize(srcImg, avatarSize, avatarSize, imaging.Lanczos)

		// 绘制到结果图像上
		x, y := positions[i][0], positions[i][1]
		log.Printf("GenerateGroupAvatar: 将头像 %d 绘制到位置 (%d,%d)", i+1, x, y)
		result = imaging.Paste(result, avatar, image.Pt(x, y))
		successCount++
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("GenerateGroupAvatar: 创建目录失败: %v", err)
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("GenerateGroupAvatar: 创建输出文件失败: %v", err)
		return err
	}
	defer outFile.Close()

	// 以JPEG格式保存
	err = jpeg.Encode(outFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Printf("GenerateGroupAvatar: 编码JPEG失败: %v", err)
		return err
	}

	log.Println("GenerateGroupAvatar: 群组头像生成成功")
	return nil
}
