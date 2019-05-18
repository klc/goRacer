package racer

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type Racer struct {
	texture      *sdl.Texture
	position     int32
	nextPosition int32
	positions    []int32
	mutex        sync.RWMutex
	crash        bool
	music        *mix.Music
	y            int32
	angle        float64
}

func NewRacer(rederer *sdl.Renderer) (*Racer, error) {
	texture, err := img.LoadTexture(rederer, "assets/racer.png")

	if err != nil {
		return nil, fmt.Errorf("racer texture load error : %v", err)
	}

	positions := []int32{50, 250, 450}

	return &Racer{texture: texture, position: 1, nextPosition: 1, angle: 0, positions: positions, crash: false}, nil
}

func (racer *Racer) Paint(renderer *sdl.Renderer) error {
	racer.mutex.RLock()
	defer racer.mutex.RUnlock()

	rect := &sdl.Rect{X: 50, Y: racer.getPosition(), W: 200, H: 100}
	err := renderer.CopyEx(racer.texture, nil, rect, racer.angle, nil, sdl.FLIP_NONE)

	if err != nil {
		return fmt.Errorf("racer render error : %v", err)
	}

	return nil
}

func (racer *Racer) Destroy() {
	racer.texture.Destroy()
}

func (racer *Racer) NewPosition(move int32) {
	newPosition := racer.position + move

	if newPosition <= 0 {
		racer.nextPosition = 0
	} else if int(newPosition) >= (len(racer.positions) - 1) {
		racer.nextPosition = int32(len(racer.positions)) - 1
	} else {
		racer.nextPosition = newPosition
	}
}

func (racer *Racer) getPosition() int32 {

	if racer.position == racer.nextPosition {
		racer.y = racer.positions[racer.position]
		racer.angle = 0
		return racer.y
	}

	if racer.positions[racer.position] > racer.positions[racer.nextPosition] {
		if racer.y == racer.positions[racer.nextPosition] {
			racer.position = racer.nextPosition
		} else {
			racer.y -= 20
			racer.angle = -20
		}
	} else {
		if racer.y == racer.positions[racer.nextPosition] {
			racer.position = racer.nextPosition
		} else {
			racer.y += 20
			racer.angle = 20
		}
	}

	return racer.y
}

func (racer *Racer) Restart() {
	racer.crash = false
	racer.position = 1
	racer.nextPosition = 1
	racer.angle = 0
}

func (racer *Racer) GetPosition() int32 {
	return racer.position
}

func (racer *Racer) GetNextPosition() int32 {
	return racer.nextPosition
}

func (racer *Racer) SetCrash(crash bool) {
	racer.crash = crash
}

func (racer *Racer) GetCrash() bool {
	return racer.crash
}
