package cars

import (
	"fmt"
	"github.com/mkilic91/goRace/racer"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

type Cars struct {
	speed    int32
	cars     []*Car
	mutex    sync.RWMutex
	textures []*sdl.Texture
}

type Car struct {
	image    int32
	position int32
	mutex    sync.RWMutex
	x        int32
}

var images = []string{
	"assets/car-1.png",
	"assets/car-2.png",
	"assets/car-3.png",
}

var positions = []int32{50, 250, 450}

func NewCars(renderer *sdl.Renderer) (*Cars, error) {

	cars := &Cars{
		speed: 7,
	}

	for _, image := range images {
		texture, err := img.LoadTexture(renderer, image)

		if err != nil {
			return cars, fmt.Errorf("cars texture1 load error :%v", err)
		}

		cars.textures = append(cars.textures, texture)
	}

	go func() {
		for {
			cars.mutex.Lock()
			cars.cars = append(cars.cars, newCar())
			cars.mutex.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return cars, nil
}

func newCar() *Car {

	return &Car{
		image:    int32(rand.Intn(len(images))),
		position: int32(rand.Intn(len(positions))),
		x:        int32(1415),
	}
}

func (cars *Cars) Paint(renderer *sdl.Renderer) error {
	cars.mutex.RLock()
	defer cars.mutex.RUnlock()

	for _, car := range cars.cars {
		err := car.paint(renderer, cars.textures)

		if err != nil {
			return fmt.Errorf("cars render error :%v", err)
		}
	}

	return nil
}

func (car *Car) paint(renderer *sdl.Renderer, textures []*sdl.Texture) error {
	car.mutex.RLock()
	defer car.mutex.RUnlock()
	var err error

	rect := &sdl.Rect{X: car.x, Y: positions[car.position], W: 200, H: 100}

	err = renderer.CopyEx(textures[car.image], nil, rect, 0, nil, sdl.FLIP_NONE)

	if err != nil {
		return fmt.Errorf("car render error :%v", err)
	}

	return nil
}

func (cars *Cars) Update() {
	var temp []*Car

	for _, car := range cars.cars {
		car.mutex.Lock()
		car.x -= cars.speed
		car.mutex.Unlock()
		if car.x > -200 {
			temp = append(temp, car)
		}
	}

	cars.cars = temp

}

func (cars *Cars) Destroy() {
	for _, texture := range cars.textures {
		texture.Destroy()
	}
}

func (cars *Cars) Crash(racer *racer.Racer) {
	for _, car := range cars.cars {
		car.crash(racer)
	}
}

func (car *Car) crash(racer *racer.Racer) {
	var carxFront int32

	if racer.GetPosition() != racer.GetNextPosition() {
		carxFront = 150
	} else {
		carxFront = 250
	}

	if car.position == racer.GetPosition() && ((car.x > 50 && car.x < carxFront) || (car.x+200 > 50 && car.x+200 < 250)) {
		racer.SetCrash(true)
	}
}

func (cars *Cars) Restart() {
	cars.cars = []*Car{}
}

func (cars *Cars) SetSpeed(speed int32) {
	cars.speed = speed
}
