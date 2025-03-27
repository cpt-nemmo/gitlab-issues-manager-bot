package utils

import (
	"errors"
	"strings"
)

const (
	START = "title:"
	END   = "desc:"
)

type Issue struct {
	Title       string
	Description string
}

func ParseIssue(s string) (*Issue, error) {
	s = strings.ReplaceAll(s, "#issue", "")
	startIdx := strings.Index(s, START) + len(START) // это первый пробел
	slice := s[startIdx:]
	endIdx := strings.Index(slice, END)
	if endIdx == -1 {
		return nil, errors.New("can't parse issue")
	}
	result := strings.TrimSpace(slice[:endIdx])
	//fmt.Printf("res: [%v]\n", result)

	newStartIndex := strings.Index(slice, END) + len(END)
	if newStartIndex == -1 {
		return nil, errors.New("can't parse issue")
	}

	result2 := strings.TrimSpace(slice[newStartIndex:])
	//fmt.Printf("res2: [%v]\n", result2)

	return &Issue{
		Title:       result,
		Description: result2,
	}, nil
}
