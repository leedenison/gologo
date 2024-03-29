package examples

import (
	"container/list"
	"fmt"
	"math/rand"
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo"
	"github.com/leedenison/gologo/obj"
	"github.com/leedenison/gologo/tags"
)

var (
	tagged       = tags.TagSet{}
	callbackMap  = []*Maze{}
	lastCallback = int(0)
)

func mazeTickCallback(tick int) {
	if tick-lastCallback > DEFAULT_TICK_INCREMENT {
		for _, maze := range callbackMap {
			if maze.Callback != nil && !HasRemainingMoves(maze) && !IsFinished(maze) {
				maze.Callback(maze)
			}
			maze.DoMove()
		}
		lastCallback = tick
	}
}

var (
	DEFAULT_SCREEN_SIZE_FACTOR = float32(0.8)
	DEFAULT_WALL_WIDTH_FACTOR  = float32(0.1)
	DEFAULT_ROOM_BORDER_FACTOR = float32(0.2)
)

var DEFAULT_TICK_INCREMENT = 200

type Maze struct {
	Gologo         *gologo.Gologo
	Size           [2]int
	Start          [2]int
	End            [2]int
	HWalls         [][]*gologo.Object
	VWalls         [][]*gologo.Object
	BottomLeft     [2]float32
	RoomSize       float32
	Player         *gologo.Object
	PlayerPosition [2]int
	Callback       func(*Maze)
	MoveQueue      []Direction
	LastMove       int
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

func Move(maze *Maze, direction Direction) {
	maze.MoveQueue = append(maze.MoveQueue, direction)
}

// Moves the player one square up.  Will not move the player if there is a wall
// above the player.
func MoveUp(maze *Maze) {
	maze.MoveQueue = append(maze.MoveQueue, UP)
}

// Moves the player one square down.  Will not move the player if there is a wall
// below the player.
func MoveDown(maze *Maze) {
	maze.MoveQueue = append(maze.MoveQueue, DOWN)
}

// Moves the player one square to the left.  Will not move the player if there is a wall
// to the left.
func MoveLeft(maze *Maze) {
	maze.MoveQueue = append(maze.MoveQueue, LEFT)
}

// Moves the player one square to the right.  Will not move the player if there is a wall
// to the right.
func MoveRight(maze *Maze) {
	maze.MoveQueue = append(maze.MoveQueue, RIGHT)
}

func HasRemainingMoves(maze *Maze) bool {
	return maze.LastMove < len(maze.MoveQueue)
}

func IsFinished(maze *Maze) bool {
	return maze.PlayerPosition == maze.End
}

func GetLastMove(maze *Maze) Direction {
	if maze.LastMove == 0 {
		return UP
	} else {
		return maze.MoveQueue[maze.LastMove-1]
	}
}

func CanMove(maze *Maze, direction Direction) bool {
	pos := maze.PlayerPosition

	switch direction {
	case UP:
		return pos[1] < maze.Size[1]-1 && maze.HWalls[pos[0]][pos[1]] == nil
	case DOWN:
		return pos[1] > 0 && maze.HWalls[pos[0]][pos[1]-1] == nil
	case LEFT:
		return pos[0] > 0 && maze.VWalls[pos[0]-1][pos[1]] == nil
	case RIGHT:
		return pos[0] < maze.Size[0]-1 && maze.VWalls[pos[0]][pos[1]] == nil
	default:
		panic(fmt.Sprintf("Unknown direction: %v\n", direction))
	}
}

func (maze *Maze) DoMove() {
	if HasRemainingMoves(maze) {
		direction := maze.MoveQueue[maze.LastMove]

		if CanMove(maze, direction) {
			switch direction {
			case UP:
				maze.Player.Position = maze.Player.Position.Add(mgl32.Vec3{0, maze.RoomSize, 0})
				maze.PlayerPosition[1] += 1
			case DOWN:
				maze.Player.Position = maze.Player.Position.Add(mgl32.Vec3{0, -maze.RoomSize, 0})
				maze.PlayerPosition[1] -= 1
			case LEFT:
				maze.Player.Position = maze.Player.Position.Add(mgl32.Vec3{-maze.RoomSize, 0, 0})
				maze.PlayerPosition[0] -= 1
			case RIGHT:
				maze.Player.Position = maze.Player.Position.Add(mgl32.Vec3{maze.RoomSize, 0, 0})
				maze.PlayerPosition[0] += 1
			default:
				panic(fmt.Sprintf("Unknown direction: %v\n", direction))
			}
		}
		maze.LastMove++
	}
}

func Run(x, y int, callback func(*Maze)) {
	g := gologo.Init()
	defer g.Close()

	maze := GenerateMaze(g, [2]int{x, y}, callback)
	callbackMap = append(callbackMap, maze)

	for !g.Window.ShouldClose() {
		g.ClearBackBuffer()

		r := tagged.GetAll("render")
		sort.Sort(gologo.ByZOrder(r))
		for _, object := range r {
			object.Draw()
		}

		g.Window.SwapBuffers()
		g.CheckForEvents()
	}
}

func GenerateMaze(g *gologo.Gologo, size [2]int, callback func(*Maze)) *Maze {
	maze := initializeMaze(g, size)
	rooms := initializeRooms(size)
	roomsList := list.New()

	rooms[maze.Start[0]][maze.Start[1]] = true
	roomsList.PushFront(maze.Start)

	for roomsList.Len() > 0 {
		randomRoom := randomRoom(roomsList)
		directions := directionsWithoutRooms(maze, rooms, randomRoom)
		randomDirection := rand.Intn(len(directions))
		extension := directions[randomDirection]

		addRoom(maze, rooms, roomsList, extension)
		removeWall(maze, randomRoom, extension)
	}

	maze.Callback = callback

	return maze
}

func calcRenderSize(maze *Maze) {
	winX, winY := maze.Gologo.Window.GetSize()
	windowCenter := maze.Gologo.GetWindowCenter()
	screenDim := float32(winX)
	if screenDim > float32(winY) {
		screenDim = float32(winY)
	}
	roomDim := maze.Size[0]
	if roomDim < maze.Size[1] {
		roomDim = maze.Size[1]
	}

	maze.RoomSize = screenDim * DEFAULT_SCREEN_SIZE_FACTOR / float32(roomDim)

	maze.BottomLeft = [2]float32{
		windowCenter[0] - maze.RoomSize*float32(maze.Size[0])/2,
		windowCenter[1] - maze.RoomSize*float32(maze.Size[1])/2,
	}
}

func initializeMaze(g *gologo.Gologo, size [2]int) *Maze {
	result := &Maze{
		Gologo: g,
		Size:   size,
		Start: [2]int{
			rand.Intn(size[0]),
			0,
		},
		End: [2]int{
			rand.Intn(size[0]),
			size[1] - 1,
		},
	}

	calcRenderSize(result)
	wallWidth := result.RoomSize * DEFAULT_WALL_WIDTH_FACTOR
	halfWallLength := (result.RoomSize + wallWidth) / 2
	halfWallWidth := wallWidth / 2
	horizontalWall := obj.Rectangle(
		gologo.Rect{
			{-halfWallLength, -halfWallWidth},
			{halfWallLength, halfWallWidth},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})
	verticalWall := obj.Rectangle(
		gologo.Rect{
			{-halfWallWidth, -halfWallLength},
			{halfWallWidth, halfWallLength},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	hWallsSize := [2]int{size[0], size[1] - 1}
	vWallsSize := [2]int{size[0] - 1, size[1]}

	initializeBorder(result.Size, result.BottomLeft, result.RoomSize, halfWallWidth)
	result.HWalls = initializeWalls(
		hWallsSize,
		horizontalWall,
		result.BottomLeft,
		result.RoomSize,
		[2]float32{result.RoomSize / 2, result.RoomSize})
	result.VWalls = initializeWalls(
		vWallsSize,
		verticalWall,
		result.BottomLeft,
		result.RoomSize,
		[2]float32{result.RoomSize, result.RoomSize / 2})

	result.Player = initializeStart(result)
	result.PlayerPosition = result.Start
	initializeEnd(result)

	return result
}

func initializeRooms(size [2]int) [][]bool {
	rooms := make([][]bool, size[0])
	for i := range rooms {
		rooms[i] = make([]bool, size[1])
	}
	return rooms
}

func initializeStart(maze *Maze) *gologo.Object {
	roomBorder := maze.RoomSize * DEFAULT_ROOM_BORDER_FACTOR
	startBottomLeft := [2]float32{
		maze.BottomLeft[0] + float32(maze.Start[0])*maze.RoomSize + roomBorder,
		maze.BottomLeft[1] + float32(maze.Start[1])*maze.RoomSize + roomBorder,
	}
	start := obj.Rectangle(
		gologo.Rect{
			startBottomLeft,
			{
				startBottomLeft[0] + maze.RoomSize - 2*roomBorder,
				startBottomLeft[1] + maze.RoomSize - 2*roomBorder,
			},
		},
		mgl32.Vec4{0.0, 1.0, 0.0, 1.0})
	tagged.Tag(start, "render")
	return start
}

func initializeEnd(maze *Maze) {
	roomBorder := maze.RoomSize * DEFAULT_ROOM_BORDER_FACTOR
	endBottomLeft := [2]float32{
		maze.BottomLeft[0] + float32(maze.End[0])*maze.RoomSize + roomBorder,
		maze.BottomLeft[1] + float32(maze.End[1])*maze.RoomSize + roomBorder,
	}
	end := obj.Rectangle(
		gologo.Rect{
			endBottomLeft,
			{
				endBottomLeft[0] + maze.RoomSize - 2*roomBorder,
				endBottomLeft[1] + maze.RoomSize - 2*roomBorder,
			},
		},
		mgl32.Vec4{1.0, 0.0, 0.0, 1.0})
	tagged.Tag(end, "render")
}

func initializeBorder(mazeSize [2]int, bottomLeft [2]float32, roomSize float32, halfWallWidth float32) {
	bottom := obj.Rectangle(
		gologo.Rect{
			{
				bottomLeft[0] - halfWallWidth,
				bottomLeft[1] - halfWallWidth,
			},
			{
				bottomLeft[0] + float32(mazeSize[0])*roomSize + halfWallWidth,
				bottomLeft[1] + halfWallWidth,
			},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	top := obj.Rectangle(
		gologo.Rect{
			{
				bottomLeft[0] - halfWallWidth,
				bottomLeft[1] + float32(mazeSize[1])*roomSize - halfWallWidth,
			},
			{
				bottomLeft[0] + float32(mazeSize[0])*roomSize + halfWallWidth,
				bottomLeft[1] + float32(mazeSize[1])*roomSize + halfWallWidth,
			},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	left := obj.Rectangle(
		gologo.Rect{
			{
				bottomLeft[0] - halfWallWidth,
				bottomLeft[1] - halfWallWidth,
			},
			{
				bottomLeft[0] + halfWallWidth,
				bottomLeft[1] + float32(mazeSize[1])*roomSize + halfWallWidth,
			},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	right := obj.Rectangle(
		gologo.Rect{
			{
				bottomLeft[0] + float32(mazeSize[0])*roomSize - halfWallWidth,
				bottomLeft[1] - halfWallWidth,
			},
			{
				bottomLeft[0] + float32(mazeSize[0])*roomSize + halfWallWidth,
				bottomLeft[1] + float32(mazeSize[1])*roomSize + halfWallWidth,
			},
		},
		mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	tagged.Tag(bottom, "render")
	tagged.Tag(top, "render")
	tagged.Tag(left, "render")
	tagged.Tag(right, "render")
}

func initializeWalls(
	size [2]int,
	wall *gologo.Object,
	mazeOffset [2]float32,
	roomSize float32,
	wallOffset [2]float32,
) [][]*gologo.Object {
	result := make([][]*gologo.Object, size[0])
	for i := range result {
		result[i] = make([]*gologo.Object, size[1])
		for j := range result[i] {
			result[i][j] = wall.Clone()
			tagged.Tag(result[i][j], "render")
			result[i][j].Position = mgl32.Vec3{
				mazeOffset[0] + float32(i)*roomSize + wallOffset[0],
				mazeOffset[1] + float32(j)*roomSize + wallOffset[1],
				0.0,
			}
		}
	}
	return result
}

func randomRoom(roomsList *list.List) [2]int {
	randomIdx := rand.Intn(roomsList.Len())

	roomElement := roomsList.Front()
	for i := 0; i < randomIdx; i++ {
		roomElement = roomElement.Next()
	}

	return roomElement.Value.([2]int)
}

func directionsWithRooms(maze *Maze, rooms [][]bool, start [2]int) [][2]int {
	return directions(maze, rooms, start, true)
}

func directionsWithoutRooms(maze *Maze, rooms [][]bool, start [2]int) [][2]int {
	return directions(maze, rooms, start, false)
}

func directions(maze *Maze, rooms [][]bool, start [2]int, findAdded bool) [][2]int {
	result := make([][2]int, 0, 4)
	if start[1] < maze.Size[1]-1 &&
		rooms[start[0]][start[1]+1] == findAdded {
		// Up
		result = append(result, [2]int{start[0], start[1] + 1})
	}
	if start[1] > 0 &&
		rooms[start[0]][start[1]-1] == findAdded {
		// Down
		result = append(result, [2]int{start[0], start[1] - 1})
	}
	if start[0] > 0 &&
		rooms[start[0]-1][start[1]] == findAdded {
		// Left
		result = append(result, [2]int{start[0] - 1, start[1]})
	}
	if start[0] < maze.Size[0]-1 &&
		rooms[start[0]+1][start[1]] == findAdded {
		// Right
		result = append(result, [2]int{start[0] + 1, start[1]})
	}

	return result
}

func addRoom(maze *Maze, rooms [][]bool, roomsList *list.List, room [2]int) {
	// Mark the room as added
	rooms[room[0]][room[1]] = true

	// Add the current room if there are extension options
	directions := directionsWithoutRooms(maze, rooms, room)
	if len(directions) > 0 {
		roomsList.PushFront(room)
	}

	// Check if surrounding rooms are now complete
	surroundingRooms := directionsWithRooms(maze, rooms, room)
	for i := range surroundingRooms {
		removeIfCompleted(maze, rooms, roomsList, surroundingRooms[i])
	}
}

func removeIfCompleted(maze *Maze, rooms [][]bool, roomsList *list.List, room [2]int) {
	directions := directionsWithoutRooms(maze, rooms, room)

	if len(directions) == 0 {
		// Room cannot be extended in any direction
		roomElement := roomsList.Front()
		for i := 0; i < roomsList.Len(); i++ {
			candidate := roomElement.Value.([2]int)

			if candidate[0] == room[0] && candidate[1] == room[1] {
				roomsList.Remove(roomElement)
				break
			}

			roomElement = roomElement.Next()
		}
	}
}

func removeWall(maze *Maze, from [2]int, to [2]int) {
	if from[0] < to[0] {
		tagged.Untag(maze.VWalls[from[0]][from[1]], "render")
		maze.VWalls[from[0]][from[1]] = nil
	}
	if from[0] > to[0] {
		tagged.Untag(maze.VWalls[to[0]][to[1]], "render")
		maze.VWalls[to[0]][to[1]] = nil
	}
	if from[1] < to[1] {
		tagged.Untag(maze.HWalls[from[0]][from[1]], "render")
		maze.HWalls[from[0]][from[1]] = nil
	}
	if from[1] > to[1] {
		tagged.Untag(maze.HWalls[to[0]][to[1]], "render")
		maze.HWalls[to[0]][to[1]] = nil
	}
}

func Clockwise(direction Direction) Direction {
	switch direction {
	case UP:
		return RIGHT
	case RIGHT:
		return DOWN
	case DOWN:
		return LEFT
	case LEFT:
		return UP
	default:
		panic(fmt.Sprintf("Unknown direction: %v\n", direction))
	}
}

func AntiClockwise(direction Direction) Direction {
	switch direction {
	case UP:
		return LEFT
	case LEFT:
		return DOWN
	case DOWN:
		return RIGHT
	case RIGHT:
		return UP
	default:
		panic(fmt.Sprintf("Unknown direction: %v\n", direction))
	}
}

func Opposite(direction Direction) Direction {
	switch direction {
	case UP:
		return DOWN
	case LEFT:
		return RIGHT
	case DOWN:
		return UP
	case RIGHT:
		return LEFT
	default:
		panic(fmt.Sprintf("Unknown direction: %v\n", direction))
	}
}
