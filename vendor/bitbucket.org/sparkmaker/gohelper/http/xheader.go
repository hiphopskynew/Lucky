package http

import (
	"strings"
)

type XHeader map[string]string

var prefix string = "x-"

func filterXHeader(key string) string {
	newKey := strings.TrimSpace(key)
	if strings.HasPrefix(strings.ToLower(newKey), strings.ToLower(prefix)) {
		return newKey
	}

	return ""
}

func NewXHeader() XHeader {
	return make(map[string]string)
}

func (xh XHeader) SetPrefix(startWith string) XHeader {
	prefix = startWith
	for k, v := range xh {
		newKey := filterXHeader(k)
		if len(newKey) == 0 {
			delete(xh, k)
		} else {
			xh[newKey] = v
		}
	}
	return xh
}

func (xh XHeader) Append(headers map[string]string) XHeader {
	for k, v := range headers {
		xh[filterXHeader(k)] = v
	}
	delete(xh, "")
	return xh
}

func (xh XHeader) Add(key string, value string) XHeader {
	xh[filterXHeader(key)] = value
	delete(xh, "")
	return xh
}

func (xh XHeader) Del(key string) XHeader {
	delete(xh, key)
	return xh
}

func (xh XHeader) ToMap() map[string]string {
	return xh
}
