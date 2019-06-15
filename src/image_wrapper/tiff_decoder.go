package image_wrapper

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
)

type decoder struct {
	r         io.ReaderAt     // указатель для чтения файла.
	byteOrder binary.ByteOrder// порядок байт
	config    image.Config    // итоговый конфиг для парсинга
	mode      imageMode       // режим изображения
	bpp       uint            // число бит на канал
	features  map[int][]uint  // список тегов
	palette   []color.Color

	buf   []byte
	off   int    // текущее смещение буфера
	v     uint32 // значение для буфера для считывания с произвольной глубиной
	nbits uint   // оставшееся число бит в буфере
}



// выпилиывает из буфера мусор(остаток строки)
func (d *decoder) flushBits() {
	d.v = 0
	d.nbits = 0
}

// считывает n бит с внутреннего буфера начиная с текущего индекса.
func (d *decoder) readBits(n uint) (v uint32, ok bool) {
	for d.nbits < n {
		d.v <<= 8
		if d.off >= len(d.buf) {
			return 0, false
		}
		d.v |= uint32(d.buf[d.off])
		d.off++
		d.nbits += 8
	}
	d.nbits -= n
	rv := d.v >> d.nbits
	d.v &^= rv << d.nbits
	return rv, true
}

// Получаем первое значение распарсеного тега заголовка по его ключу.
// Возвращает первое найденное значения с нужным тегом или 0.
func (d *decoder) firstVal(tag int) uint {
	f := d.features[tag]
	if len(f) == 0 {
		return 0
	}
	return f[0]
}

// считываем с буфера декодера  сырые данные по индексам.
// и пишем стрип или tile в dst изображения.
func (d *decoder) decode(dst image.Image,
	                     xmin, ymin, xmax, ymax int) error {
	d.off = 0

	// в спеках 64-65 кейс с обработкой разницы в случае с tPredictor равным prHorizontal
	if d.firstVal(tPredictor) == prHorizontal {
		switch d.bpp {
		case 8:
			var off int
			n := 1 * len(d.features[tBitsPerSample]) // bytes per sample times samples per pixel
			for y := ymin; y < ymax; y++ {
				off += n
				for x := 0; x < (xmax-xmin-1)*n; x++ {
					if off >= len(d.buf) {
						return errNoPixels
					}
					d.buf[off] += d.buf[off-n]
					off++
				}
			}
		case 1:
			return UnsupportedError("horizontal predictor with 1 BitsPerSample")
		}
	}

	rMaxX := minInt(xmax, dst.Bounds().Max.X)
	rMaxY := minInt(ymax, dst.Bounds().Max.Y)

	switch d.mode {
	case mRGB:
		if d.bpp == 8 {
			img := dst.(*image.RGBA)
			for y := ymin; y < rMaxY; y++ {
				// PixOffset returns the index of the first element of Pix that corresponds to the pixel at (x, y).
				min := img.PixOffset(xmin, y)
				max := img.PixOffset(rMaxX, y)
				off := (y - ymin) * (xmax - xmin) * 3
				for i := min; i < max; i += 4 {
					if off+3 > len(d.buf) {
						return errNoPixels
					}
					img.Pix[i+0] = d.buf[off+0]
					img.Pix[i+1] = d.buf[off+1]
					img.Pix[i+2] = d.buf[off+2]
					img.Pix[i+3] = 0xff
					off += 3
				}
			}
		}
	}

	return nil
}


// парсим ImageFolderDIrectory определяя поддерживаем этот тег или нет.
// Возвращаем номер тега и ошибку или nil.
func (d *decoder) parseIFD(p []byte) (int, error) {
	tag := d.byteOrder.Uint16(p[0:2])
	switch tag {
	case tBitsPerSample,
		tExtraSamples,
		tPhotometricInterpretation,
		tCompression,
		tPredictor,
		tStripOffsets,
		tStripByteCounts,
		tRowsPerStrip,
		tTileWidth,
		tTileLength,
		tTileOffsets,
		tTileByteCounts,
		tImageLength,
		tImageWidth:
		val, err := d.ifdUint(p)
		if err != nil {
			return 0, err
		}
		d.features[int(tag)] = val
	case tColorMap:
		val, err := d.ifdUint(p)
		if err != nil {
			return 0, err
		}
		numcolors := len(val) / 3
		if len(val)%3 != 0 || numcolors <= 0 || numcolors > 256 {
			return 0, FormatError("bad ColorMap length")
		}
		d.palette = make([]color.Color, numcolors)
		for i := 0; i < numcolors; i++ {
			d.palette[i] = color.RGBA64{
				uint16(val[i]),
				uint16(val[i+numcolors]),
				uint16(val[i+2*numcolors]),
				0xffff,
			}
		}
	case tSampleFormat:
		// Page 27 of the spec: If the SampleFormat is present and
		// the value is not 1 [= unsigned integer data], a Baseline
		// TIFF reader that cannot handle the SampleFormat value
		// must terminate the import process gracefully.
		val, err := d.ifdUint(p)
		if err != nil {
			return 0, err
		}
		for _, v := range val {
			if v != 1 {
				return 0, UnsupportedError("sample format")
			}
		}
	}
	return int(tag), nil
}


