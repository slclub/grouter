package grouter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var path_decoder = NewPath()

func TestPathParse(t *testing.T) {
	var path_arr = [][]string{
		{"/abc", "/abc/"},
		{"/abc/:uid", "/abc/"},
		{"/abc?a=1&b=2", "/abc/"},
		{"//abc?a=1&b=2", "/abc/"},
		{"/abc?", "/abc/"},
	}
	for _, v := range path_arr {
		result, _ := path_decoder.Parse(v[0])
		assert.Equal(t, v[1], result)
	}

	// panics
	assert.Panics(t, func() { path_decoder.Parse("") })
}

// test path decode
func TestPathDecodeQuestion(t *testing.T) {
	path_type, result, param_str := path_decoder.Decode("/abc?uid=2&pid=1")
	assert.Equal(t, "/abc/", result)
	assert.Equal(t, "uid=2&pid=1", param_str)
	assert.Equal(t, PATH_T_QUESTION, path_type)
}

func TestPathDecodeCommon(t *testing.T) {

	path_type, result, param_str := path_decoder.Decode("/abc//uid/pid")
	assert.Equal(t, "/abc//uid/pid/", result)
	assert.Equal(t, "", param_str)
	assert.Equal(t, PATH_T_COMMON, path_type)

}

func TestAnGolangExaminationQuestions(t *testing.T) {
	num := []int{0, 1, 2, 3, 4}
	index := []int{0, 1, 2, 2, 1}
	//期待
	// [0]
	// [0,1]
	// [0,1, 2]
	// [0,1, 3, 2]
	// [0,4,1, 3, 2]

	var target = make([]int, 1)
	for i, v := range index {
		//	fmt.Println("TestAnGolangExaminationQuestions", v, target)
		if i == 0 {
			target[0] = num[0]
			continue
		}
		tmp_len := len(target)
		if v >= tmp_len {
			target = append(target, num[i])
			continue
		}
		if v < tmp_len {
			target = append(target, 0)
			var tmp int
			for j := v; j <= tmp_len; j++ {
				if v == j {
					tmp = target[j]
					target[j] = num[i]
					continue
				}
				target[j], tmp = tmp, target[j]
			}
			//// append 第一个参数不能 切片中部分， 否则会修改那个切片后面的元素值
			//var tmp = make([]int, 0)
			//tmp = append(tmp, target[:v]...)
			//tmp = append(tmp, num[i])

			////fmt.Println("TestAnGolangExaminationQuestions", "tmp", tmp, target[v:])
			//target = append(tmp, target[v:]...)
		}
	}

	fmt.Println("TestAnGolangExaminationQuestions", "last", target)
}
