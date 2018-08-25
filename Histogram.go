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

/*
#include "libimagequant.h"
*/
import "C"
import "runtime"

type Histogram struct {
	p *C.struct_liq_histogram
	released bool
}

// "Learns" colors from the image, which will be later used to generate the palette.
// After the image is added to the histogram it may be freed to save memory (but it's more efficient to keep the image object if it's going to be used for remapping).
// Fixed colors added to the image are also added to the histogram.
func (this *Histogram) AddImage(attr *Attributes, img *Image) error {
	return translateError(C.liq_histogram_add_image(this.p, attr.p, img.p))
}

// Quantize generate palette from the histogram.
func (this *Histogram) Quantize(attr *Attributes) (*Result, error) {
	res := Result{}
	runtime.SetFinalizer(res, res.release)
	liqerr := C.liq_histogram_quantize(this.p, attr.p, &res.p)
	if liqerr != C.LIQ_OK {
		return nil, translateError(liqerr)
	}

	return &res, nil
}

// Saved for backward capability. You should not call it.
func (this *Histogram) Release() {
	return
}

func (this *Histogram) release() {
	if !this.released {
		C.liq_histogram_destroy(this.p)
		this.released = true
	}
}