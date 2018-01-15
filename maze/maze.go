package maze

import (
    "container/list"
    "math/rand"
    "github.com/go-gl/mathgl/mgl32"
    "github.com/leedenison/gologo"
)

var player *gologo.Object

var DEFAULT_SCREEN_SIZE_FACTOR = float32(0.8)
var DEFAULT_WALL_WIDTH_FACTOR = float32(0.1)
var DEFAULT_ROOM_BORDER_FACTOR = float32(0.2)

type Maze struct {
    Size [2]int
    Start [2]int
    End [2]int
    HWalls [][]*gologo.Object
    VWalls [][]*gologo.Object
    BottomLeft [2]float32
    RoomSize float32
}

func GenerateMaze(size [2]int) *Maze {
    maze := initializeMaze(size)
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

    return maze
}

func calcRenderSize(maze *Maze) {
    windowSize := gologo.GetWindowSize()
    windowCenter := gologo.GetWindowCenter()
    screenDim := windowSize[0]
    if screenDim > windowSize[1] {
        screenDim = windowSize[1]
    }
    roomDim := maze.Size[0]
    if roomDim < maze.Size[1] {
        roomDim = maze.Size[1]
    }

    maze.RoomSize = screenDim * DEFAULT_SCREEN_SIZE_FACTOR / float32(roomDim)

    maze.BottomLeft = [2]float32 {
        windowCenter[0] - maze.RoomSize * float32(maze.Size[0]) / 2,
        windowCenter[1] - maze.RoomSize * float32(maze.Size[1]) / 2,
    }
}

func initializeMaze(size [2]int) *Maze {
    result := &Maze {
        Size: size,
        Start: [2]int {
            rand.Intn(size[0]),
            0,
        },
        End: [2]int {
            rand.Intn(size[0]),
            size[1] - 1,
        },
    }

    calcRenderSize(result)
    wallWidth := result.RoomSize * DEFAULT_WALL_WIDTH_FACTOR
    halfWallLength := (result.RoomSize + wallWidth) / 2
    halfWallWidth := wallWidth / 2
    horizontalWall := gologo.Rectangle(
      gologo.Rect {
          { -halfWallLength, -halfWallWidth },
          { halfWallLength, halfWallWidth },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })
    verticalWall := gologo.Rectangle(
      gologo.Rect {
          { -halfWallWidth, -halfWallLength },
          { halfWallWidth, halfWallLength },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })

    hWallsSize := [2]int { size[0], size[1] - 1 }
    vWallsSize := [2]int { size[0] - 1, size[1] }

    initializeBorder(result.Size, result.BottomLeft, result.RoomSize, halfWallWidth)
    result.HWalls = initializeWalls(
        hWallsSize,
        horizontalWall,
        result.BottomLeft,
        result.RoomSize,
        [2]float32 { result.RoomSize / 2, result.RoomSize })
    result.VWalls = initializeWalls(
        vWallsSize,
        verticalWall,
        result.BottomLeft,
        result.RoomSize,
        [2]float32 { result.RoomSize, result.RoomSize / 2 })

    player = initializeStart(result)
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
    startBottomLeft := [2]float32 {
        maze.BottomLeft[0] + float32(maze.Start[0]) * maze.RoomSize + roomBorder,
        maze.BottomLeft[1] + float32(maze.Start[1]) * maze.RoomSize + roomBorder,
    }
    start := gologo.Rectangle(
      gologo.Rect {
          startBottomLeft,
          {
              startBottomLeft[0] + maze.RoomSize - 2 * roomBorder,
              startBottomLeft[1] + maze.RoomSize - 2 * roomBorder,
          },
      },
      mgl32.Vec4 { 0.0, 1.0, 0.0, 1.0 })
    gologo.TagRender(start)
    return start
}

func initializeEnd(maze *Maze) {
    roomBorder := maze.RoomSize * DEFAULT_ROOM_BORDER_FACTOR
    endBottomLeft := [2]float32 {
        maze.BottomLeft[0] + float32(maze.End[0]) * maze.RoomSize + roomBorder,
        maze.BottomLeft[1] + float32(maze.End[1]) * maze.RoomSize + roomBorder,
    }
    end := gologo.Rectangle(
      gologo.Rect {
          endBottomLeft,
          {
              endBottomLeft[0] + maze.RoomSize - 2 * roomBorder,
              endBottomLeft[1] + maze.RoomSize - 2 * roomBorder,
          },
      },
      mgl32.Vec4 { 1.0, 0.0, 0.0, 1.0 })
    gologo.TagRender(end)
}

