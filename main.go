package main

import (
	"fmt"
	"math/rand"
)

type MazeFieldType uint

const (
	Undef MazeFieldType = iota
	Empty
	Wall
)

type MazeField struct {
	FieldType          MazeFieldType
	PartOfWinningRoute bool
}

type Maze struct {
	Width  uint
	Height uint

	Start struct {
		x uint
		y uint
	}

	Target struct {
		x uint
		y uint
	}

	Fields []MazeField
}

func main() {
	maze := GenerateSolvableMaze(60, 20)

	maze.print()
}

func GenerateSolvableMaze(width uint, height uint) Maze {
	maze := Maze{
		Width:  width,
		Height: height,
		Fields: make([]MazeField, width*height),
		Start: struct {
			x uint
			y uint
		}{x: 0, y: 0},
		Target: struct {
			x uint
			y uint
		}{x: width - 1, y: height - 1},
	}

	iter := rand.Perm(len(maze.Fields))

	for j := range iter {
		i := iter[j]
		for !maze.solvable() || maze.Fields[i].FieldType == Undef {
			fieldType := MazeFieldType(uint(rand.Intn(2) + 1))

			maze.Fields[i] = MazeField{FieldType: fieldType}
		}
	}

	maze.Fields[0] = MazeField{FieldType: Empty, PartOfWinningRoute: true}
	maze.Fields[len(maze.Fields)-1] = MazeField{FieldType: Empty, PartOfWinningRoute: true}

	return maze
}

func (field MazeFieldType) String() string {
	switch field {
	case Undef:
		return "X"

	case Empty:
		return " "

	case Wall:
		return "█"
	}

	panic(fmt.Sprintf("Invalid field type %d", field))
}

func (field MazeField) String() string {
	return field.FieldType.String()
}

func (maze Maze) print() {
	for i := uint(0); i < maze.Height; i++ {
		for j := uint(0); j < maze.Width; j++ {
			fmt.Printf("%s", maze.Fields[j+i*maze.Width])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (maze Maze) fieldAt(x uint, y uint) *MazeField {
	return &maze.Fields[x+(maze.Width*y)]
}

func (maze Maze) solvable() bool {
	visited := map[uint]map[uint]struct{}{}
	visited[maze.Start.x] = map[uint]struct{}{}
	visited[maze.Start.x][maze.Start.y] = struct{}{}

	return maze.solvableFrom(maze.Start.x, maze.Start.y, visited)
}

type coords struct {
	x uint
	y uint
}

func (maze Maze) solvableFrom(x uint, y uint, visited map[uint]map[uint]struct{}) bool {
	if x == maze.Target.x && y == maze.Target.y {
		return true
	} else {
		visited[x][y] = struct{}{}

		maybePairs := []coords{
			{x: x + 1, y: y},
			{x: x - 1, y: y},
			{x: x, y: y + 1},
			{x: x, y: y - 1},
		}

		var pairs []coords
		for j := range maybePairs {
			if maybePairs[j].x >= 0 && maybePairs[j].y >= 0 && maybePairs[j].x < maze.Width && maybePairs[j].y < maze.Height {
				_, pairVisited := visited[maybePairs[j].x][maybePairs[j].y]
				field := maze.fieldAt(maybePairs[j].x, maybePairs[j].y)
				if pairVisited == false && (field.FieldType == Empty || field.FieldType == Undef) {
					pairs = append(pairs, maybePairs[j])
				}
			}
		}

		anySolvable := false
		for i := range pairs {
			_, ok := visited[pairs[i].x]
			if ok == false {
				visited[pairs[i].x] = map[uint]struct{}{}
			}
			visited[pairs[i].x][pairs[i].y] = struct{}{}
		}

		for i := range pairs {
			if maze.solvableFrom(pairs[i].x, pairs[i].y, visited) {
				anySolvable = true
			} else {
			}
		}

		return anySolvable
	}
}
