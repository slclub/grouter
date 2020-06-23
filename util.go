package grouter

import (
	"bytes"
	//"sync"
)

const (
	LOWER_MIN_LETER_UINT8   = uint8(97) // lowercase start leter ascll.
	CAPTIAL_MIN_LETER_UINT8 = uint8(65) // captial start leter ascll.
	DIFF_LETER_UINT8        = uint8(32) // diff num = 32

)

// Handle url leter captail to lowercase for lookup.
func lowercase(indices string) string {
	//return indices
	lenp := len(indices)
	buf := bytes.NewBufferString(indices)
	buf.Reset()
	var leter_num uint8 = 26
	for i := 0; i < lenp; i++ {
		tmp_l := byte(indices[i])
		// captail leter
		// 26 个大写英文字母 转小写
		if tmp_l >= CAPTIAL_MIN_LETER_UINT8 && tmp_l < CAPTIAL_MIN_LETER_UINT8+leter_num {
			tmp_l += DIFF_LETER_UINT8
		}
		buf.WriteByte(tmp_l)
	}
	// for performance testing code .
	// return indices

	// 3 allocs/op
	// return string(buf.Bytes())

	return buf.String()
}

func lowercaseWithBuffer(buf []byte, b byte, w int) {
	var leter_num uint8 = 26
	if b >= CAPTIAL_MIN_LETER_UINT8 && b < CAPTIAL_MIN_LETER_UINT8+leter_num {
		b += DIFF_LETER_UINT8
	}

	buf[w] = b
}
