package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type racer struct {
	texture      *sdl.Texture
	position     int32
	nextPosition int32
	positions    []int32
	mutex        sync.RWMutex
	crash        bool
	music        *mix.Music
	y            int32
}

func newRacer(rederer *sdl.Renderer) (*racer, error) {
	texture, err := img.LoadTexture(rederer, "assets/racer.png")

	if err != nil {
		return nil, fmt.Errorf("racer texture load error : %v", err)
	}

	positions := []int32{50, 250, 450}

	return &racer{texture: texture, position: 1, nextPosition: 1, positions: positions, crash: false}, nil
}

func (racer *racer) paint(renderer *sdl.Renderer) error {
	racer.mutex.RLock()
	defer racer.mutex.RUnlock()

	rect := &sdl.Rect{X: 50, Y: racer.getPosition(), W: 200, H: 100}
	err := renderer.CopyEx(racer.texture, nil, rect, 0, nil, sdl.FLIP_NONE)

	if err != nil {
		return fmt.Errorf("racer render error : %v", err)
	}

	return nil
}

func (racer *racer) destroy() {
	racer.texture.Destroy()
}

func (racer *racer) newPosition(move int32) {
	newPosition := racer.position + move

	if newPosition <= 0 {
		racer.nextPosition = 0
	} else if int(newPosition) >= (len(racer.positions) - 1) {
		racer.nextPosition = int32(len(racer.positions)) - 1
	} else {
		racer.nextPosition = newPosition
	}
}

func (racer *racer) getPosition() int32 {

	if racer.position == racer.nextPosition {
		racer.y = racer.positions[racer.position]

		return racer.y
	}

	if racer.positions[racer.position] > racer.positions[racer.nextPosition] {
		if racer.y == racer.positions[racer.nextPosition] {
			racer.position = racer.nextPosition
		} else {
			racer.y -= 20
		}
	} else {
		if racer.y == racer.positions[racer.nextPosition] {
			racer.position = racer.nextPosition
		} else {
			racer.y += 20
		}
	}

	return racer.y
}

func (racer *racer) restart() {
	racer.crash = false
	racer.position = 1
	racer.nextPosition = 1
}
