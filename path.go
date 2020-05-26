package grouter

import (
	"strings"
)

const (
	PATH_T_COMMON   = 1
	PATH_T_QUESTION = 2
	PATH_T_FILE     = 3
)

// url decode.
// implement path.
type urldecoder struct {
	// support router case sensitive.
	case_sensitive bool
}

func NewPath() Path {
	return &urldecoder{
		case_sensitive: false,
	}
}

// invoked by router add route.
// here the path first character must be /. we judgy it by Router.handle.
// TODO:file upload and donwload
func (u *urldecoder) Parse(path string) (string, []string) {
	if u.case_sensitive {
		path = strings.ToLower(path)
	}
	path = strings.Replace(path, "//", "/", -1)
	lenp := len(path)
	path_rtn, path_left := "", ""
	for i, v := range path {
		if v == '?' {
			return u.ParseQuestion(path, i)
		}

		if v == ':' && lenp > 2 && (i+2) < lenp {
			path_rtn = path[:i]
			path_left = path[i-1:]
			break
		}
	}

	// no param keys.
	if path_rtn == "" {
		if path[lenp-1] != '/' {
			return path + "/", nil
		}
		if path_left != "" {
			panic("[ERROR][GROUTER][PATH][NOT_VALID]")
		}
		return path, nil
	}

	if len(path_left) > 0 && path_left[len(path_left)-1] == '/' {
		path_left = path_left[:len(path_left)-2]
	}
	param_keys := strings.Split(path_left, "/:")
	return path_rtn, param_keys
}

// we can use ? define the router.
func (u *urldecoder) ParseQuestion(path string, position int) (string, []string) {
	var path_rtn string
	// TODO: path valid check
	if path[len(path)-1] != '/' {
		path_rtn = path[:position] + "/"
	}

	param_keys := strings.Split(path[position+1:], "&")

	// get key string before equal sign.
	// if force user observe the rule "?first&seconde&otherRouter", no equal sign. the following code is unnecessary.
	for i, v := range param_keys {
		for k, vv := range v {
			if vv == '=' {
				param_keys[i] = v[:k]
			}
		}
	}
	return path_rtn, param_keys
}

// whether url support case sensitive.
func (u *urldecoder) CaseSensitive(status bool) {
	u.case_sensitive = status
}

// parse client request url. example http request.
func (u *urldecoder) Decode(path string) (int, string, interface{}) {
	// TODO: encoding url  %f etc convert to string.
	//
	// 加号处理可以提供方法尽量不在这里处理，会影响性能
	// Processing plus signs in other places can improve some performance

	// deal split path.
	// first section is url path. second section is params.
	path_buf, param_str, path_type := u.convPath(path)
	return path_type, string(path_buf), param_str
}

func (u *urldecoder) convPath(path string) ([]byte, string, int) {
	lenp := len(path)
	if lenp > 0 && path[0] == '*' {
		return nil, "", 0
	}

	buf := make([]byte, lenp+1)
	// write buf size.
	w := 0
	if path[0] != '/' {
		buf[w] = '/'
		w++
	}

	// [65, 97) if byte uin8 in this section, that is capital letter.
	// convert to lower letter need -32. cant handle it here.
	// param value should not to be converted.
	for i := 0; i < lenp; i++ {
		// replace // or /// to /
		//if path[i] == '/' && w > 0 && buf[w-1] == '/' {
		//	continue
		//}

		if path[i] == '?' {
			buf[w] = '/'
			w++
			if i == (lenp - 1) {
				return buf[:w], "", PATH_T_QUESTION
			} else {
				return buf[:w], path[i+1:], PATH_T_QUESTION
			}
		}
		buf[w] = path[i]
		w++
		if i == (lenp-1) && path[i] != '/' {
			buf[w] = '/'
			w++
		}
	}
	return buf[:w], "", PATH_T_COMMON
}
