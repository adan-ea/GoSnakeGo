package audio

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	//go:embed hit.ogg
	Hit_ogg []byte

	//go:embed eat.ogg
	Eat_ogg []byte

	//go:embed game_over.wav
	GameOver_wav []byte

	//go:embed tetris-theme.wav
	TetrisTheme_wav []byte
)

var (
	AudioContext   *audio.Context
	HitPlayer      *audio.Player
	EatPlayer      *audio.Player
	GameOverPlayer *audio.Player
	ThemePlayer    *audio.Player
)

func InitAudio() {
	if AudioContext == nil {
		AudioContext = audio.NewContext(48000)
	}

	hitD, err := vorbis.DecodeWithoutResampling(bytes.NewReader(Hit_ogg))
	if err != nil {
		log.Fatal(err)
	}

	HitPlayer, err = AudioContext.NewPlayer(hitD)
	if err != nil {
		log.Fatal(err)
	}

	eatD, err := vorbis.DecodeWithoutResampling(bytes.NewReader(Eat_ogg))
	if err != nil {
		log.Fatal(err)
	}

	EatPlayer, err = AudioContext.NewPlayer(eatD)
	if err != nil {
		log.Fatal(err)
	}

	gameOverD, err := wav.DecodeWithoutResampling(bytes.NewReader(GameOver_wav))
	if err != nil {
		log.Fatal(err)
	}

	themeD, err := wav.DecodeWithoutResampling(bytes.NewReader(TetrisTheme_wav))
	if err != nil {
		log.Fatal(err)
	}

	ThemePlayer, err = AudioContext.NewPlayer(themeD)
	if err != nil {
		log.Fatal(err)
	}

	GameOverPlayer, err = AudioContext.NewPlayer(gameOverD)
	if err != nil {
		log.Fatal(err)
	}
}

func PlayOnce(p *audio.Player) {
	p.Rewind()
	p.Play()
}

func PlayLoop(p *audio.Player) {
	if !p.IsPlaying() {
		p.Rewind()
		p.Play()
	}
}
