// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fractal

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"strings"
)

var (
	extra = map[string]http.Handler{
		"/":             Serve("mandelbrot.html"),
		"/juliaviewer":  Serve("julia.html"),
		"/newtonviewer": Serve("newton.html"),
		"/behavior.js":  Serve("behavior.js"),
		"/center.gif":   Serve("center.gif"),
		"/gsv.js":       Serve("gsv.js"),
		"/none.png":     Serve("none.png"),
	}
)

// MainPage serves the main fractal viewing HTML page.
//var MainPage http.Handler = http.HandlerFunc(mainPage)

func rootPage(w http.ResponseWriter, req *http.Request) {
	log.Printf("request:%v\n", req.URL)
	h := extra[req.URL.Path]
	if h != nil {
		h.ServeHTTP(w, req)
		return
	}
	http.NotFound(w, req)
	w.WriteHeader(http.StatusNotFound)
}

func HandleRoot() http.Handler {
	return http.HandlerFunc(rootPage)
}

func HandleMandelbrot() http.Handler {
	return MyImage{requestParserFunc: mandelbrotReqParser}
}

func HandleJulia() http.Handler {
	return MyImage{requestParserFunc: juliaReqParser}
}

// Structs and methods for Mandelbrot and Julia Set viewer exercise -- Start

type Colorizer interface {
	At(x, y int) color.Color
}

type mandelbrotColorizer struct {
	props ColorProps
}

type juliaColorizer struct {
	props ColorProps
}

type ColorProps struct {
	width, height, iterations int
	origin, crange, c         complex128
}

type MyImage struct {
	requestParserFunc func(string) (ColorProps, Colorizer, error)
	colorProps        ColorProps
	Colorizer
}

func (m MyImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (m MyImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.colorProps.width, m.colorProps.height)
}

func (m MyImage) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	colorProps, colorizer, _ := m.requestParserFunc(req.FormValue("p"))
	m.colorProps = colorProps
	m.Colorizer = colorizer

	var buf bytes.Buffer
	png.Encode(&buf, m)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-rau-pod", envPod)
	w.Header().Set("x-rau-podip", envPodIP)
	w.Header().Set("X-rau-svc", envPodSvc)
	w.Write(buf.Bytes())
}

func mandelbrotReqParser(p string) (ColorProps, Colorizer, error) {
	h, _ := os.Hostname()
	log.Printf("%s: params: %v\n", h, p)
	var (
		x, y, iterations int
		origin, crange   complex128
	)

	r := strings.NewReader(p)
	_, err := fmt.Fscanf(r, "%d %d %g %g %d", &x, &y, &origin, &crange, &iterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fscanf: %v\n", err)
		return ColorProps{}, mandelbrotColorizer{}, err
	}

	return ColorProps{x, y, iterations, origin, crange, 0},
		mandelbrotColorizer{
			ColorProps{x, y, iterations, origin, crange, 0}}, nil
}

func (m mandelbrotColorizer) At(x, y int) color.Color {
	fx := float64(x) / float64(m.props.width)
	fy := float64(y) / float64(m.props.height)
	c := m.props.origin + complex(real(m.props.crange)*fx, imag(m.props.crange)*fy)

	return mandelbrotColor(c, m.props.iterations)
}

func mandelbrotColor(c complex128, iterations int) color.Color {
	const contrast = 15
	var z complex128
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			//return color.Gray{255 - uint8(contrast*i)}
			return Cycle(i, iterations)
			//return Ramp(i, iterations)
		}
	}

	return color.Black
}

func juliaReqParser(p string) (ColorProps, Colorizer, error) {
	var (
		x, y, iterations  int
		origin, crange, c complex128
	)

	r := strings.NewReader(p)
	_, err := fmt.Fscanf(r, "%d %d %g %g %g %d", &x, &y, &origin, &crange, &c, &iterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fscanf: %v, parsing the input : %s\n", err, p)
		return ColorProps{}, juliaColorizer{}, err
	}

	return ColorProps{x, y, iterations, origin, crange, c},
		juliaColorizer{ColorProps{x, y, iterations, origin, crange, c}}, nil
}

func (m juliaColorizer) At(x, y int) color.Color {
	fx := float64(x) / float64(m.props.width)
	fy := float64(y) / float64(m.props.height)
	z := m.props.origin + complex(real(m.props.crange)*fx, imag(m.props.crange)*fy)

	return juliaColor(z, m.props.c, m.props.iterations)
}

func juliaColor(z, c complex128, iterations int) color.Color {
	const contrast = 15
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			//return color.Gray{255 - uint8(contrast*i)}
			//return fractal.Cycle(i, iterations)
			return Ramp(i, iterations)
		}
	}
	return color.Black
}
