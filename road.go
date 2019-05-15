package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type road struct {
	texture *sdl.Texture
	speed   int32
	mutex   sync.RWMutex

	x int32
}

func newRoad(renderer *sdl.Renderer) (*road, error) {
	texture, err := img.LoadTexture(renderer, "assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("load backgroun error : %v", err)
	}

	return &road{texture: texture, speed: 15, x: 0}, nil
}

func (road *road) paint(renderer *sdl.Renderer) error {
	road.mutex.RLock()
	defer road.mutex.RUnlock()

	rect := &sdl.Rect{X: (0 - road.x), Y: 0, W: 1600, H: 600}

	err := renderer.CopyEx(road.texture, nil, rect, 0, nil, sdl.FLIP_NONE)

	if err != nil {
		return fmt.Errorf("road render error :%v", err)
	}

	return nil
}

func (road *road) update() {
	if road.x >= 195 {
		road.x = 0
	} else {
		road.x += road.speed
	}
}

func (road *road) destroy() {
	road.texture.Destroy()
}

func (road *road) restart() {
	road.x = 0
}
