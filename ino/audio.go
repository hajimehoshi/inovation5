package ino

import (
	"bytes"
	"strings"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/go-inovation/ino/internal/assets"
)

var (
	audioContext   *audio.Context
	soundFilenames = []string{
		"damage.wav",
		"heal.wav",
		"ino1.ogg",
		"ino2.ogg",
		"itemget.wav",
		"itemget2.wav",
		"jump.wav",
	}
	soundPlayers = map[string]*audio.Player{}
)

type bytesReadSeekCloser struct {
	r *bytes.Reader
}

func (b *bytesReadSeekCloser) Read(data []byte) (int, error) {
	return b.r.Read(data)
}

func (b *bytesReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return b.r.Seek(offset, whence)
}

func (b *bytesReadSeekCloser) Close() error {
	return nil
}

func (g *Game) loadAudio() {
	var err error
	defer func() {
		if err != nil {
			g.audioLoadedCh <- err
		}
		close(g.audioLoadedCh)
	}()

	const sampleRate = 44100
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		return
	}
	for _, n := range soundFilenames {
		var b []byte
		b, err = assets.Asset("resources/sound/" + n)
		if err != nil {
			return
		}
		f := &bytesReadSeekCloser{bytes.NewReader(b)}
		var s audio.ReadSeekCloser
		switch {
		case strings.HasSuffix(n, ".ogg"):
			var stream *vorbis.Stream
			stream, err = vorbis.Decode(audioContext, f)
			if err != nil {
				return
			}
			s = NewLoop(stream, stream.Size())
		case strings.HasSuffix(n, ".wav"):
			s, err = wav.Decode(audioContext, f)
			if err != nil {
				return
			}
		default:
			panic("invalid file name")
		}
		var p *audio.Player
		p, err = audio.NewPlayer(audioContext, s)
		if err != nil {
			return
		}
		soundPlayers[n] = p
	}
}

func finalizeAudio() error {
	for _, p := range soundPlayers {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

type BGM string

const (
	BGM0 BGM = "ino1.ogg"
	BGM1 BGM = "ino2.ogg"
)

func SetBGMVolume(volume float64) {
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		if !p.IsPlaying() {
			continue
		}
		p.SetVolume(volume)
		return
	}
}

func PauseBGM() error {
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		if err := p.Pause(); err != nil {
			return err
		}
	}
	return nil
}

func ResumeBGM(bgm BGM) error {
	if err := PauseBGM(); err != nil {
		return err
	}
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	return p.Play()
}

func PlayBGM(bgm BGM) error {
	if err := PauseBGM(); err != nil {
		return err
	}
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	if err := p.Rewind(); err != nil {
		return err
	}
	return p.Play()
}

type SE string

const (
	SE_DAMAGE   SE = "damage.wav"
	SE_HEAL     SE = "heal.wav"
	SE_ITEMGET  SE = "itemget.wav"
	SE_ITEMGET2 SE = "itemget2.wav"
	SE_JUMP     SE = "jump.wav"
)

func PlaySE(se SE) error {
	p := soundPlayers[string(se)]
	if err := p.Rewind(); err != nil {
		return err
	}
	if err := p.Play(); err != nil {
		return err
	}
	return nil
}
