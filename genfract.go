package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"genfract/fractal"
)

var pixelFuns = map[string]func(x, y uint8) uint8{
	"X+Y":     func(x, y uint8) uint8 { return x + y },
	"X*Y":     func(x, y uint8) uint8 { return x * y },
	"(X+Y)/2": func(x, y uint8) uint8 { return (x + y) / 2 },
	"X^Y":     func(x, y uint8) uint8 { return x ^ y },
}

func main() {
	var mandelbrot = flag.Bool("mandelbrot", false, "flag to run Mandelbrot Set Viewer exercise")
	var julia = flag.Bool("julia", false, "flag to run Julia Set Viewer exercise")

	fmt.Println(os.Args)

	flag.Parse()

	http.Handle("/", fractal.HandleRoot())

	switch {

	case *mandelbrot:
		http.Handle("/mandelbrot", fractal.HandleMandelbrot())

	case *julia:
		http.Handle("/julia", fractal.HandleJulia())

	default:
		fmt.Println("...done!")
		return //exit
	}

	fmt.Println("Server running at 127.0.0.1:4000...")
	log.Fatal(http.ListenAndServe(":4000", nil))

}
