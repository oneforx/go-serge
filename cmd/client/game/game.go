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
}

func (g *SquirrelGame) Update() error {
	// Send MOUSE_POSITIOn

	// screenWidth, screenHeight := ebiten.WindowSize()

	x, y := ebiten.CursorPosition()

	bytes, err := engine.MessageToBytes(messages.CS_ORIENTATION(x, y))
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

	return nil
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
			imageOption.GeoM.Translate(x, y)
			screen.DrawImage(img, imageOption)
		}
	}
}

// Layout retourne la taille de l'écran Ebiten
func (g *SquirrelGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
