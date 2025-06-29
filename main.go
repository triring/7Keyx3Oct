package main

/*
tinygo flash --target waveshare-rp2040-zero --size short -monitor .
mkdir uf2
tinygo build -o uf2/7Keyx3Oct.uf2 --target waveshare-rp2040-zero --size short ./main.go
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
| row0 | O3  |  7  | ♭  |  #  |
| row1 | O4  |  4  |  5  |  6  |
| row2 | O5  |  1  |  2  |  3  |


|      |col0 |col1 |col2 |col3 |
| ---- | --- | --- | --- | --- |
| row0 | O3  |  71 | ♭  |  #  |
| row1 | O4  |  65 |  67 |  69 |
| row2 | O5  |  60 |  62 |  64 |


＃	Sharp	半音上げる
♭	Flat	半音下げる

ノートナンバー 60：88鍵ピアノにおける中央のド
ノートナンバー 69：チューニングの標準音（440Hz）
中央のド：C4（88鍵ピアノにおける左から4番目のC）
*/

// zero-kb02用出力ポート等の定期
var pinToPWM = map[machine.Pin]tone.PWM{
	machine.GPIO14: machine.PWM7, // for EX01
	machine.GPIO15: machine.PWM7, // for EX02
	machine.GPIO26: machine.PWM5, // for EX01
	machine.GPIO27: machine.PWM5, // for EX01
}

var (
	note int = 0
	//	key        int = 0
	flag_sharp int = 0
	flag_flat  int = 0
	octabe     int = 0
)

var colPins = []machine.Pin{
	machine.GPIO5,
	machine.GPIO6,
	machine.GPIO7,
	machine.GPIO8,
}

var rowPins = []machine.Pin{
	machine.GPIO9,
	machine.GPIO10,
	machine.GPIO11,
}

