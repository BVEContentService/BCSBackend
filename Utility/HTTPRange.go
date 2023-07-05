package Utility

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ParseRange(c *gin.Context, count int) (int, int, error) {
	header := strings.TrimSpace(c.GetHeader("Range"))
	if header == "" {
		return 0, count, nil
	}
	parted1 := strings.Split(header, "=")
	if len(parted1) != 2 {
		return 0, 0, ERR_BAD_RANGE
	}
	parted2 := strings.Split(strings.TrimSpace(parted1[1]), "-")
	if len(parted2) != 2 {
		return 0, 0, ERR_BAD_RANGE
	}
	begin, err1 := strconv.Atoi(strings.TrimSpace(parted2[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parted2[1]))
	if err1 != nil || err2 != nil || end < begin {
		return 0, 0, ERR_BAD_RANGE
	}
	return min(max(begin, 0), count), min(max(end+1, 0), count), nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
