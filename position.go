package main

type position_t struct {
	lstart, lend, cstart, cend int
}

func Position(lstart, cstart int) *position_t {
	pos := new(position_t)
	pos.lstart = lstart
	pos.lend = lstart
	pos.cstart = cstart
	pos.cend = cstart
	return pos
}

func PositionFromStartAndEnd(start, end *position_t) *position_t {
	pos := new(position_t)
	pos.lstart = start.lstart
	pos.lend = end.lend
	pos.cstart = start.cstart
	pos.cend = end.cend
	return pos
}

func (p *position_t) merge(other *position_t) *position_t {
	return PositionFromStartAndEnd(p, other)
}