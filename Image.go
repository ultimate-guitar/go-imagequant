/*
Copyright (c) 2016, The go-imagequant author(s)

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND ISC DISCLAIMS ALL WARRANTIES WITH REGARD TO
THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS.
IN NO EVENT SHALL ISC BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR
CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS
ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS
SOFTWARE.
*/

package imagequant

import (
	"errors"
	"unsafe"
	"runtime"
)

/*
#include "libimagequant.h"
#include "stdlib.h"
*/
import "C"

type Image struct {
	p        *C.struct_liq_image
	dataP    unsafe.Pointer
	w, h     int
	released bool
}

// Creates an object that represents the image pixels to be used for quantization and remapping.
// The first argument is attributes object from NewAttributes(). The same attr object should be used for the entire process, from creation of images to quantization.
// The rgba32data string must be contiguous run of RGBA pixels (alpha is the last component, 0 = transparent, 255 = opaque).
// The Image.dataP that contain CString of rgba32data must not be modified or freed until this object is freed with Image.Release().
// Callers MUST call Release() on the returned object to free memory.
func NewImage(attr *Attributes, rgba32data string, width, height int, gamma float64) (*Image, error) {
	dataP := unsafe.Pointer(C.CString(rgba32data))
	pImg := C.liq_image_create_rgba(attr.p, dataP, C.int(width), C.int(height), C.double(gamma))
	if pImg == nil {
		C.free(dataP)
		return nil, errors.New("Failed to create image (invalid argument)")
	}
	img := &Image{
		p:        pImg,
		w:        width,
		dataP:    dataP,
		h:        height,
		released: false,
	}

	runtime.SetFinalizer(img, img.release)
	return img, nil
}

// Saved for backward capability. You should not call it.
func (this *Image) Release() {
	return
}

func (this *Image) release() {
	if !this.released{
		C.liq_image_destroy(this.p)
		C.free(this.dataP)
		this.released = true
	}
}

// Performs quantization (palette generation) based on settings in attr and pixels of the image.
func (this *Image) Quantize(attr *Attributes) (*Result, error) {
	res := Result{
		im: this,
	}
	runtime.SetFinalizer(res, res.release)
	liqerr := C.liq_image_quantize(this.p, attr.p, &res.p)
	if liqerr != C.LIQ_OK {
		return nil, translateError(liqerr)
	}

	return &res, nil
}
