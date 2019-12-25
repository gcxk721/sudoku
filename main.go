package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

type Cell struct {
	Row            int   // 行
	Line           int   // 列
	Piece          int   // 块
	Value          int   // 确定的值
	PossibleValues []int // 有可能出现的值
}

func main() {
	c, err := ini.Load("input.ini")
	if err != nil {
		log.Panic(err.Error())
	}
	lst := make([]Cell, 0)
	r := c.Section("topic")
	for i := 0; i < 9; i++ {
		line := r.Key(fmt.Sprintf("line%d", i+1)).String()
		for j := 0; j < 9; j++ {
			lst = append(lst, Cell{
				Row:   i,
				Line:  j,
				Piece: StringToInt(fmt.Sprintf("%d%d", i/3, j/3)),
				Value: StringToInt(line[j : j+1]),
			})
		}
	}
	lst = Puzzle(lst)
	if !Judge(lst) {
		// 不能确定 开始猜想解谜
		lst = Guess(lst)
	}
	// 打印结果
	PrintResult(lst)
}

func Puzzle(lst []Cell) []Cell {
	lst = PossibleValues(lst)
	lst, flag := Exclude(Sure(lst))
	if flag {
		Puzzle(lst)
	}
	return lst
}

func Guess(lst []Cell) []Cell {
	index, last := FindGuessCell(lst, 2)
	for _, may := range last.PossibleValues {
		last.Value = may
		lst[index] = last
		nl := Puzzle(lst)
		if !Judge(nl) {
			// 无解还是继续猜
			if !NoResult(nl) { // 接着猜
				return Guess(nl)
			}
		} else {
			return nl
		}
	}
	return lst
}

// 找到可能值最少的块
func FindGuessCell(lst []Cell, count int) (int, Cell) {
	for i, cell := range lst {
		if len(cell.PossibleValues) == count {
			return i, cell
		}
	}
	return FindGuessCell(lst, count+1)
}

// Sure 根据可能性列表，找到只有一个可能的块，成为确定值
func Sure(lst []Cell) []Cell {
	for i, cell := range lst {
		if cell.Value == 0 && len(cell.PossibleValues) == 1 {
			cell.Value = cell.PossibleValues[0]
			cell.PossibleValues = nil
			lst[i] = cell
			// 有确定值，重新计算可能性列表
			return Sure(PossibleValues(lst))
		}
	}
	return lst
}

// PossibleValues 计算每个块可能出现的值
func PossibleValues(lst []Cell) []Cell {
	for i, cell := range lst {
		if cell.Value == 0 {
			cell.PossibleValues = make([]int, 0)
			have := make([]int, 0)
			for _, other := range lst {
				if !Same(cell, other) && other.Value != 0 {
					// 找到线，块，行已经出现的值
					if other.Row == cell.Row || other.Line == cell.Line || other.Piece == cell.Piece {
						have = append(have, other.Value)
					}
				}
			}
			for v := 1; v < 10; v++ {
				if !Contain(v, have) {
					cell.PossibleValues = append(cell.PossibleValues, v)
				}
			}
		} else {
			cell.PossibleValues = nil
		}
		lst[i] = cell
	}
	return lst
}

// Exclude 通过排除算法，确定块的值
func Exclude(lst []Cell) ([]Cell, bool) {
	excludePiece := func(val int, cell Cell) bool {
		for _, other := range lst {
			if !Same(other, cell) && other.Piece == cell.Piece && other.Value == 0 && Contain(val, other.PossibleValues) {
				return false
			}
		}
		return true
	}
	excludeLine := func(val int, cell Cell) bool {
		for _, other := range lst {
			if !Same(other, cell) && other.Line == cell.Line && other.Value == 0 && Contain(val, other.PossibleValues) {
				return false
			}
		}
		return true
	}
	excludeRow := func(val int, cell Cell) bool {
		for _, other := range lst {
			if !Same(other, cell) && other.Row == cell.Row && other.Value == 0 && Contain(val, other.PossibleValues) {
				return false
			}
		}
		return true
	}
	for i, cell := range lst {
		for _, psb := range cell.PossibleValues {
			if excludePiece(psb, cell) || excludeLine(psb, cell) || excludeRow(psb, cell) {
				cell.Value = psb
				lst[i] = cell
				return lst, true
			}
		}
	}
	return lst, false
}

// 判断结果
func Judge(lst []Cell) bool {
	for _, cell := range lst {
		if cell.Value == 0 {
			return false
		}
	}
	return true
}

// NoResult 测试数独是否有解
func NoResult(lst []Cell) bool {
	for _, cell := range lst {
		if cell.Value == 0 && len(cell.PossibleValues) == 0 {
			return true
		}
	}
	return false
}

func Contain(val int, lst []int) bool {
	for _, v := range lst {
		if v == val {
			return true
		}
	}
	return false
}

func StringToInt(val string) int {
	val = strings.Replace(val, " ", "", -1)
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		log.Println(val, err.Error())
		return 0
	}
	return int(i)
}

func Same(c1, c2 Cell) bool {
	return c1.Line == c2.Line && c1.Row == c2.Row
}

func PrintResult(lst []Cell) {
	for i, cell := range lst {
		if i != 0 && i%9 == 0 {
			println()
		}
		if i != 0 && i%27 == 0 {
			println()
		}
		if i%3 == 0 {
			print("    ", cell.Value)
		} else {
			print(cell.Value)
		}
	}
}

