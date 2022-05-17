package api

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"golang.org/x/image/draw"
)

func resizeImage(out io.Writer, in io.Reader, maxWidth float64) error {

	img, format, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("decoding image: %w", err)
	}

	imgRect := img.Bounds()
	imgW, imgH := float64(imgRect.Dx()), float64(imgRect.Dy())

	if imgW > maxWidth {

		newH := imgH * maxWidth / imgW
		newImg := image.NewRGBA(image.Rect(0, 0, int(maxWidth), int(newH)))

		// Resize
		draw.BiLinear.Scale(
			newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil,
		)

		img = newImg
	}

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, nil)
	case "png":
		return png.Encode(out, img)
	default:
		return fmt.Errorf("unknown image format %q", format)
	}
}

func (a *api) servePostAssetHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := filepath.Base(r.URL.Path)

		maxWidth, err := apiutil.StrToInt(r.FormValue("w"), 0)
		if err != nil {
			apiutil.BadRequest(rw, r, fmt.Errorf("invalid w parameter: %w", err))
			return
		}

		buf := new(bytes.Buffer)

		err = a.params.PostAssetStore.Get(id, buf)

		if errors.Is(err, post.ErrAssetNotFound) {
			http.Error(rw, "Asset not found", 404)
			return
		} else if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching asset with id %q: %w", id, err),
			)
			return
		}

		if maxWidth == 0 {

			if _, err := io.Copy(rw, buf); err != nil {
				apiutil.InternalServerError(
					rw, r,
					fmt.Errorf(
						"copying asset with id %q to response writer: %w",
						id, err,
					),
				)
			}

			return
		}

		switch ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(id), ".")); ext {
		case "jpg", "jpeg", "png":

			if err := resizeImage(rw, buf, float64(maxWidth)); err != nil {
				apiutil.InternalServerError(
					rw, r,
					fmt.Errorf(
						"resizing image with id %q to size %d: %w",
						id, maxWidth, err,
					),
				)
			}

		default:
			apiutil.BadRequest(rw, r, fmt.Errorf("cannot resize file with extension %q", ext))
			return
		}

	})
}
