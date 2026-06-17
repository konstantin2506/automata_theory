package game

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gridW     = 30
	gridH     = 30
	btnHeight = 80
)

var dirs = [6][2]int{
	{1, 0},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{0, -1},
	{1, -1},
}

type Game struct {
	maze     [][]bool
	start    [2]int
	exit     [2]int
	robotPos [2]int
	robotDir int
	alive    bool
	won      bool
	message  string

	windowW   int
	windowH   int
	originX   float64
	originY   float64
	hexRadius float64

	img     *image.RGBA
	buttons []button

	moveCh    chan struct{}
	rotateCh  chan int
	surrCh    chan struct{}
	surrResCh chan [3][][2]bool
}

type button struct {
	x, y, w, h int
	label      string
	action     func()
}

func NewGame() *Game {
	g := &Game{}
	g.moveCh = make(chan struct{})
	g.surrCh = make(chan struct{})
	g.rotateCh = make(chan int)
	g.surrResCh = make(chan [3][][2]bool)

	g.windowW, g.windowH = 1024, 768
	g.resize(g.windowW, g.windowH)
	g.regenerate()
	return g
}

func (g *Game) regenerate() {
	rand.Seed(time.Now().UnixNano())
	g.maze, g.start, g.exit = generateMaze(gridW, gridH)
	g.robotPos = g.start
	g.robotDir = 0
	g.alive = true
	g.won = false
	g.message = ""
}

func (g *Game) resize(width, height int) {
	g.windowW = width
	g.windowH = height

	mazeAreaHeight := height - btnHeight
	if mazeAreaHeight < 100 {
		mazeAreaHeight = 100
	}

	reqW := math.Sqrt(3)*(float64(gridW-1)+float64(gridH-1)*0.5) + 2
	reqH := 1.5*float64(gridH-1) + 2
	scaleW := float64(width-20) / reqW
	scaleH := float64(mazeAreaHeight-20) / reqH
	scale := math.Min(scaleW, scaleH)
	if scale < 5 {
		scale = 5
	}
	g.hexRadius = math.Floor(scale)

	totalW := (math.Sqrt(3) * (float64(gridW-1) + float64(gridH-1)*0.5)) * g.hexRadius
	totalH := (1.5 * float64(gridH-1)) * g.hexRadius
	g.originX = float64(width)*0.5 - totalW*0.5 + g.hexRadius
	g.originY = float64(mazeAreaHeight)*0.5 - totalH*0.5 + g.hexRadius

	// Кнопки
	btnW := 120
	btnH := 45
	spacing := 20
	totalBtnW := 4*btnW + 3*spacing
	startX := (width - totalBtnW) / 2
	btnY := height - btnHeight + (btnHeight-btnH)/2

	g.buttons = []button{
		{startX, btnY, btnW, btnH, "Move", func() { g.Move(1) }},
		{startX + btnW + spacing, btnY, btnW, btnH, "Rot L", func() { g.Rotate(-1) }},
		{startX + 2*(btnW+spacing), btnY, btnW, btnH, "Rot R", func() { g.Rotate(1) }},
		{startX + 3*(btnW+spacing), btnY, btnW, btnH, "Reset", func() { g.regenerate() }},
	}

	g.img = image.NewRGBA(image.Rect(0, 0, width, height))
}

func (g *Game) Move(n int) error {
	if !g.alive || g.won {
		return errors.New("игра уже окончена")
	}
	for i := 0; i < n; i++ {
		next := [2]int{g.robotPos[0] + dirs[g.robotDir][0], g.robotPos[1] + dirs[g.robotDir][1]}
		if next[0] < 0 || next[1] < 0 || next[0] >= gridW || next[1] >= gridH || g.maze[next[1]][next[0]] {
			g.alive = false
			g.message = "Робот разбился!"
			return errors.New("стена")
		}
		g.robotPos = next
		if g.robotPos == g.exit {
			g.won = true
			g.message = "Выход найден!"
			return nil
		}
	}
	fmt.Println("Surroundings: ", g.Surroundings()[0][0][0], g.Surroundings()[1][0][0], g.Surroundings()[2][0][0])
	return nil
}

