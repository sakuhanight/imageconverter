package converter

import (
	"go.uber.org/zap"
	"gopkg.in/gographics/imagick.v3/imagick"
	"path/filepath"
	"strconv"
)

type Option func(wand *imagick.MagickWand)

func Write(outputPath string, options ...Option) {
	// Wandを作成
	wand := imagick.NewMagickWand()

	// オプションを適用
	for _, option := range options {
		option(wand)
	}

	n := wand.GetNumberImages()
	zap.S().Debugf("number of images: %d", n)

	if n == 0 {
		zap.S().Error("no images")
		return
	} else if n == 1 {
		err := wand.WriteImage(outputPath)
		if err != nil {
			zap.S().Errorf("failed to write image: %v", err)
			return
		}
	} else if n > 1 {
		ext := filepath.Ext(outputPath)
		base := outputPath[:len(outputPath)-len(ext)]
		for i := 0; i < int(n); i++ {
			// 画像を書き込む
			ret := wand.SetIteratorIndex(i)
			if !ret {
				zap.S().Error("failed to set iterator index.")
				return
			}
			err := wand.WriteImage(base + "_" + strconv.Itoa(i) + ext)
			if err != nil {
				zap.S().Errorf("failed to write image: %v", err)
				return
			}
		}
	}
}

func ReadImage(path string) Option {
	return func(wand *imagick.MagickWand) {
		// 画像を読み込む
		err := wand.ReadImage(path)
		if err != nil {
			zap.S().Errorf("failed to read image: %v", err)
			return
		}
	}
}

// SetFormat ...
// フォーマットを設定する
func SetFormat(format string) Option {
	return func(wand *imagick.MagickWand) {
		err := wand.SetFormat(format)
		if err != nil {
			zap.S().Errorf("failed to set format: %v", err)
		}
		return
	}
}

// SetDPI ...
// dpiを設定する
func SetDPI(dpi float64) Option {
	return func(wand *imagick.MagickWand) {
		err := wand.SetResolution(dpi, dpi)
		if err != nil {
			zap.S().Errorf("failed to set dpi: %v", err)
		}
		return
	}
}

// SetHeight ...
// 高さを設定する
func SetSize(width, height uint) Option {
	return func(wand *imagick.MagickWand) {
		err := wand.ResizeImage(width, height, imagick.FILTER_LANCZOS)
		if err != nil {
			zap.S().Errorf("failed to set size: %v", err)
		}
		return
	}
}

// SetPassphrase ...
// パスフレーズを設定する
func SetPassphrase(passphrase string) Option {
	return func(wand *imagick.MagickWand) {
		err := wand.SetPassphrase(passphrase)
		if err != nil {
			zap.S().Errorf("failed to set passphrase: %v", err)
		}
		return
	}
}
