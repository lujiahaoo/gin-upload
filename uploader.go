package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	maxFormSize = int64(27 << 20) 										//允许表单最大长度 27MiB

	maxImages = 9 														//允许最大上传图片数量

	supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".gif"}    //支持的图片类型

	distPath          = "./static"           							//普通图片存放根目录
	thumbnailDistPath = "./static/thumbnail" 							//缩略图片存放目录

	maxWidthThum  = uint(300)
	maxHeightThum = uint(200)

	timeout = time.Second * 5				//注意:这里的超时时间不包括图片从客户端发送到服务器所消耗的时间
)

func UploadImage(ctx *gin.Context) ([]string, error) {
	err := ctx.Request.ParseMultipartForm(maxFormSize)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	finishCh := make(chan bool)
	fileNameCh := make(chan string)
	var arr = []string{}

	fhs := ctx.Request.MultipartForm.File["image"]

	length := len(fhs)
	if length > maxImages {
		return nil, errors.New("too many images")
	}

	dayDir := path.Join(distPath, time.Now().Format("2006-01-02"))
	err = os.MkdirAll(dayDir, 0777)
	if err != nil {
		return nil, err
	}

	wg.Add(length)
	for _, fheader := range fhs {
		go saveUploadImage(dayDir, fheader, &wg, fileNameCh)
	}

	//任务全部完成，主动结束
	go func() {
		wg.Wait()
		finishCh <- true
	}()

	//监听任务主动结束或超时
	for {
		select {
		case <-time.After(timeout):
			return nil, errors.New("超时")
		case response := <-fileNameCh: //读取返回的文件名，到时候获取图片的话直接拼一下路径就好了，不用记录完整的路径
			arr = append(arr, response)
		case <-finishCh:
			return arr, nil
		}
	}

}

func isAllowImage(extName string) bool {
	for _, allowExt := range supportImageExtNames {
		if extName == allowExt {
			return true
		}
	}
	return false
}

func saveUploadImage(dayDir string, file *multipart.FileHeader, wg *sync.WaitGroup, fileNameCh chan<- string) {
	defer wg.Done()
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	extName := strings.ToLower(path.Ext(file.Filename))
	if isAllowImage(extName) == false {
		return
	}

	fileName := string(genRandomString(10)) + extName
	distPath := path.Join(dayDir, fileName)
	dist, err := os.Create(distPath)
	if err != nil {
		return
	}
	defer dist.Close()

	_, err = io.Copy(dist, src)
	if err != nil {
		return
	}

	fileNameCh <- fileName

	//缩略图片
	if err = thumbnailify(distPath); err != nil {

	}
}

func thumbnailify(imagePath string) error {
	var (
		file     *os.File
		img      image.Image
		fileName = path.Base(imagePath)
		extName  = strings.ToLower(path.Ext(imagePath))

		err error
	)

	outputPath := path.Join(thumbnailDistPath, fileName)

	//读取文件
	if file, err = os.Open(imagePath); err != nil {
		return err
	}
	defer file.Close()

	switch extName {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		break
	case ".png":
		img, err = png.Decode(file)
		break
	case ".gif":
		img, err = gif.Decode(file)
		break
	default:
		err = errors.New("不支持的类型" + extName)
		return err
	}

	if img == nil {
		err = errors.New("生成缩略图失败")
		return err
	}

	m := resize.Thumbnail(maxWidthThum, maxHeightThum, img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	switch extName {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, m, nil)
		break
	case ".png":
		png.Encode(out, m)
		break
	case ".gif":
		gif.Encode(out, m, nil)
		break
	default:
		err = errors.New("不支持的类型" + extName)
		return err
	}

	return nil
}


func GenRandomString(length int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return result
}