func (g *Game) Rotate(sectors int) {
	if !g.alive || g.won {
		return
	}
	g.robotDir = (g.robotDir + sectors) % 6
	if g.robotDir < 0 {
		g.robotDir += 6
	}

	fmt.Println("Surroundings: ", g.Surroundings()[0][0][0], g.Surroundings()[1][0][0], g.Surroundings()[2][0][0])
}

func (g *Game) Surroundings() [3][][2]bool {
	res := [3][][2]bool{}
	directions := [3]int{
		(g.robotDir - 1 + 6) % 6,
		g.robotDir,
		(g.robotDir + 1) % 6,
	}
	for i, d := range directions {
		view := make([][2]bool, 5)
		for dist := 1; dist <= 5; dist++ {
			q := g.robotPos[0] + dirs[d][0]*dist
			r := g.robotPos[1] + dirs[d][1]*dist
			if q < 0 || r < 0 || q >= gridW || r >= gridH || g.maze[r][q] {
				view[dist-1] = [2]bool{true, false}
			} else {
				isExit := q == g.exit[0] && r == g.exit[1]
				view[dist-1] = [2]bool{false, isExit}
			}
		}
		res[i] = view
	}
	return res
}

/*
	func (g *Game) Surroundings() [3][][2]int {
		res := [3][][2]int{}
		directions := [3]int{
			(g.robotDir - 1 + 6) % 6,
			g.robotDir,
			(g.robotDir + 1) % 6,
		}
		for i, d := range directions {
			view := make([][2]int, 5)
			for dist := 1; dist <= 5; dist++ {
				q := g.robotPos[0] + dirs[d][0]*dist
				r := g.robotPos[1] + dirs[d][1]*dist
				if q < 0 || r < 0 || q >= gridW || r >= gridH || g.maze[r][q] {
					view[dist-1] = [2]int{1, 0}
				} else {
					exitFlag := 0
					if q == g.exit[0] && r == g.exit[1] {
						exitFlag = 1
					}
					view[dist-1] = [2]int{0, exitFlag}
				}
			}
			res[i] = view
		}
		return res
	}
*/
func generateMaze(w, h int) ([][]bool, [2]int, [2]int) {
	maze := make([][]bool, h)
	for r := range maze {
		maze[r] = make([]bool, w)
		for q := range maze[r] {
			maze[r][q] = true
		}
	}
	startQ := rand.Intn(w-2) + 1
	startR := rand.Intn(h-2) + 1
	start := [2]int{startQ, startR}
	maze[startR][startQ] = false

	type cell struct{ q, r int }
	frontier := []cell{}

	addFrontier := func(q, r int) {
		if q < 1 || r < 1 || q >= w-1 || r >= h-1 || !maze[r][q] {
			return
		}
		count := 0
		for _, d := range dirs {
			nq, nr := q+d[0], r+d[1]
			if nq >= 0 && nr >= 0 && nq < w && nr < h && !maze[nr][nq] {
				count++
			}
		}
		if count == 1 {
			for _, f := range frontier {
				if f.q == q && f.r == r {
					return
				}
			}
			frontier = append(frontier, cell{q, r})
		}
	}

	for _, d := range dirs {
		addFrontier(startQ+d[0], startR+d[1])
	}

	for len(frontier) > 0 {
		idx := rand.Intn(len(frontier))
		c := frontier[idx]
		frontier = append(frontier[:idx], frontier[idx+1:]...)

		openCount := 0
		for _, d := range dirs {
			nq, nr := c.q+d[0], c.r+d[1]
			if nq >= 0 && nr >= 0 && nq < w && nr < h && !maze[nr][nq] {
				openCount++
			}
		}
		if openCount == 1 {
			maze[c.r][c.q] = false
			for _, d := range dirs {
				addFrontier(c.q+d[0], c.r+d[1])
			}
		}
	}

	var openCells []cell
	for r := 1; r < h-1; r++ {
		for q := 1; q < w-1; q++ {
			if !maze[r][q] && (q != startQ || r != startR) {
				openCells = append(openCells, cell{q, r})
			}
		}
	}
	exit := [2]int{}
	if len(openCells) == 0 {
		for _, d := range dirs {
			nq, nr := startQ+d[0], startR+d[1]
			if nq >= 1 && nr >= 1 && nq < w-1 && nr < h-1 {
				maze[nr][nq] = false
				exit = [2]int{nq, nr}
				break
			}
		}
	} else {
		c := openCells[rand.Intn(len(openCells))]
		exit = [2]int{c.q, c.r}
	}
	return maze, start, exit
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != g.windowW || outsideHeight != g.windowH {
		g.resize(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, btn := range g.buttons {
			if mx >= btn.x && mx <= btn.x+btn.w && my >= btn.y && my <= btn.y+btn.h {
				btn.action()
				break
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.regenerate()
		return nil
	}
	if !g.alive || g.won {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Move(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Rotate(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Rotate(1)
	}
	/*
		select {
		case _, ok := <-g.moveCh:
			if !ok {
				return nil
			}
			g.Move(1)
		case sectors, ok := <-g.rotateCh:
			if !ok {
				return nil
			}
			g.Rotate(sectors)
		case _, ok := <-g.surrCh:
			if !ok {
				return nil
			}
			fmt.Println(g.Surroundings())
		default:
			return nil
		}
	*/
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bg := color.RGBA{30, 30, 30, 255}
	for y := 0; y < g.windowH; y++ {
		for x := 0; x < g.windowW; x++ {
			g.img.Set(x, y, bg)
		}
	}

	r := g.hexRadius

	for row := 0; row < gridH; row++ {
		for col := 0; col < gridW; col++ {
			cx := g.originX + r*(math.Sqrt(3)*float64(col)+math.Sqrt(3)/2*float64(row))
			cy := g.originY + r*1.5*float64(row)

			var colColor color.Color
			if g.maze[row][col] {
				colColor = color.RGBA{15, 15, 15, 255}
			} else if col == g.start[0] && row == g.start[1] {
				colColor = color.RGBA{80, 220, 80, 255}
			} else if col == g.exit[0] && row == g.exit[1] {
				colColor = color.RGBA{255, 90, 90, 255}
			} else {
				colColor = color.RGBA{220, 220, 220, 255}
			}
			drawFilledHex(g.img, cx, cy, r, colColor)
		}
	}

	robotCX := g.originX + r*(math.Sqrt(3)*float64(g.robotPos[0])+math.Sqrt(3)/2*float64(g.robotPos[1]))
	robotCY := g.originY + r*1.5*float64(g.robotPos[1])
	robotColor := color.RGBA{60, 120, 255, 255}
	drawFilledCircle(g.img, robotCX, robotCY, r*0.45, robotColor)

	if g.alive || g.won {
		dir := dirs[g.robotDir]
		nx := g.originX + r*(math.Sqrt(3)*float64(g.robotPos[0]+dir[0])+math.Sqrt(3)/2*float64(g.robotPos[1]+dir[1]))
		ny := g.originY + r*1.5*float64(g.robotPos[1]+dir[1])
		dx := nx - robotCX
		dy := ny - robotCY
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 0 {
			scale := r * 0.5 / length
			ex := robotCX + dx*scale
			ey := robotCY + dy*scale
			drawLine(g.img, int(robotCX), int(robotCY), int(ex), int(ey), color.RGBA{255, 255, 0, 255})
		}
	}

	for _, btn := range g.buttons {
		btnCol := color.RGBA{70, 70, 70, 255}
		for y := btn.y; y < btn.y+btn.h; y++ {
			for x := btn.x; x < btn.x+btn.w; x++ {
				g.img.Set(x, y, btnCol)
			}
		}
		white := color.RGBA{255, 255, 255, 255}
		for y := btn.y; y <= btn.y+btn.h; y++ {
			g.img.Set(btn.x, y, white)
			g.img.Set(btn.x+btn.w, y, white)
		}
		for x := btn.x; x <= btn.x+btn.w; x++ {
			g.img.Set(x, btn.y, white)
			g.img.Set(x, btn.y+btn.h, white)
		}
	}

	screen.ReplacePixels(g.img.Pix)

	status := "Жив"
	if g.won {
		status = "Победил!"
	} else if !g.alive {
		status = "Разбит"
	}
	dirNames := []string{"E", "NE", "NW", "W", "SW", "SE"}
	msg := "Статус: " + status + " | Направление: " + dirNames[g.robotDir]
	if g.message != "" {
		msg += "\n" + g.message
	}
	// msg += "\nУправление: [↑] Move  [←] Rot -1  [→] Rot +1  [R] рестарт"
	// ebitenutil.DebugPrint(screen, msg)

	for _, btn := range g.buttons {
		ebitenutil.DebugPrintAt(screen, btn.label, btn.x+10, btn.y+btn.h/2-8)
	}
}

// Примитивы рисования (гексагон, круг, линия)
func drawFilledHex(img *image.RGBA, cx, cy, r float64, col color.Color) {
	const n = 6
	verts := make([][2]float64, n)
	for i := 0; i < n; i++ {
		angle := math.Pi/6 + float64(i)*math.Pi/3
		verts[i][0] = cx + r*math.Cos(angle)
		verts[i][1] = cy + r*math.Sin(angle)
	}
	minX := int(cx - r)
	maxX := int(cx + r)
	minY := int(cy - r)
	maxY := int(cy + r)
	bounds := img.Bounds()
	for y := minY; y <= maxY; y++ {
		if y < bounds.Min.Y || y >= bounds.Max.Y {
			continue
		}
		for x := minX; x <= maxX; x++ {
			if x < bounds.Min.X || x >= bounds.Max.X {
				continue
			}
			if pointInPolygon(float64(x)+0.5, float64(y)+0.5, verts) {
				img.Set(x, y, col)
			}
		}
	}
}

func drawFilledCircle(img *image.RGBA, cx, cy, r float64, col color.Color) {
	minX := int(cx - r)
	maxX := int(cx + r)
	minY := int(cy - r)
	maxY := int(cy + r)
	bounds := img.Bounds()
	for y := minY; y <= maxY; y++ {
		if y < bounds.Min.Y || y >= bounds.Max.Y {
			continue
		}
		for x := minX; x <= maxX; x++ {
			if x < bounds.Min.X || x >= bounds.Max.X {
				continue
			}
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, col)
			}
		}
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx := 1
	if x1 >= x2 {
		sx = -1
	}
	sy := 1
	if y1 >= y2 {
		sy = -1
	}
	err := dx + dy
	bounds := img.Bounds()
	for {
		if x1 >= bounds.Min.X && x1 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
			img.Set(x1, y1, col)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

func pointInPolygon(x, y float64, poly [][2]float64) bool {
	n := len(poly)
	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
			inside = !inside
		}
		j = i
	}
	return inside
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type Controller struct {
	game *Game
}

func NewController(game *Game) Controller {
	return Controller{game}
}

func (c *Controller) Move(n int) error {
	time.Sleep(75 * time.Millisecond)
	return c.game.Move(n)
}

func (c *Controller) Rotate(n int) {
	c.game.Rotate(n)
}

func (c *Controller) Surroundings() [3][][2]bool {
	return c.game.Surroundings()
}

func get3Walls(surr [3][][2]bool) [3]bool {
	return [3]bool{surr[0][0][0], surr[1][0][0], surr[2][0][0]}
}

func main() {
	ebiten.SetWindowTitle("Шестиугольный лабиринт 30x30")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)
	game := NewGame()
	go alg(game)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func alg(game *Game) {
	ctrl := NewController(game)
	var err error = nil
	for err == nil {
		walls := get3Walls(game.Surroundings())
		if !walls[2] {
			ctrl.Rotate(1)
			err = ctrl.Move(1)
			continue
		}
		if !walls[1] {
			err = ctrl.Move(1)
			continue
		}
		if !walls[0] {
			ctrl.Rotate(-1)
			err = ctrl.Move(1)
			continue
		}
		ctrl.Rotate(3)
	}
}