// Формируем тег и его описанию по последовательности байт IFD
// тип один из Byte, Short, Long type.
func (d *decoder) ifdUint(p []byte) (u []uint, err error) {
	var raw []byte
	if len(p) < ifdLen {
		return nil, FormatError("bad IFD entry")
	}

	datatype := d.byteOrder.Uint16(p[2:4])
	if dt := int(datatype); dt <= 0 || dt >= len(lengths) {
		return nil, UnsupportedError("IFD entry datatype")
	}

	count := d.byteOrder.Uint32(p[4:8])
	if count > math.MaxInt32/lengths[datatype] {
		return nil, FormatError("IFD data too large")
	}
	// Из спеков. стр 15
	//To save time and space the Value Offset contains the Value instead of pointing to
	//the Value if and only if the Value fits into 4 bytes. If the Value is shorter than 4
	//bytes, it is left-justified within the 4-byte Value Offset, i.e., stored in the lowernumbered bytes.
	// Whether the Value fits within 4 bytes is determined by the Type and Count of the field.
	if datalen := lengths[datatype] * count; datalen > 4 {
		// СОдержит указатель на значение
		// Согласно спекам, указатель может указывтаь на любое место в файле, даже после изображения.
		raw = make([]byte, datalen)
		_, err = d.r.ReadAt(raw, int64(d.byteOrder.Uint32(p[8:12])))
	} else {
		// Содержит значение
		raw = p[8 : 8+datalen]
	}
	if err != nil {
		return nil, err
	}

	u = make([]uint, count)
	switch datatype {
	case dtByte:
		for i := uint32(0); i < count; i++ {
			u[i] = uint(raw[i])
		}
	case dtShort:
		for i := uint32(0); i < count; i++ {
			u[i] = uint(d.byteOrder.Uint16(raw[2*i : 2*(i+1)]))
		}
	case dtLong:
		for i := uint32(0); i < count; i++ {
			u[i] = uint(d.byteOrder.Uint32(raw[4*i : 4*(i+1)]))
		}
	default:
		return nil, UnsupportedError("data type")
	}
	return u, nil
}


func newDecoder(r io.Reader) (*decoder, error) {
	// Чтение заголовка
	d := &decoder{
		r:        newReaderAt(r),
		features: make(map[int][]uint),
	}
	p := make([]byte, 8)
	if _, err := d.r.ReadAt(p, 0); err != nil {
		return nil, err
	}
	// Первые 4 байта - определение порядка байт для чтения
	switch string(p[0:4]) {
	case leHeader:
		// little-endian, от младшего к старшему, intel
		d.byteOrder = binary.LittleEndian
	case beHeader:
		// big-endian, от старшего к младшему, Motorola
		d.byteOrder = binary.BigEndian
	default:
		return nil, FormatError("malformed header")
	}
	// Следующие 4 байта - после определения порядка чтения может
	// смщение для начала считывания imageFileDirectory
	ifdOffset := int64(d.byteOrder.Uint32(p[4:8]))

	// Первые два байта начиная с ifdOffest - число изображений(директорий).
	if _, err := d.r.ReadAt(p[0:2], ifdOffset); err != nil {
		return nil, err
	}

	// Число изображений(директорий)
	numItems := int(d.byteOrder.Uint16(p[0:2]))
	// Игнорируем, будем считывать только первое
	// Создаём буфер по общему размеру описаний всех директорий
	p = make([]byte, ifdLen*numItems)
	if _, err := d.r.ReadAt(p, ifdOffset+2); err != nil {
		return nil, err
	}

	prevTag := -1
	for i := 0; i < len(p); i += ifdLen {
		// Считывание тега:
		// Скармливаем кусок буфера с описанием очередного тега
		// Если всё хорошо то пишем его по ключу в features.
		tag, err := d.parseIFD(p[i : i+ifdLen])
		if err != nil {
			return nil, err
		}
		if tag <= prevTag {
			// ПО спеке нужно проверить, что ключ тега больше либо равен предыдущему ключу
			return nil, FormatError("tags are not sorted in ascending order")
		}
		prevTag = tag
	}

	d.config.Width = int(d.firstVal(tImageWidth))
	d.config.Height = int(d.firstVal(tImageLength))

	// Проверяем, что в заголовке содержался тег число бит в канале
	if _, ok := d.features[tBitsPerSample]; !ok {
		return nil, FormatError("BitsPerSample tag missing")
	}
	d.bpp = d.firstVal(tBitsPerSample)
	switch d.bpp {
	case 0:
		return nil, FormatError("BitsPerSample must not be 0")
	case 1, 8, 16:
		log.Println("Bits per sample is fine")
		// Ожидаем что попадём именно сюда
	default:
		return nil, UnsupportedError(fmt.Sprintf("BitsPerSample of %v", d.bpp))
	}

	// Определяем режим изображения
	switch d.firstVal(tPhotometricInterpretation) {
	case pRGB:
		if d.bpp == 8 {
			for _, b := range d.features[tBitsPerSample] {
				if b != 8 {
					return nil, FormatError("wrong number of samples for 8bit RGB")
				}
			}
		}
		// 3 канала RGB каждый по 8 бит
		if len(d.features[tBitsPerSample]) == 3{
			d.mode = mRGB
			if d.bpp == 8 {
				log.Println("Color model parsed...: RGBA")
				d.config.ColorModel = color.RGBAModel
			}
		} else{
			return nil, FormatError("wrong number of samples for RGB")
		}
	default:
		return nil, UnsupportedError("color model")
	}
	return d, nil
}
