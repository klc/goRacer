package main

import (
	"fmt"
	"github.com/mkilic91/goRace/scene"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"runtime"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	var err error

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	/*	err = mix.Init(mix.INIT_MP3)
		if err != nil {
			return fmt.Errorf("could not initialize MIX: %v", err)
		}
		mix.Quit()

		err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
		if err != nil {
			return fmt.Errorf("mix audio error :%v", err)
		}
		mix.CloseAudio()*/

	var window *sdl.Window
	var renderer *sdl.Renderer

	window, renderer, err = sdl.CreateWindowAndRenderer(1415, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer window.Destroy()

	err = intro(renderer)
	if err != nil {
		return fmt.Errorf("intro error :%v", err)
	}

	var newScene *scene.Scene

	newScene, err = scene.NewScene(renderer)
	if err != nil {
		return fmt.Errorf("scene create error :%v", err)
	}
	defer newScene.Destroy()

	var events chan sdl.Event
	events = make(chan sdl.Event)
	errc := newScene.Run(events, renderer)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}

}

func drawTitle(renderer *sdl.Renderer, text string) error {
	renderer.Clear()

	font, err := ttf.OpenFont("assets/race.ttf", 500)

	if err != nil {
		return fmt.Errorf("font load error : %v", err)
	}
	defer font.Close()

	color := sdl.Color{R: 255, G: 40, B: 0, A: 255}
	surface, err := font.RenderUTF8Solid(text, color)

	if err != nil {
		return fmt.Errorf("font render error : %v", err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)

	if err != nil {
		return fmt.Errorf("create texture error : %v", err)
	}
	defer texture.Destroy()

	err = renderer.Copy(texture, nil, nil)

	if err != nil {
		return fmt.Errorf("render copy error : %v", err)
	}

	renderer.Present()

	return nil
}

func intro(renderer *sdl.Renderer) error {
	var err error

	err = drawTitle(renderer, "Speedy Car")

	if err != nil {
		return fmt.Errorf("draw title error :%v", err)
	}

	time.Sleep(time.Second * 3)

	return nil
}

func outro(renderer *sdl.Renderer) error {
	var err error

	err = drawTitle(renderer, "Game Over")
	if err != nil {
		return fmt.Errorf("draw title error :%v", err)
	}

	time.Sleep(time.Second * 3)

	return nil
}
