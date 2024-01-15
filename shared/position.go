package shared

import (
	"fmt"
	"os"
)

type Position_t struct {
	Lstart, Lend, Cstart, Cend int
}

func Position(Lstart, Cstart int) *Position_t {
	pos := new(Position_t)
	pos.Lstart = Lstart
	pos.Lend = Lstart
	pos.Cstart = Cstart
	pos.Cend = Cstart
	return pos
}

func PositionFromStartAndEnd(start, end *Position_t) *Position_t {
	pos := new(Position_t)

	//lint:ignore SA4031 possible not nil
	if pos == nil {
		fmt.Println("Out of memory!!!")
		os.Exit(0x1)
	}

	pos.Lstart = start.Lstart
	pos.Lend = end.Lend
	pos.Cstart = start.Cstart
	pos.Cend = end.Cend
	return pos
}

func (p *Position_t) Merge(other *Position_t) *Position_t {
	return PositionFromStartAndEnd(p, other)
}
