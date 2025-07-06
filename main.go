package main

/*
tinygo flash --target waveshare-rp2040-zero --size short -monitor .
mkdir uf2
tinygo build -o uf2/7Keyx3Oct.uf2 --target waveshare-rp2040-zero --size short .
*/

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/tone"
)

/*
|      | col0  | col1  | col2  | col3  |
| ---- | ----- | ----  | ----  | ----- |
| row0 | sw1   | sw2   | sw3   | sw4   |
| row1 | sw5   | sw6   | sw7   | sw8   |
| row2 | sw9   | sw10  | sw11  | sw12  |

|      | col0  | col1  | col2  | col3  |
| ---- | ----- | ----- | ----- | ----- |
| row0 | Oct3  | Key7  | Flat  | Sharp |
| row1 | Oct4  | Key4  | Key5  | Key6  |
| row2 | Oct5  | Key1  | Key2  | Key3  |

|      |col0 |col1 |col2 |col3 |
| ---- | --- | --- | --- | --- |
| row0 | O3  |  7  | ♭  |  ＃  |
| row1 | O4  |  4  |  5  |  6  |
| row2 | O5  |  1  |  2  |  3  |

＃	Sharp	半音上げる
♭	Flat	半音下げる

ノートナンバー 60：88鍵ピアノにおける中央のド
ノートナンバー 69：チューニングの標準音（440Hz）
中央のド：C4（88鍵ピアノにおける左から4番目のC）
*/

// zero-kb02用出力ポート等の定義
var pinToPWM = map[machine.Pin]tone.PWM{
	machine.GPIO14: machine.PWM7, // for EX01
	machine.GPIO15: machine.PWM7, // for EX02
	machine.GPIO26: machine.PWM5, // for EX01
	machine.GPIO27: machine.PWM5, // for EX01
}

// 入力のあった音階キーと拡張キーの状態を保持する変数群
var (
	note         int    = 0  // 音階キー
	number_score int    = 0  // 数字譜
	note_name    string = "" // 音名
	flag_sharp   int    = 0  // ＃キー
	flag_flat    int    = 0  // ♭キー
	octabe       int    = 0  // オクターブキー
)

// 列(縦方向)のPin定義
var colPins = []machine.Pin{
	machine.GPIO5,
	machine.GPIO6,
	machine.GPIO7,
	machine.GPIO8,
}

// 行(横方向)のPin定義
var rowPins = []machine.Pin{
	machine.GPIO9,
	machine.GPIO10,
	machine.GPIO11,
}

func NoteNoToPeriod(NoteNo int) uint64 {
	freq := freq_table[NoteNo]
	val := 1e9 / freq
	//	fmt.Printf("%v,%v,%v,%v\n", NoteNo, freq, val, uint64(val+0.5))
	return uint64(val)
}

func main() {
	// 列(縦方向)のPinの初期化
	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	// 行(横方向)のPinの初期化
	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	// ブザーが接続されたPinの設定と初期化
	bzrPin := machine.GPIO14
	pwm := pinToPWM[bzrPin]
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	for {
		note = 0
		number_score = 0
		flag_sharp = 0
		flag_flat = 0
		octabe = 0
		// COL1
		colPins[0].High()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			//	fmt.Printf("Oct5:sw1 pressed\n")
			octabe = 12
		}
		if rowPins[1].Get() {
			//	fmt.Printf("Oct4:sw5 pressed\n")
			octabe = 0
		}
		if rowPins[2].Get() {
			//	fmt.Printf("Oct3:sw9 pressed\n")
			octabe = -12
		}

		// COL2
		colPins[0].Low()
		colPins[1].High()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			//	fmt.Printf("Key7:sw2 pressed\n")
			note = 71 // Oct4,B  シ
		}
		if rowPins[1].Get() {
			//	fmt.Printf("Key4:sw6 pressed\n")
			note = 65 // Oct4,F  ファ
		}
		if rowPins[2].Get() {
			//	fmt.Printf("Key1:sw10 pressed\n")
			note = 60 // Oct4,C  ド
		}

		// COL3
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].High()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			//	fmt.Printf("Flat :sw3 pressed\n")
			flag_flat = -1 // 半音下げる
		}
		if rowPins[1].Get() {
			//	fmt.Printf("Key5:sw7 pressed\n")
			note = 67 // Oct4,G  ソ
		}
		if rowPins[2].Get() {
			//	fmt.Printf("Key2:sw11 pressed\n")
			note = 62 // Oct4,D  レ
		}

		// COL4
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].High()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			//	fmt.Printf("Sharp:sw4 pressed\n")
			flag_sharp = 1 // 半音上げる
		}
		if rowPins[1].Get() {
			//	fmt.Printf("Key6:sw8 pressed\n")
			note = 69 // Oct4,A  ラ
		}
		if rowPins[2].Get() {
			//	fmt.Printf("Key3:sw12 pressed\n")
			note = 64 // Oct4,E  ミ
		}
		// 入力されたキーを元に音階を決定する。
		// オクターブと半音の情報を反映させる。

		if 0 != note {
			switch note {
			case 60: // C,ド
				note = octabe + note + flag_sharp
				number_score = 1
				note_name = "Do "
				break
			case 62: // D,レ
				note = octabe + note + flag_sharp + flag_flat
				number_score = 2
				note_name = "Re "
				break
			case 64: // E,ミ
				note = octabe + note + flag_flat
				number_score = 3
				note_name = "Mi "
				break
			case 65: // F,ファ
				note = octabe + note + flag_sharp
				number_score = 4
				note_name = "Fa "
				break
			case 67: // G,ソ
				note = octabe + note + flag_sharp + flag_flat
				number_score = 5
				note_name = "Sol"
				break
			case 69: // A,ラ
				note = octabe + note + flag_sharp + flag_flat
				number_score = 6
				note_name = "La "
				break
			case 71: // B,シ
				note = octabe + note + flag_flat
				number_score = 7
				note_name = "Si "
				break
			}
			semitone := " "
			if 0 != flag_flat {
				semitone = "♭"
			}
			if 0 != flag_sharp {
				semitone = "#"
			}

			var octabe_table = map[int]string{
				-12: "Oct3", // 3オクターブ
				0:   "Oct4", // 4オクターブ
				12:  "Oct5", // 5オクターブ
			}
			fmt.Printf("| %s | %d | %s | %3s |\n", octabe_table[octabe], number_score, semitone, note_name)
			speaker.SetPeriod(NoteNoToPeriod(note))
			time.Sleep(100 * time.Millisecond)
		} else {
			speaker.Stop()
		}
	}
}
