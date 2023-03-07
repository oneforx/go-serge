package game

import (
	"image/color"
	"log"
	"net"
	"os"

	"github.com/bawdeveloppement/squirrelgameserver/ecs"
	"github.com/bawdeveloppement/squirrelgameserver/engine"
	"github.com/bawdeveloppement/squirrelgameserver/messages"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ClientServer struct {
	TCP *net.Conn
	UDP *net.UDPConn
}

type SquirrelGame struct {
	Latence      float64
	AccountID    *uuid.UUID
	Token        *uuid.UUID
	World        ecs.IWorld
	TcpConnected bool
	Resources    map[string]*ebiten.Image
	PressedKeys  []ebiten.Key
	ClientServer ClientServer
	Entities     []ecs.Entity
	Camera       *Camera
}

func (g *SquirrelGame) Update() error {

	if len(g.World.GetEntities()) > 0 {
		log.Println(len(g.World.GetEntities()))
	}

	width, height := ebiten.WindowSize()

	if entities := g.World.GetEntitiesPossessedBy(*g.Token); len(entities) > 0 {
		player := *entities[0]
		playerPosition, ok := player.GetComponent("position").Data.(map[string]interface{})
		if !ok {
			log.Println("Could not parse player position to map[string]interface{}")
		}
		x, ok := playerPosition["x"].(float64)
		if !ok {
			log.Println("Could not parse x to float64")
		}
		y, ok := playerPosition["y"].(float64)
		if !ok {
			log.Println("Could not parse y to float64")
		}
		g.Camera.X = int(x) - width/2
		g.Camera.Y = int(y) - height/2
	} else {
		g.Camera.X = 0
		g.Camera.Y = 0
	}
	// Si on a une entité que l'on posséde
	if entitiesWeOwn := g.World.GetEntitiesPossessedBy(*g.Token); len(entitiesWeOwn) > 0 {
		player := *entitiesWeOwn[0]
		playerPosition, ok := player.GetComponent("position").Data.(map[string]interface{})
		if !ok {
			log.Println("Could not parse player position to map[string]interface{}")
		}
		_, ok = playerPosition["x"].(float64)
		if !ok {
			log.Println("Could not parse x to float64")
		}
		_, ok = playerPosition["y"].(float64)
		if !ok {
			log.Println("Could not parse y to float64")
		}

		mx, my := ebiten.CursorPosition()

		bytes, err := engine.MessageToBytes(messages.CS_ORIENTATION(mx, my))
		if err != nil {
			log.Println("COuld not parse ")
		}

		_, err = g.ClientServer.UDP.Write(bytes)
		if err != nil {
			log.Println("Could not send CS_ORIENTATION")
		}

		if g.PressedKeys = inpututil.AppendPressedKeys(g.PressedKeys[:0]); len(g.PressedKeys) != 0 {
			for _, k := range g.PressedKeys {
				switch k.String() {
				case "D":
					move_right_message, err := engine.MessageToBytes(engine.Message{
						MessageType: "MOVE",
						Data:        "RIGHT",
					})
					if err != nil {
						log.Println("Could not parse message")
					}

					if err == nil {
						g.ClientServer.UDP.Write(move_right_message)
					}
				case "S":
					move_right_message, err := engine.MessageToBytes(engine.Message{
						MessageType: "MOVE",
						Data:        "DOWN",
					})
					if err != nil {
						log.Println("Could not parse message")
					}

					if err == nil {
						g.ClientServer.UDP.Write(move_right_message)
					}
				case "W":
					move_right_message, err := engine.MessageToBytes(engine.Message{
						MessageType: "MOVE",
						Data:        "UP",
					})
					if err != nil {
						log.Println("Could not parse message")
					}

					if err == nil {
						g.ClientServer.UDP.Write(move_right_message)
					}
				case "A":
					move_left_message, err := engine.MessageToBytes(engine.Message{
						MessageType: "MOVE",
						Data:        "LEFT",
					})
					if err != nil {
						log.Println("Could not parse message")
					}

					if err == nil {
						g.ClientServer.UDP.Write(move_left_message)
					}
				}
			}
		}

		leftButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		if leftButtonPressed {
			tcpConnexion := *g.ClientServer.TCP

			bytes, err := engine.MessageToBytes(messages.CS_SHOOT())
			if err != nil {
				log.Println("Could not parse message CS_SHOOT to bytes")
			}

			tcpConnexion.Write(bytes)
		}
	}

	return nil
}
func GetMouseWorldPosition(mouseX, mouseY int, camera Camera, screen map[string]int, world map[string]int) (int, int) {
	worldX := camera.X + (mouseX/screen["width"])*world["width"]
	worldY := camera.Y + (mouseY/screen["height"])*world["height"]
	return worldX, worldY
}

func (g *SquirrelGame) Save() error {
	if g.Token != nil {
		bytes := []byte(g.Token.String())
		if err := os.WriteFile("./cmd/client/token", bytes, 0666); err != nil {
			return err
		}
	}

	return nil
}

type Camera struct {
	X int
	Y int
}

// Draw dessine le jeu Ebiten
func (g *SquirrelGame) Draw(screen *ebiten.Image) {

	if len(g.World.GetEntities()) > 0 {
		for _, entity := range g.World.GetEntities() {
			entityLocalisation := *entity
			positionComponent := entityLocalisation.GetComponent("position")
			if positionComponent == nil {
				log.Println("dazddzdzzdddda")
				continue
			}
			positionComponentLocalised := *positionComponent
			positionComponentData, ok := positionComponentLocalised.GetData().(map[string]interface{})
			if !ok {
				log.Println("daxdsxxzda")
				continue
			}
			x, ok := positionComponentData["x"].(float64)
			if !ok {
				log.Println("daxdsxxzda")
				continue
			}
			y, ok := positionComponentData["y"].(float64)
			if !ok {
				log.Println("daxdsxxzda")
				continue
			}
			img, ok := g.Resources["squirrel"]
			if !ok {
				log.Println("test", img)
				log.Println("skip")
				continue
			}

			imageOption := &ebiten.DrawImageOptions{
				GeoM:          ebiten.GeoM{},
				ColorM:        ebiten.ColorM{},
				CompositeMode: 0,
				Filter:        0,
			}
			ebitenutil.DrawRect(screen, 0, 0, 32, 32, color.RGBA{255, 255, 255, 255})
			imageOption.Filter = ebiten.FilterLinear
			imageOption.GeoM.Translate(x-float64(g.Camera.X), y-float64(g.Camera.Y))
			screen.DrawImage(img, imageOption)
		}
	}
}

// Layout retourne la taille de l'écran Ebiten
func (g *SquirrelGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
