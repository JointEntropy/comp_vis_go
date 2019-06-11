package image_wrapper

import (
	"bufio"
	"fmt"
	"image"
	"io"
)

// Стандартный reader + интерфейс дл  заглядывания вперёд
type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

type format struct {
	name string
	decode       func(io.Reader) (image.Image, error)
	decodeConfig func(io.Reader) (image.Config, error)
}


// Основная функция парсинга файла типа tiff из указателя на чтение r.
// TODO: Возвращаёт объект изображения
func parse(r io.Reader) (img image.Image, err error) {
	d, err := newDecoder(r)
	if err != nil {
		return
	}

	blockPadding := false
	blockWidth := d.config.Width
	blockHeight := d.config.Height
	blocksAcross := 1
	blocksDown := 1

	/// TODO: Уточнить
	if d.config.Width == 0 {
		blocksAcross = 0
	}
	if d.config.Height == 0 {
		blocksDown = 0
	}

	var blockOffsets, blockCounts []uint

	if int(d.firstVal(tTileWidth)) != 0 {
		blockPadding = true

		blockWidth = int(d.firstVal(tTileWidth)) // число столбцов в блоке
		blockHeight = int(d.firstVal(tTileLength)) // число строк в блоке

		// Число блоков по ширине
		if blockWidth != 0 {
			blocksAcross = (d.config.Width + blockWidth - 1) / blockWidth
		}
		// Число блоков по высоте
		if blockHeight != 0 {
			blocksDown = (d.config.Height + blockHeight - 1) / blockHeight
		}

		blockCounts = d.features[tTileByteCounts]
		blockOffsets = d.features[tTileOffsets]

	} else {
		if int(d.firstVal(tRowsPerStrip)) != 0 {
			blockHeight = int(d.firstVal(tRowsPerStrip))
		}

		if blockHeight != 0 {
			blocksDown = (d.config.Height + blockHeight - 1) / blockHeight
		}

		blockOffsets = d.features[tStripOffsets]
		blockCounts = d.features[tStripByteCounts]
	}

	// Проверяем, что число блоков валидно. Check if we have the right number of strips/tiles, offsets and counts.
	if n := blocksAcross * blocksDown; len(blockOffsets) < n || len(blockCounts) < n {
		return nil, FormatError("inconsistent header")
	}

	imgRect := image.Rect(0, 0, d.config.Width, d.config.Height)
	// на основе d.mode создаём изображения некоторой размерности
	img = image.NewRGBA(imgRect)
	//matr := CreateMatrix(d.config.Height, d.config.Width, float64[][]{})

	// Для каждого строки блоков
	for i := 0; i < blocksAcross; i++ {
		blkW := blockWidth
		if !blockPadding && i == blocksAcross-1 && d.config.Width%blockWidth != 0 {
			blkW = d.config.Width % blockWidth
		}
		// Для каждого столбца блоков.
		for j := 0; j < blocksDown; j++ {
			blkH := blockHeight
			if !blockPadding && j == blocksDown-1 && d.config.Height%blockHeight != 0 {
				blkH = d.config.Height % blockHeight
			}
			offset := int64(blockOffsets[j*blocksAcross+i])
			n := int64(blockCounts[j*blocksAcross+i])
			switch d.firstVal(tCompression) {


			// Поддерживаем единственное значение режима сжатия- его отсутствие.
			case cNone, 0:
				if b, ok := d.r.(*buffer); ok {
					d.buf, err = b.Slice(int(offset), int(n)) // выбираем n байт начиная с offset и пишем в буфер
				} else {
					d.buf = make([]byte, n)
					_, err = d.r.ReadAt(d.buf, offset)
				}
			default:
				err = UnsupportedError(fmt.Sprintf("compression value %d", d.firstVal(tCompression)))
			}
			if err != nil {
				return nil, err
			}

			xmin := i * blockWidth
			ymin := j * blockHeight
			xmax := xmin + blkW
			ymax := ymin + blkH
			err = d.decode(img, xmin, ymin, xmax, ymax)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

// DecodeConfig returns the color model and dimensions of a TIFF image without decoding the entire image.
func parseConfig(r io.Reader) (image.Config, error) {
	d, err := newDecoder(r)
	if err != nil {
		return image.Config{}, err
	}
	return d.config, nil
}

