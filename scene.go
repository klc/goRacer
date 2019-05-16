package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
	"time"
)

type scene struct {
	frameRate  int32
	background *sdl.Texture
	road       *road
	racer      *Racer
	cars       *Cars
}

func newScene(renderer *sdl.Renderer) (*scene, error) {

	var frameRate int32 = 60

	background, err := img.LoadTexture(renderer, "assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("load backgroun error : %v", err)
	}

	road, err := newRoad(renderer)
	if err != nil {
		return nil, fmt.Errorf("newroad error :%v", err)
	}
	road.speed = 1200 / frameRate

	racer, err := newRacer(renderer)
	if err != nil {
		return nil, fmt.Errorf("newracer error :%v", err)
	}

	cars, err := newCars(renderer)
	if err != nil {
		return nil, fmt.Errorf("newcars error :v", err)
	}
	cars.speed = 600 / frameRate

	return &scene{frameRate: frameRate, background: background, road: road, racer: racer, cars: cars}, nil
}

func (scene *scene) run(renderer *sdl.Renderer) error {
	errc := make(chan error)
	defer close(errc)
	events := make(chan sdl.Event)
	donec := make(chan bool)
	defer close(donec)

	go func() {
		for {
			select {
			case <-donec:
				close(events)
			default:
				events <- sdl.WaitEvent()
			}
		}
	}()

	wg := sync.WaitGroup{}
	for {
		select {
		case err := <-errc:
			return err
		case <-donec:
			return nil
		default:
			wg.Add(1)
			go func() {
				select {
				case e := <-events:
					if done := scene.handleEvent(e); done {
						go func() {
							donec <- done
						}()
					}
				default:
					scene.update()

					if scene.racer.crash {
						outro(renderer)
						scene.restart()
						time.Sleep(time.Second)
					}

					err := scene.paint(renderer)

					if err != nil {
						errc <- err
					}
				}
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

func (scene *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if event.GetType() == sdl.KEYDOWN {
			switch e.Keysym.Sym {
			case sdl.K_UP:
				scene.racer.newPosition(-1)
				break
			case sdl.K_DOWN:
				scene.racer.newPosition(1)
				break
			}
		}

	}
	return false
}

func (scene *scene) paint(renderer *sdl.Renderer) error {
	renderer.Clear()
	var err error

	sdl.Do(func() {
		err = scene.road.paint(renderer)
	})
	if err != nil {
		return fmt.Errorf("road paint error :%v", err)
	}

	sdl.Do(func() {
		err = scene.racer.paint(renderer)
	})
	if err != nil {
		return fmt.Errorf("racer paint error :%v", err)
	}

	sdl.Do(func() {
		err = scene.cars.paint(renderer)
	})
	if err != nil {
		return fmt.Errorf("cars paint error :%v", err)
	}

	sdl.Do(func() {
		renderer.Present()
		sdl.Delay(1000 / uint32(scene.frameRate))
	})

	return nil
}

func (scene *scene) destroy() {
	scene.road.destroy()
	scene.racer.destroy()
	scene.cars.destroy()
}

func (scene *scene) update() {
	scene.road.update()
	scene.cars.update()
	scene.cars.crash(scene.racer)
}

func (scene *scene) restart() {
	scene.road.restart()
	scene.racer.restart()
	scene.cars.restart()
}
