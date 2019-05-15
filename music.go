package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/mix"
)

type Musics struct {
	musics map[string]*mix.Music
}

func newMusics() *Musics {
	musicMap := make(map[string]*mix.Music)
	fmt.Println(len(musicMap))

	return &Musics{musics: musicMap}
}

func (musics *Musics) playIntro() error {
	var err error
	fmt.Println(musics)

	musics.musics["intro"], err = mix.LoadMUS("assets/intro.mp3")
	if err != nil {
		return fmt.Errorf("intro music load error :%v", err)
	}

	err = musics.musics["intro"].Play(1)
	if err != nil {
		return fmt.Errorf("intro music play error :%v", err)
	}

	return nil
}

func (musics *Musics) playOutro() error {
	musics.destroy()
	var err error

	musics.musics["outro"], err = mix.LoadMUS("assets/intro.mp3")
	if err != nil {
		return fmt.Errorf("intro music load error :%v", err)
	}

	err = musics.musics["outro"].Play(1)

	if err != nil {
		return fmt.Errorf("outro music play error :%v", err)
	}

	return nil
}

func (musics *Musics) playCar1() error {
	musics.destroy()
	var err error

	musics.musics["car1"], err = mix.LoadMUS("assets/loop_1.wav")
	if err != nil {
		return fmt.Errorf("car1 music load error :%v", err)
	}

	err = musics.musics["car1"].Play(1)

	if err != nil {
		return fmt.Errorf("car1 music play error :%v", err)
	}

	return nil
}

func (musics *Musics) playCar2() error {
	musics.destroy()
	var err error

	musics.musics["car2"], err = mix.LoadMUS("assets/loop_5.wav")
	if err != nil {
		return fmt.Errorf("car2 music load error :%v", err)
	}

	err = musics.musics["car2"].Play(1)

	if err != nil {
		return fmt.Errorf("car2 music play error :%v", err)
	}

	return nil
}

func (musics *Musics) destroy() {
	fmt.Println(len(musics.musics))
	if len(musics.musics) > 0 {
		for _, music := range musics.musics {
			music.Free()
		}
	}
}
