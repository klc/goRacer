package scene

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"time"

	"github.com/mkilic91/goRace/cars"
	"github.com/mkilic91/goRace/racer"
	"github.com/mkilic91/goRace/road"
)

type Scene struct {
	frameRate  int32
	background *sdl.Texture
	road       *road.Road
	racer      *racer.Racer
	cars       *cars.Cars
}

func NewScene(renderer *sdl.Renderer) (*Scene, error) {

	var frameRate int32 = 60

	background, err := img.LoadTexture(renderer, "assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("load backgroun error : %v", err)
	}

	newRoad, err := road.NewRoad(renderer)
	if err != nil {
		return nil, fmt.Errorf("newroad error :%v", err)
	}
	newRoad.SetSpeed(1200 / frameRate)

	newRacer, err := racer.NewRacer(renderer)
	if err != nil {
		return nil, fmt.Errorf("newracer error :%v", err)
	}

	newCars, err := cars.NewCars(renderer)
	if err != nil {
		return nil, fmt.Errorf("newcars error :v", err)
	}
	newCars.SetSpeed(600 / frameRate)

	return &Scene{frameRate: frameRate, background: background, road: newRoad, racer: newRacer, cars: newCars}, nil
}

func (scene *Scene) Run(events <-chan sdl.Event, renderer *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := scene.handleEvent(e); done {
					return
				}
			case <-tick:
				scene.update()

				if scene.racer.GetCrash() {
					//outro(renderer)
					time.Sleep(time.Second)
					scene.restart()
				}

				err := scene.paint(renderer)

				if err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (scene *Scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if event.GetType() == sdl.KEYDOWN {
			switch e.Keysym.Sym {
			case sdl.K_UP:
				scene.racer.NewPosition(-1)
				break
			case sdl.K_DOWN:
				scene.racer.NewPosition(1)
				break
			}
		}

	}
	return false
}

func (scene *Scene) paint(renderer *sdl.Renderer) error {
	renderer.Clear()
	var err error

	err = scene.road.Paint(renderer)
	if err != nil {
		return fmt.Errorf("road paint error :%v", err)
	}

	err = scene.racer.Paint(renderer)
	if err != nil {
		return fmt.Errorf("racer paint error :%v", err)
	}

	err = scene.cars.Paint(renderer)
	if err != nil {
		return fmt.Errorf("cars paint error :%v", err)
	}

	renderer.Present()

	return nil
}

func (scene *Scene) Destroy() {
	scene.road.Destroy()
	scene.racer.Destroy()
	scene.cars.Destroy()
}

func (scene *Scene) update() {
	scene.road.Update()
	scene.cars.Update()
	scene.cars.Crash(scene.racer)
}

func (scene *Scene) restart() {
	scene.road.Restart()
	scene.racer.Restart()
	scene.cars.Restart()
}
