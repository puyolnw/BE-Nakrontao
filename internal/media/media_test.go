package media

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestSanitizeJPEGDropsExif(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(1, 1, color.RGBA{200, 10, 10, 255})
	var plain bytes.Buffer
	if err := jpeg.Encode(&plain, img, nil); err != nil {
		t.Fatal(err)
	}
	// Splice an APP1 "Exif" segment with a fake GPS marker right after SOI.
	payload := []byte("Exif\x00\x00GPS-SECRET-13.75,100.50")
	seg := []byte{0xFF, 0xE1, 0, 0}
	binary.BigEndian.PutUint16(seg[2:], uint16(len(payload)+2))
	withExif := append(append(append([]byte{}, plain.Bytes()[:2]...), append(seg, payload...)...), plain.Bytes()[2:]...)

	out, err := Sanitize(".jpg", withExif)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(out, []byte("GPS-SECRET")) || bytes.Contains(out, []byte("Exif")) {
		t.Fatal("EXIF survived sanitizing")
	}
	if _, _, err := image.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("output is not a valid image: %v", err)
	}
}

func chunk(id string, data []byte) []byte {
	c := append([]byte(id), 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(c[4:], uint32(len(data)))
	c = append(c, data...)
	if len(data)%2 == 1 {
		c = append(c, 0)
	}
	return c
}

func TestStripWebP(t *testing.T) {
	vp8x := make([]byte, 10)
	vp8x[0] = 0x08 | 0x04
	body := append(chunk("VP8X", vp8x), chunk("VP8 ", []byte("pixels"))...)
	body = append(body, chunk("EXIF", []byte("GPS-SECRET"))...)
	body = append(body, chunk("XMP ", []byte("<author>x</author>"))...)
	file := append([]byte("RIFF"), 0, 0, 0, 0)
	file = append(file, "WEBP"...)
	file = append(file, body...)
	binary.LittleEndian.PutUint32(file[4:], uint32(len(file)-8))

	out, err := stripWebP(file)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(out, []byte("GPS-SECRET")) || bytes.Contains(out, []byte("author")) {
		t.Fatal("metadata survived")
	}
	if !bytes.Contains(out, []byte("pixels")) {
		t.Fatal("image data lost")
	}
	if out[20]&0x0C != 0 {
		t.Fatal("VP8X metadata flags not cleared")
	}
	if int(binary.LittleEndian.Uint32(out[4:])) != len(out)-8 {
		t.Fatal("RIFF size not updated")
	}
	if _, err := stripWebP([]byte("RIFF\x10\x00\x00\x00WEBPVP8 \xff\xff\xff\x7f")); err == nil {
		t.Fatal("truncated chunk accepted")
	}
}
