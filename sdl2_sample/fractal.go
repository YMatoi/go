package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"os"
	"strconv"
)

const WINDOW_W, WINDOW_H = 800, 600

type Point struct {
	X float64
	Y float64
}

type Line struct {
	P1 Point
	P2 Point
}

type Lines []Line

func DrawLine(renderer *sdl.Renderer, line Line) {
	renderer.DrawLine(int(line.P1.X), int(line.P1.Y), int(line.P2.X), int(line.P2.Y))
}

func DrawLines(renderer *sdl.Renderer, lines Lines) {
	for _, line := range lines {
		DrawLine(renderer, line)
	}
}

func (p1 Point) Minus(p2 Point) Point {
	return Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func (p1 Point) Plus(p2 Point) Point {
	return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func (p Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point) Norm() Point {
	length := p.Length()
	return Point{X: p.X / length, Y: p.Y / length}
}

func (p Point) Scale(s float64) Point {
	return Point{X: p.X * s, Y: p.Y * s}
}

func (p Point) Rotate(rad float64) Point {
	return Point{
		X: p.X*math.Cos(rad) - p.Y*math.Sin(rad),
		Y: p.X*math.Sin(rad) + p.Y*math.Cos(rad)}
}

func (p Point) Print(str string) {
	fmt.Printf("%v,%v,%v\n", str, p.X, p.Y)
}

func (line Line) Print(str string) {
	fmt.Printf("line,%v\n", str)
	line.P1.Print("p1")
	line.P2.Print("p2")
}

func (lines Lines) Print() {
	fmt.Printf("lines\n")
	for i, line := range lines {
		line.Print(strconv.Itoa(i))
	}
}

func (line Line) DragonNext() Lines {
	lines := Lines{}
	p := line.P2.Minus(line.P1)
	p = p.Rotate(-math.Pi / 4)
	p = p.Scale(math.Sqrt(2) / 2)
	p = line.P1.Plus(p)

	lines = append(lines, Line{P1: line.P1, P2: p})
	lines = append(lines, Line{P1: line.P2, P2: p})

	return lines
}

func (line Line) CNext() Lines {
	lines := Lines{}
	p := line.P2.Minus(line.P1)
	p = p.Rotate(-math.Pi / 4)
	p = p.Scale(math.Sqrt(2) / 2)
	p = line.P1.Plus(p)

	lines = append(lines, Line{P1: line.P1, P2: p})
	lines = append(lines, Line{P1: p, P2: line.P2})
	return lines
}

func (line Line) KochNext() Lines {
	lines := Lines{}
	p1 := line.P2.Minus(line.P1)
	p1 = p1.Scale(1.0 / 3.0)

	p3 := p1.Scale(2.0)

	p2 := p3.Minus(p1)
	p2 = p2.Rotate(-math.Pi / 3)

	p1 = line.P1.Plus(p1)
	p2 = p1.Plus(p2)
	p3 = line.P1.Plus(p3)

	lines = append(lines, Line{P1: line.P1, P2: p1})
	lines = append(lines, Line{P1: p1, P2: p2})
	lines = append(lines, Line{P1: p2, P2: p3})
	lines = append(lines, Line{P1: p3, P2: line.P2})

	return lines
}

func (line Line) Next(algorithm string) Lines {
	switch algorithm {
	case "C":
		return line.CNext()
	case "Dragon":
		return line.DragonNext()
	case "Koch":
		return line.KochNext()
	}
	return Lines{}
}

func (lines Lines) Next(algorithm string) Lines {
	ret := Lines{}
	for _, line := range lines {
		nextLines := line.Next(algorithm)
		ret = append(ret, nextLines...)
	}
	return ret
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WINDOW_W, WINDOW_H, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}

	lines := Lines{Line{P1: Point{X: 250, Y: 300}, P2: Point{X: 550, Y: 300}}}
	for i := 0; i < 10; i += 1 {
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)

		//lines.Print()

		DrawLines(renderer, lines)
		lines = lines.Next(os.Args[1])
		renderer.Present()
		sdl.Delay(1000)
	}

	sdl.Delay(1000)
	sdl.Quit()
}
