package field

import (
	"strings"

	"github.com/hajimehoshi/ebiten"

	"github.com/hajimehoshi/go-inovation/ino/internal/draw"
)

type FieldType int

const (
	FIELD_NONE         FieldType = iota // なし
	FIELD_HIDEPATH                      // 隠しルート(見えるけど判定のないブロック)
	FIELD_UNVISIBLE                     // 不可視ブロック(見えないけど判定があるブロック)
	FIELD_BLOCK                         // 通常ブロック
	FIELD_BAR                           // 床。降りたり上ったりできる
	FIELD_SCROLL_L                      // ベルト床左
	FIELD_SCROLL_R                      // ベルト床右
	FIELD_SPIKE                         // トゲ
	FIELD_SLIP                          // すべる
	FIELD_ITEM_BORDER                   // アイテムチェック用
	FIELD_ITEM_POWERUP                  // パワーアップ
	// ふじ系
	FIELD_ITEM_FUJI
	FIELD_ITEM_BUSHI
	FIELD_ITEM_APPLE
	FIELD_ITEM_V
	// たか系
	FIELD_ITEM_TAKA
	FIELD_ITEM_SHUOLDER
	FIELD_ITEM_DAGGER
	FIELD_ITEM_KATAKATA
	// なす系
	FIELD_ITEM_NASU
	FIELD_ITEM_BONUS
	FIELD_ITEM_NURSE
	FIELD_ITEM_NAZUNA
	// くそげー系
	FIELD_ITEM_GAMEHELL
	FIELD_ITEM_GUNDAM
	FIELD_ITEM_POED
	FIELD_ITEM_MILESTONE
	FIELD_ITEM_1YEN
	FIELD_ITEM_TRIANGLE
	FIELD_ITEM_OMEGA      // 隠し
	FIELD_ITEM_LIFE       // ハート
	FIELD_ITEM_STARTPOINT // 開始地点
	FIELD_ITEM_MAX
)

const (
	CHAR_SIZE = 16
	maxFieldX = 128
	maxFieldY = 128
)

type Field struct {
	field [maxFieldX * maxFieldY]FieldType
	timer int
}

func New(data string) *Field {
	f := &Field{}
	xm := strings.Split(data, "\n")
	const decoder = " HUB~<>*I PabcdefghijklmnopqrzL@"

	for yy, line := range xm {
		for xx, c := range line {
			n := strings.IndexByte(decoder, byte(c))
			f.field[yy*maxFieldX+xx] = FieldType(n)
		}
	}
	return f
}

func (f *Field) Update() {
	f.timer++
}

func (f *Field) GetStartPoint() (int, int) {
	for yy := 0; yy < maxFieldY; yy++ {
		for xx := 0; xx < maxFieldX; xx++ {
			if f.GetField(xx, yy) == FIELD_ITEM_STARTPOINT {
				x := xx * CHAR_SIZE
				y := yy * CHAR_SIZE
				f.EraseField(xx, yy)
				return x, y
			}
		}
	}
	panic("no start point")
}

func (f *Field) IsWall(x, y int) bool {
	return f.field[y*maxFieldX+x] != FIELD_NONE &&
		f.field[y*maxFieldX+x] != FIELD_HIDEPATH &&
		f.field[y*maxFieldX+x] != FIELD_BAR &&
		!f.IsItem(x, y)
}
func (f *Field) IsRidable(x, y int) bool {
	return f.field[y*maxFieldX+x] != FIELD_NONE &&
		f.field[y*maxFieldX+x] != FIELD_HIDEPATH &&
		!f.IsItem(x, y)
}

func (f *Field) IsSpike(x, y int) bool {
	return f.field[y*maxFieldX+x] == FIELD_SPIKE
}

func (f *Field) GetField(x, y int) FieldType {
	return f.field[y*maxFieldX+x]
}

func (f *Field) IsItem(x, y int) bool {
	return f.field[y*maxFieldX+x] >= FIELD_ITEM_BORDER &&
		f.field[y*maxFieldX+x] != FIELD_ITEM_STARTPOINT
}

func (f *Field) IsItemGettable(x, y int, gameData GameData) bool {
	if !f.IsItem(x, y) {
		return false
	}
	if f.field[y*maxFieldX+x] == FIELD_ITEM_OMEGA && gameData.IsHiddenSecret() {
		return false
	}
	return true
}

func (f *Field) EraseField(x, y int) {
	f.field[y*maxFieldX+x] = FIELD_NONE
}

type GameData interface {
	IsHiddenSecret() bool
}

func (f *Field) Draw(screen *ebiten.Image, gameData GameData, viewPositionX, viewPositionY int) {
	const (
		graphicOffsetX = -16 - 16*2
		graphicOffsetY = 8 - 16*2
	)
	vx, vy := viewPositionX, viewPositionY
	ofs_x := CHAR_SIZE - vx%CHAR_SIZE
	ofs_y := CHAR_SIZE - vy%CHAR_SIZE
	for xx := -(draw.ScreenWidth/CHAR_SIZE/2 + 2); xx < (draw.ScreenWidth/CHAR_SIZE/2 + 2); xx++ {
		fx := xx + vx/CHAR_SIZE
		if fx < 0 || fx >= maxFieldX {
			continue
		}
		for yy := -(draw.ScreenHeight/CHAR_SIZE/2 + 2); yy < (draw.ScreenHeight/CHAR_SIZE/2 + 2); yy++ {
			fy := yy + vy/CHAR_SIZE
			if fy < 0 || fy >= maxFieldY {
				continue
			}

			gy := (f.timer / 10) % 4
			gx := int(f.field[fy*maxFieldX+fx])

			if f.IsItem(fx, fy) {
				gx = gx - (int(FIELD_ITEM_BORDER) + 1)
				gy = 4 + gx/16
				gx = gx % 16
			}

			if gameData.IsHiddenSecret() && f.field[fy*maxFieldX+fx] == FIELD_ITEM_OMEGA {
				continue
			}

			draw.Draw(screen, "ino",
				(xx+12)*CHAR_SIZE+ofs_x+graphicOffsetX+(draw.ScreenWidth-320)/2,
				(yy+8)*CHAR_SIZE+ofs_y+graphicOffsetY+(draw.ScreenHeight-240)/2,
				gx*16, gy*16, 16, 16)
		}
	}
}
