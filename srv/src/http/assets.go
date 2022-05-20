package http

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

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"golang.org/x/image/draw"
)

func isImgResizable(id string) bool {
	switch strings.ToLower(filepath.Ext(id)) {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

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

func (a *api) renderPostAssetsIndexHandler() http.Handler {

	tpl := a.mustParseBasedTpl("assets.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ids, err := a.params.PostAssetStore.List()

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("getting list of asset ids: %w", err),
			)
			return
		}

		tplPayload := struct {
			IDs []string
		}{
			IDs: ids,
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}

func (a *api) getPostAssetHandler() http.Handler {

	renderIndexHandler := a.renderPostAssetsIndexHandler()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := filepath.Base(r.URL.Path)

		if id == "/" {
			renderIndexHandler.ServeHTTP(rw, r)
			return
		}

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

		if !isImgResizable(id) {
			apiutil.BadRequest(rw, r, fmt.Errorf("cannot resize file %q", id))
			return
		}

		if err := resizeImage(rw, buf, float64(maxWidth)); err != nil {
			apiutil.InternalServerError(
				rw, r,
				fmt.Errorf(
					"resizing image with id %q to size %d: %w",
					id, maxWidth, err,
				),
			)
		}

	})
}

func (a *api) postPostAssetHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := r.PostFormValue("id")
		if id == "/" {
			apiutil.BadRequest(rw, r, errors.New("id is required"))
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			apiutil.BadRequest(rw, r, fmt.Errorf("reading multipart file: %w", err))
			return
		}
		defer file.Close()

		if err := a.params.PostAssetStore.Set(id, file); err != nil {
			apiutil.InternalServerError(rw, r, fmt.Errorf("storing file: %w", err))
			return
		}

		a.executeRedirectTpl(rw, r, "assets/")
	})
}

func (a *api) deletePostAssetHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := filepath.Base(r.URL.Path)

		if id == "/" {
			apiutil.BadRequest(rw, r, errors.New("id is required"))
			return
		}

		err := a.params.PostAssetStore.Delete(id)

		if errors.Is(err, post.ErrAssetNotFound) {
			http.Error(rw, "Asset not found", 404)
			return
		} else if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("deleting asset with id %q: %w", id, err),
			)
			return
		}

		a.executeRedirectTpl(rw, r, "assets/")
	})
}
