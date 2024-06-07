package audio

import (
	_ "embed"
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