var freq_table = map[int]float64{
	0:   8.25000000,    // Oct-1,C  ド
	1:   8.80000000,    // Oct-1,C#
	2:   9.28125000,    // Oct-1,D  レ
	3:   9.75000000,    // Oct-1,D#
	4:   10.31250000,   // Oct-1,E  ミ
	5:   11.00000000,   // Oct-1,F  ファ
	6:   11.60156250,   // Oct-1,F#
	7:   12.37500000,   // Oct-1,G  ソ
	8:   13.20000000,   // Oct-1,G#
	9:   13.75000000,   // Oct-1,A  ラ
	10:  14.85000000,   // Oct-1,A#
	11:  15.46875000,   // Oct-1,B  シ
	12:  16.50000000,   // Oct0,C  ド
	13:  17.60000000,   // Oct0,C#
	14:  18.56250000,   // Oct0,D  レ
	15:  19.50000000,   // Oct0,D#
	16:  20.62500000,   // Oct0,E  ミ
	17:  22.00000000,   // Oct0,F  ファ
	18:  23.20312500,   // Oct0,F#
	19:  24.75000000,   // Oct0,G  ソ
	20:  26.40000000,   // Oct0,G#
	21:  27.50000000,   // Oct0,A  ラ
	22:  29.70000000,   // Oct0,A#
	23:  30.93750000,   // Oct0,B  シ
	24:  33.00000000,   // Oct1,C  ド
	25:  35.20000000,   // Oct1,C#
	26:  37.12500000,   // Oct1,D  レ
	27:  39.00000000,   // Oct1,D#
	28:  41.25000000,   // Oct1,E  ミ
	29:  44.00000000,   // Oct1,F  ファ
	30:  46.40625000,   // Oct1,F#
	31:  49.50000000,   // Oct1,G  ソ
	32:  52.80000000,   // Oct1,G#
	33:  55.00000000,   // Oct1,A  ラ
	34:  59.40000000,   // Oct1,A#
	35:  61.87500000,   // Oct1,B  シ
	36:  66.00000000,   // Oct2,C  ド
	37:  70.40000000,   // Oct2,C#
	38:  74.25000000,   // Oct2,D  レ
	39:  78.00000000,   // Oct2,D#
	40:  82.50000000,   // Oct2,E  ミ
	41:  88.00000000,   // Oct2,F  ファ
	42:  92.81250000,   // Oct2,F#
	43:  99.00000000,   // Oct2,G  ソ
	44:  105.60000000,  // Oct2,G#
	45:  110.00000000,  // Oct2,A  ラ
	46:  118.80000000,  // Oct2,A#
	47:  123.75000000,  // Oct2,B  シ
	48:  132.00000000,  // Oct3,C  ド
	49:  140.80000000,  // Oct3,C#
	50:  148.50000000,  // Oct3,D  レ
	51:  156.00000000,  // Oct3,D#
	52:  165.00000000,  // Oct3,E  ミ
	53:  176.00000000,  // Oct3,F  ファ
	54:  185.62500000,  // Oct3,F#
	55:  198.00000000,  // Oct3,G  ソ
	56:  211.20000000,  // Oct3,G#
	57:  220.00000000,  // Oct3,A  ラ
	58:  237.60000000,  // Oct3,A#
	59:  247.50000000,  // Oct3,B  シ
	60:  264.00000000,  // Oct4,C  ド
	61:  281.60000000,  // Oct4,C#
	62:  297.00000000,  // Oct4,D  レ
	63:  312.00000000,  // Oct4,D#
	64:  330.00000000,  // Oct4,E  ミ
	65:  352.00000000,  // Oct4,F  ファ
	66:  371.25000000,  // Oct4,F#
	67:  396.00000000,  // Oct4,G  ソ
	68:  422.40000000,  // Oct4,G#
	69:  440.00000000,  // Oct4,A  ラ
	70:  475.20000000,  // Oct4,A#
	71:  495.00000000,  // Oct4,B  シ
	72:  528.00000000,  // Oct5,C  ド
	73:  563.20000000,  // Oct5,C#
	74:  594.00000000,  // Oct5,D  レ
	75:  624.00000000,  // Oct5,D#
	76:  660.00000000,  // Oct5,E  ミ
	77:  704.00000000,  // Oct5,F  ファ
	78:  742.50000000,  // Oct5,F#
	79:  792.00000000,  // Oct5,G  ソ
	80:  844.80000000,  // Oct5,G#
	81:  880.00000000,  // Oct5,A  ラ
	82:  950.40000000,  // Oct5,A#
	83:  990.00000000,  // Oct5,B  シ
	84:  1056.00000000, // Oct6,C  ド
	85:  1126.40000000, // Oct6,C#
	86:  1188.00000000, // Oct6,D  レ
	87:  1248.00000000, // Oct6,D#
	88:  1320.00000000, // Oct6,E  ミ
	89:  1408.00000000, // Oct6,F  ファ
	90:  1485.00000000, // Oct6,F#
	91:  1584.00000000, // Oct6,G  ソ
	92:  1689.60000000, // Oct6,G#
	93:  1760.00000000, // Oct6,A  ラ
	94:  1900.80000000, // Oct6,A#
	95:  1980.00000000, // Oct6,B  シ
	96:  2112.00000000, // Oct7,C  ド
	97:  2252.80000000, // Oct7,C#
	98:  2376.00000000, // Oct7,D  レ
	99:  2496.00000000, // Oct7,D#
	100: 2640.00000000, // Oct7,E  ミ
	101: 2816.00000000, // Oct7,F  ファ
	102: 2970.00000000, // Oct7,F#
	103: 3168.00000000, // Oct7,G  ソ
	104: 3379.20000000, // Oct7,G#
	105: 3520.00000000, // Oct7,A  ラ
	106: 3801.60000000, // Oct7,A#
	107: 3960.00000000, // Oct7,B  シ
	108: 4224.00000000, // Oct8,C  ド
	109: 4505.60000000, // Oct8,C#
	110: 4752.00000000, // Oct8,D  レ
	111: 4992.00000000, // Oct8,D#
	112: 5280.00000000, // Oct8,E  ミ
	113: 5632.00000000, // Oct8,F  ファ
	114: 5940.00000000, // Oct8,F#
	115: 6336.00000000, // Oct8,G  ソ
	116: 6758.40000000, // Oct8,G#
	117: 7040.00000000, // Oct8,A  ラ
	118: 7603.20000000, // Oct8,A#
	119: 7920.00000000, // Oct8,B  シ
	120: 4224.00000000, // Oct9,C  ド
	121: 4505.60000000, // Oct9,C#
	122: 4752.00000000, // Oct9,D  レ
	123: 4992.00000000, // Oct9,D#
	124: 5280.00000000, // Oct9,E  ミ
	125: 5632.00000000, // Oct9,F  ファ
	126: 5940.00000000, // Oct9,F#
	127: 6336.00000000, // Oct9,G  ソ
}

func NoteNoToPeriod(NoteNo int) uint64 {
	freq := freq_table[NoteNo]
	val := 1e9 / freq
	fmt.Printf("%v,%v,%v,%v\n", NoteNo, freq, val, uint64(val+0.5))
	return uint64(val)
}

func main() {
	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	bzrPin := machine.GPIO14
	pwm := pinToPWM[bzrPin]
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}
	octabe = 0
	for {
		note = 0
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
				break
			case 62: // D,レ
				note = octabe + note + flag_sharp + flag_flat
				break
			case 64: // E,ミ
				note = octabe + note + flag_flat
				break
			case 65: // F,ファ
				note = octabe + note + flag_sharp
				break
			case 67: // G,ソ
				note = octabe + note + flag_sharp + flag_flat
				break
			case 69: // A,ラ
				note = octabe + note + flag_sharp + flag_flat
				break
			case 71: // B,シ
				note = octabe + note + flag_flat
				break
			}
			// speaker.SetNote(scale_table[note])
			speaker.SetPeriod(NoteNoToPeriod(note))
			fmt.Printf("scale:%v,note:%d,octabe:%d,flag_sharp:%d,flag_flat:%d,\n", freq_table[note], note, octabe, flag_sharp, flag_flat)
			time.Sleep(100 * time.Millisecond)
		} else {
			speaker.Stop()
		}
	}
}
