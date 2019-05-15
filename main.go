package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

func main() {
	sdl.Main(func() {
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(2)
		}
	})
}

func run() error {
	var err error

	sdl.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			sdl.Quit()
		})
	}()

	sdl.Do(func() {
		err = ttf.Init()
	})
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			ttf.Quit()
		})
	}()

	sdl.Do(func() {
		err = mix.Init(mix.INIT_MP3)
	})
	if err != nil {
		return fmt.Errorf("could not initialize MIX: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			mix.Quit()
		})
	}()

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		return fmt.Errorf("mix audio error :%v", err)
	}
	defer func() {
		sdl.Do(func() {
			mix.CloseAudio()
		})
	}()

	var window *sdl.Window
	var renderer *sdl.Renderer

	sdl.Do(func() {
		window, renderer, err = sdl.CreateWindowAndRenderer(1415, 600, sdl.WINDOW_SHOWN)
	})
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			window.Destroy()
		})
	}()

	err = intro(renderer)

	if err != nil {
		return fmt.Errorf("intro error :%v", err)
	}

	var scene *scene

	sdl.Do(func() {
		scene, err = newScene(renderer)
	})
	if err != nil {
		return fmt.Errorf("scene create error :%v", err)
	}

	defer sdl.Do(func() {
		scene.destroy()
	})

	err = scene.run(renderer)

	if err != nil {
		return fmt.Errorf("scene run error :%v", err)
	}
	return err
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

	sdl.Do(func() {
		err = drawTitle(renderer, "Speedy Car")
	})

	if err != nil {
		return fmt.Errorf("draw title error :%v", err)
	}

	time.Sleep(time.Second * 3)

	return nil
}

func outro(renderer *sdl.Renderer) error {
	var err error

	sdl.Do(func() {
		err = drawTitle(renderer, "Game Over")
	})

	if err != nil {
		return fmt.Errorf("draw title error :%v", err)
	}

	time.Sleep(time.Second * 3)

	return nil
}
