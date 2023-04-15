// Author: Maximilian Floto, Yannick Kirschen
package utils

import (
	"errors"
	"strconv"
	"strings"
)

func PathParameterFilter(path string, prefix string) (int64, error) {
	suffix := strings.TrimPrefix(path, prefix)
	var parameter string
	for _, char := range suffix {
		if char == '/' {
			return 0, errors.New("too many parameters (found something after /)")
		}

		parameter += string(char)
	}

	if len(parameter) > 0 {
		return strconv.ParseInt(parameter, 10, 64)
	}

	return 0, errors.New("no parameter found")
}