func initializeBorder(mazeSize [2]int, bottomLeft [2]float32, roomSize float32, halfWallWidth float32) {
    bottom := gologo.Rectangle(
      gologo.Rect {
          {
            bottomLeft[0] - halfWallWidth,
            bottomLeft[1] - halfWallWidth,
          },
          {
            bottomLeft[0] + float32(mazeSize[0]) * roomSize + halfWallWidth,
            bottomLeft[1] + halfWallWidth,
          },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })

    top := gologo.Rectangle(
      gologo.Rect {
          {
            bottomLeft[0] - halfWallWidth,
            bottomLeft[1] + float32(mazeSize[1]) * roomSize - halfWallWidth,
          },
          {
            bottomLeft[0] + float32(mazeSize[0]) * roomSize + halfWallWidth,
            bottomLeft[1] + float32(mazeSize[1]) * roomSize + halfWallWidth,
          },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })

    left := gologo.Rectangle(
      gologo.Rect {
          {
            bottomLeft[0] - halfWallWidth,
            bottomLeft[1] - halfWallWidth,
          },
          {
            bottomLeft[0] + halfWallWidth,
            bottomLeft[1] + float32(mazeSize[1]) * roomSize + halfWallWidth,
          },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })

    right := gologo.Rectangle(
      gologo.Rect {
          {
            bottomLeft[0] + float32(mazeSize[0]) * roomSize - halfWallWidth,
            bottomLeft[1] - halfWallWidth,
          },
          {
            bottomLeft[0] + float32(mazeSize[0]) * roomSize + halfWallWidth,
            bottomLeft[1] + float32(mazeSize[1]) * roomSize + halfWallWidth,
          },
      },
      mgl32.Vec4 { 1.0, 1.0, 1.0, 1.0 })

    gologo.TagRender(bottom)
    gologo.TagRender(top)
    gologo.TagRender(left)
    gologo.TagRender(right)
}

func initializeWalls(
        size [2]int,
        wall *gologo.Object,
        mazeOffset [2]float32,
        roomSize float32,
        wallOffset [2]float32) [][]*gologo.Object {
    result := make([][]*gologo.Object, size[0])
    for i := range result {
        result[i] = make([]*gologo.Object, size[1])
        for j := range result[i] {
            result[i][j] = wall.Clone()
            gologo.TagRender(result[i][j])
            result[i][j].Model.SetCol(3, mgl32.Vec4 {
                mazeOffset[0] + float32(i) * roomSize + wallOffset[0],
                mazeOffset[1] + float32(j) * roomSize + wallOffset[1],
                0.0,
                1.0,
            })
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
    if start[1] < maze.Size[1] - 1 &&
        rooms[start[0]][start[1] + 1] == findAdded {
        // Up
        result = append(result, [2]int { start[0], start[1] + 1 })
    }
    if start[1] > 0 &&
        rooms[start[0]][start[1] - 1] == findAdded {
        // Down
        result = append(result, [2]int { start[0], start[1] - 1 })
    }
    if start[0] > 0 &&
        rooms[start[0] - 1][start[1]] == findAdded {
        // Left
        result = append(result, [2]int { start[0] - 1, start[1] })
    }
    if start[0] < maze.Size[0] - 1 &&
        rooms[start[0] + 1][start[1]] == findAdded {
        // Right
        result = append(result, [2]int { start[0] + 1, start[1] })
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
        gologo.UntagRender(maze.VWalls[from[0]][from[1]])
        maze.VWalls[from[0]][from[1]] = nil
    }
    if from[0] > to[0] {
        gologo.UntagRender(maze.VWalls[to[0]][to[1]])
        maze.VWalls[to[0]][to[1]] = nil
    }
    if from[1] < to[1] {
        gologo.UntagRender(maze.HWalls[from[0]][from[1]])
        maze.HWalls[from[0]][from[1]] = nil
    }
    if from[1] > to[1] {
        gologo.UntagRender(maze.HWalls[to[0]][to[1]])
        maze.HWalls[to[0]][to[1]] = nil
    }
}
