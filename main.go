package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/go-p5/p5"
)

type particle struct {
	x  float64
	y  float64
	vx float64
	vy float64
	c  color.Color
	cc int
}

type d struct {
	distx float64
	disty float64
	dist  float64
}

var (
	particles []particle        //see struct above
	k         float64    = 0.03 //spring constant of particle interaction
	kw        float64    = 0.4  //Spring coefficient of the walls
	ar        float64    = 0.01 //drag coefficient
	n         int        = 200  //number of particles
	minrange  float64    = 50   //inside of this range, particles repel one another
	distance  d                 //see struct above
	ca        = [5]color.Color{
		color.RGBA{R: 255, A: 255},
		color.RGBA{G: 255, A: 255},
		color.RGBA{B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 100, B: 255, A: 255}} //array of possible particle colors
	cm = [5][5]float64{
		{-0.02, -0.01, 0.01, -0.01, 0.01}, //r2r,r2g,r2b,r2w,r2p
		{0.01, -0.02, -0.01, 0.01, -0.01}, //g2r,g2g,g2b,g2w,g2p
		{-0.01, 0.01, -0.02, -0.01, 0.01}, //b2r,b2g,b2b,b2w,b2p
		{0.01, -0.01, 0.01, -0.02, -0.01}, //w2r,w2g,w2b,w2w,w2p
		{-0.01, 0.01, -0.01, 0.01, -0.02}} //p2r,p2g,p2b,p2w,p2p
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(1920, 1080)
	p5.Background(color.Gray{Y: 0})
	rand.Seed(time.Now().UnixNano())
	particles = Makeparticles(n)

}

func draw() {
	//draw dynamic particle
	for i := range particles {
		p := &particles[i]
		p5.Fill(p.c)
		p5.Stroke(p.c)
		p5.Circle(p.x, p.y, 10)
		////make interact with appropriate particles
		for j := range particles {
			op := &particles[j]
			d := &distance
			d.distx = math.Abs(p.x - op.x)
			d.disty = math.Abs(p.y - op.y)
			// calculate distance
			d.dist = math.Sqrt(math.Pow(d.distx, 2) + math.Pow(d.disty, 2))
			// Disqualify force under defined conditions
			if d.dist > 3*minrange || i == j || d.dist == 0 {
				continue
			}
			// choose proper force
			if d.dist < minrange { //universal repulsive force for small distances
				p.vx += k * (minrange - d.dist) * (p.x - op.x) / d.dist
				p.vy += k * (minrange - d.dist) * (p.y - op.y) / d.dist
			} else if d.dist >= minrange && d.dist <= 3*minrange { //color dependant force
				p.vx += cm[p.cc][op.cc] * (minrange - math.Abs(2*minrange-d.dist)) * (p.x - op.x) / d.dist
				p.vy += cm[p.cc][op.cc] * (minrange - math.Abs(2*minrange-d.dist)) * (p.y - op.y) / d.dist
			}
			//
		}
		////
		p.vx += -ar * math.Pow(p.vx, 2) * (math.Abs(p.vx) / p.vx)
		p.vy += -ar * math.Pow(p.vy, 2) * (math.Abs(p.vy) / p.vy)
		//reflecting off walls
		if p.y < 0 {
			p.vy += kw * (-p.y)
		} else if p.y > 1080 {
			p.vy += kw * (1080 - p.y)
		}
		if p.x < 0 {
			p.vx += kw * (-p.x)
		} else if p.x > 1920 {
			p.vx += kw * (1920 - p.x)
		}
		p.x += p.vx
		p.y += p.vy
	}
	//
}

func Newparticle(x float64, y float64, vx float64, vy float64, c color.Color, cc int) particle {
	return particle{
		x: x, y: y, vx: vx, vy: vy, c: c, cc: cc,
	}
}

func Makeparticles(n int) []particle {
	particles := make([]particle, n)
	for i := 0; i < n; i++ {
		r := rand.Intn(5)
		particles[i] = Newparticle(rand.Float64()*1920, rand.Float64()*1080, 2.5-rand.Float64()*5, 2.5-rand.Float64()*5, ca[r], r)
	}
	return particles
}
