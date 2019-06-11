package image_wrapper


func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}


type FormatError string
func (e FormatError) Error() string {
	return "tiff: invalid format: " + string(e)
}

// An UnsupportedError reports that the input uses a valid but unimplemented feature.
type UnsupportedError string
func (e UnsupportedError) Error() string {
	return "tiff: unsupported feature: " + string(e)
}

var errNoPixels = FormatError("not enough pixel data")

