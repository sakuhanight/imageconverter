package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"imageconverter/constant"
	"imageconverter/converter"
	"imageconverter/docs"
	"os"
	"path/filepath"
	"strconv"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", Ping)
	r.POST("/convert", Convert)
	r.GET("/files", GetFileList)
	r.GET("/download", Download)
	r.GET("/admin/delete", Delete)

	docs.SwaggerInfo.BasePath = ""
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

// Ping godoc
// @Summary ping
// @Schemes
// @Description do ping
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Convert godoc
// @Summary convert file
// @Schemes
// @Description convert file
// @Accept multipart/form-data
// @Produce json
// @Param format query string true "出力ファイルフォーマット" Enums(pdf, png, gif, png8, webp, bmp, jpeg, jpg, tiff)
// @Param file formData file true "入力ファイル"
// @Param dpi query string false "dpi"
// @Param width query string false "width"
// @Param height query string false "height"
// @Param x query string false "x"
// @Param y query string false "y"
// @Param transformMode query string false "`width`もしくは`height`を指定した際の変形方法。指定無しの場合は`resize`として動作。" Enums(resize, crop)
// @Success 200 {object} string
// @Router /convert [post]
func Convert(c *gin.Context) {
	format := c.Query("format")
	file, _ := c.FormFile("file")
	id := uuid.NewString()
	filename := id + filepath.Ext(file.Filename)
	path := filepath.Join(constant.UPLOAD_FILE_PATH, filename)
	outputFilename := filename

	err := c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(500, gin.H{
			"id":      id,
			"message": "upload failed",
		})
		return
	}

	options := make([]converter.Option, 0)

	switch format {
	case "pdf":
		options = append(options, converter.SetFormat("pdf"))
		outputFilename = id + ".pdf"
	case "png":
		options = append(options, converter.SetFormat("png"))
		outputFilename = id + ".png"
	case "gif":
		options = append(options, converter.SetFormat("gif"))
		outputFilename = id + ".gif"
	case "png8":
		options = append(options, converter.SetFormat("png8"))
		outputFilename = id + ".png"
	case "webp":
		options = append(options, converter.SetFormat("webp"))
		outputFilename = id + ".webp"
	case "bmp":
		options = append(options, converter.SetFormat("bmp"))
		outputFilename = id + ".bmp"
	case "jpeg", "jpg":
		options = append(options, converter.SetFormat("jpg"))
		outputFilename = id + ".jpeg"
	case "tiff":
		options = append(options, converter.SetFormat("tiff"))
		outputFilename = id + ".tiff"
	default:
		c.JSON(400, gin.H{
			"message": "invalid format",
		})
		return
	}

	if c.Query("dpi") != "" {
		dpi, err := strconv.ParseFloat(c.Query("dpi"), 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid dpi",
			})
			return
		}
		options = append(options, converter.SetDPI(dpi))
	}

	options = append(options, converter.ReadImage(path))

	// サイズ変更
	if c.Query("width") != "" || c.Query("height") != "" {
		width, height := converter.GetSize(path)
		x, y := 0, 0
		switch c.DefaultQuery("transformMode", "resize") {
		case "resize":
			if c.Query("width") != "" {
				w, _ := strconv.Atoi(c.Query("width"))
				width = uint(w)
			}
			if c.Query("height") != "" {
				h, _ := strconv.Atoi(c.Query("height"))
				height = uint(h)
			}
			options = append(options, converter.SetSize(uint(width), uint(height)))
		case "crop":
			if c.Query("width") != "" {
				w, _ := strconv.Atoi(c.Query("width"))
				width = uint(w)
			}
			if c.Query("height") != "" {
				h, _ := strconv.Atoi(c.Query("height"))
				height = uint(h)
			}
			if c.Query("x") != "" {
				x, _ = strconv.Atoi(c.Query("x"))
			}
			if c.Query("y") != "" {
				y, _ = strconv.Atoi(c.Query("y"))
			}
			options = append(options, converter.Clop(width, height, x, y))
		}

	}

	go converter.Write(filepath.Join(constant.CONVERTED_FILE_PATH, outputFilename), options...)
	c.JSON(200, gin.H{
		"id":      id,
		"message": "accepted",
	})
	return
}

// GetFileList godoc
// @Summary get file list
// @Schemes
// @Description get file list
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} string
// @Router /files [get]
func GetFileList(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"message": "invalid id",
		})
		return
	}
	files := converter.GetFileList(id, constant.CONVERTED_FILE_PATH)
	c.JSON(200, gin.H{
		"files": files,
	})
}

// Download godoc
// @Summary download file
// @Schemes
// @Description download file
// @Accept json
// @Produce json
// @Param filename query string true "filename"
// @Success 200 {object} string
// @Router /download [get]
func Download(c *gin.Context) {
	filename := c.Query("filename")
	_, err := os.Open(filepath.Join(constant.CONVERTED_FILE_PATH, filename))
	if err != nil {
		c.JSON(404, gin.H{
			"message": "file not found",
		})
		return
	}
	// ファイルダウンロード
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath.Join(constant.CONVERTED_FILE_PATH, filename))
}

// Delete godoc
// @Summary delete file
// @Schemes
// @Description delete file
// @Accept json
// @Produce json
// @Param filename query string false "filename"
// @Param kind query string false "保存種類。指定無しの場合`converted`として動作。" Enums(converted, upload)
// @Param id query string false "id"
// @Success 200 {object} string
// @Router /admin/delete [get]
func Delete(c *gin.Context) {
	filename := c.Query("filename")
	id := c.Query("id")
	if filename == "" && id == "" {
		c.JSON(400, gin.H{
			"message": "invalid filename or id",
		})
		return
	} else if filename == "" && id != "" {
		if c.Query("kind") == "" || c.Query("kind") == "converted" {
			files := converter.GetFileList(id, constant.CONVERTED_FILE_PATH)
			for _, file := range files {
				err := os.Remove(filepath.Join(constant.CONVERTED_FILE_PATH, file))
				if err != nil {
					zap.S().Errorf("file (%s) delete failed: %+v", filepath.Join(constant.CONVERTED_FILE_PATH, file), err)
				}
			}
		}
		if c.Query("kind") == "" || c.Query("kind") == "upload" {
			files := converter.GetFileList(id, constant.UPLOAD_FILE_PATH)
			for _, file := range files {
				err := os.Remove(filepath.Join(constant.UPLOAD_FILE_PATH, file))
				if err != nil {
					zap.S().Errorf("file (%s) delete failed: %+v", filepath.Join(constant.UPLOAD_FILE_PATH, file), err)
				}
			}
		}
		c.JSON(200, gin.H{
			"message": "done",
		})
		return
	} else if filename != "" && id == "" {
		kind := c.DefaultQuery("kind", "converted")
		basepath := constant.CONVERTED_FILE_PATH
		switch kind {
		case "converted":
			basepath = constant.CONVERTED_FILE_PATH
		case "upload":
			basepath = constant.UPLOAD_FILE_PATH
		}
		err := os.Remove(filepath.Join(basepath, filename))
		if err != nil {
			zap.S().Errorf("file (%s) delete failed: %+v", filepath.Join(constant.CONVERTED_FILE_PATH, filename), err)
			c.JSON(500, gin.H{
				"message": "delete failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "done",
		})
	}
}
