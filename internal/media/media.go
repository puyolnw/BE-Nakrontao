// Package media stores admin-uploaded images on local disk (scope 2.6, 7.2).
package media

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/reclover/sinakhontao/internal/httpx"
)

const maxUpload = 8 << 20

var allowed = map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"}

type Handler struct{ Dir string }

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload+1024)
	f, _, err := r.FormFile("file")
	if err != nil {
		return httpx.BadRequest("file field is required (max 8 MB)")
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, maxUpload+1))
	if err != nil || len(raw) > maxUpload {
		return httpx.BadRequest("file field is required (max 8 MB)")
	}
	// Trust the bytes, not the client's declared type or filename.
	ext, ok := allowed[http.DetectContentType(raw)]
	if !ok {
		return httpx.BadRequest("only JPEG, PNG or WebP images are allowed")
	}
	clean, err := Sanitize(ext, raw)
	if err != nil {
		return httpx.BadRequest(err.Error())
	}
	name := make([]byte, 16)
	rand.Read(name)
	fname := hex.EncodeToString(name) + ext
	if err := os.MkdirAll(h.Dir, 0o750); err != nil {
		return err
	}
	out, err := os.OpenFile(filepath.Join(h.Dir, fname), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := out.Write(clean); err != nil {
		os.Remove(out.Name())
		return err
	}
	return httpx.JSON(w, http.StatusCreated, map[string]string{"url": "/uploads/" + fname})
}

const maxPixels = 25_000_000 // decoding is bounded so a small file cannot exhaust memory

// Sanitize removes metadata (EXIF GPS position, camera serial, author) before an image is published.
// JPEG and PNG are decoded and re-encoded, which also drops anything appended after the image data;
// WebP keeps its pixel chunks and loses EXIF/XMP.
func Sanitize(ext string, b []byte) ([]byte, error) {
	if ext == ".webp" {
		return stripWebP(b)
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(b))
	if err != nil {
		return nil, errors.New("the image could not be read")
	}
	if cfg.Width*cfg.Height > maxPixels {
		return nil, errors.New("image is too large (max 25 megapixels)")
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, errors.New("the image could not be read")
	}
	var out bytes.Buffer
	if ext == ".png" {
		err = png.Encode(&out, img)
	} else {
		err = jpeg.Encode(&out, img, &jpeg.Options{Quality: 88})
	}
	return out.Bytes(), err
}

// stripWebP rewrites a RIFF/WEBP file without its EXIF and XMP chunks and clears their VP8X flags.
func stripWebP(b []byte) ([]byte, error) {
	bad := errors.New("the image could not be read")
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WEBP" {
		return nil, bad
	}
	out := append([]byte("RIFF"), 0, 0, 0, 0, 'W', 'E', 'B', 'P')
	for i := 12; i < len(b); {
		if i+8 > len(b) {
			return nil, bad
		}
		id := string(b[i : i+4])
		size := int(binary.LittleEndian.Uint32(b[i+4 : i+8]))
		end := i + 8 + size + size%2
		if size < 0 || end > len(b) {
			return nil, bad
		}
		chunk := append([]byte(nil), b[i:end]...)
		switch id {
		case "EXIF", "XMP ":
		case "VP8X":
			if len(chunk) > 8 {
				chunk[8] &^= 0x08 | 0x04 // EXIF and XMP present flags
			}
			out = append(out, chunk...)
		default:
			out = append(out, chunk...)
		}
		i = end
	}
	binary.LittleEndian.PutUint32(out[4:8], uint32(len(out)-8))
	return out, nil
}

// Serve exposes uploaded files without directory listing.
func (h *Handler) Serve() http.Handler {
	fs := http.StripPrefix("/uploads/", http.FileServer(http.Dir(h.Dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/uploads/" || r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		fs.ServeHTTP(w, r)
	})
}
