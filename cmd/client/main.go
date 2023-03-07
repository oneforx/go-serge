package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/oneforx/go-serge/cmd/client/game"
	"github.com/oneforx/go-serge/ecs"
	"github.com/oneforx/go-serge/engine"
	"github.com/oneforx/go-serge/messages"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

var squirrelGame *game.SquirrelGame = &game.SquirrelGame{
	Latence:      0,
	Token:        nil,
	TcpConnected: false,
	World: &ecs.World{
		Id: string(uuid.New().String()),
	},
	Resources: map[string]*ebiten.Image{},
	Entities:  []ecs.Entity{},
	Camera: &game.Camera{
		X: 0,
		Y: 0,
	},
}

func init() {
	bytes, err := os.ReadFile("./cmd/client/token")
	if err != nil {
		fmt.Println("Could not read token file:", err)
	}

	if len(bytes) != 0 {
		fmt.Println("We found a token in token file, setting up game.Token")
		uuid, err := uuid.Parse(string(bytes))
		if err != nil {
			fmt.Println(uuid)
		}
		squirrelGame.Token = &uuid
	} else {
		fmt.Println("Token file is empty, please type: \"CONNECT email password\", and press Enter")
	}

	// Chargez l'image à partir d'un fichier PNG
	file, err := ebitenutil.OpenFile("image.png")
	if err != nil {
		panic(err)
	}

	imgSrc, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	img := ebiten.NewImageFromImage(imgSrc)

	squirrelGame.Resources["squirrel"] = img
}

func main() {
	addr := "localhost:8080"

	tcpConnexion, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Failed to connect to TCP server:", err)
		return
	}

	defer tcpConnexion.Close()

	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Failed to connect to UDP server:", err)
		return
	}

	remoteUdpConnexion, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Println("Error connecting to UDP server:", err)
		return
	}

	defer remoteUdpConnexion.Close()

	squirrelGame.ClientServer = game.ClientServer{
		TCP: &tcpConnexion,
		UDP: remoteUdpConnexion,
	}

	var wait_group sync.WaitGroup
	var mutex sync.Mutex

	var tcpConnexionDone chan<- error

	wait_group.Add(3)
	go func() {
		defer wait_group.Done()
		// Après que nous écoutons les messages, nous pouvons envoyer un message

		if squirrelGame.Token != nil {
			bytes, err := engine.MessageToBytes(messages.CS_CONNECT_TOKEN(squirrelGame.Token.String()))
			if err != nil {
				log.Println("Could not parse message to bytes")
			}
			tcpConnexion.Write(bytes)
		}

		for {
			buf := make([]byte, 20480000)
			n, err := tcpConnexion.Read(buf)
			if err != nil {
				continue
			}
			var message engine.Message

			if err := json.Unmarshal(buf[:n], &message); err != nil {
				log.Println("Could not read message", err.Error())
				continue
			}

			log.Println("[TCP]: Received: ")

			switch message.MessageType {
			case "CREATE_WORLD":
				type CustomMessage struct {
					MessageType string              `json:"message"`
					Data        []ecs.EntityNoCycle `json:"data"`
				}

				var createWorldMessage CustomMessage

				if err := json.Unmarshal(buf[:n], &createWorldMessage); err != nil {
					continue
				}

				squirrelGame.World.AddEntitiesFromEntitiesNoCycle(createWorldMessage.Data)
			case "CONNECT_SUCCESS":
				token, ok := message.Data.(string)
				if !ok {
					log.Println("Could not parse message.Data to string")
				}

				tokenId, err := uuid.Parse(token)
				if err != nil {
					log.Println("Could not parse message.Data to uuid.UUID")
				}

				if err == nil {
					squirrelGame.Token = &tokenId

					bytesForUdp, err := engine.MessageToBytes(messages.CS_CONNECT_TOKEN(squirrelGame.Token.String()))
					if err != nil {
						log.Println("Could not parse message to bytes")
					}
					remoteUdpConnexion.Write(bytesForUdp)

					squirrelGame.Save()
				}
			case "CREATE_ENTITY":
				type CustomMessage struct {
					MessageType string            `json:"message"`
					Data        ecs.EntityNoCycle `json:"data"`
					Target      engine.TargetType
					NetMode     engine.NET_MODE
				}

				var customMessage CustomMessage

				if err := json.Unmarshal(buf[:n], &customMessage); err != nil {
					log.Println("Could not parse custom message", err)
					continue
				}

				var sync_mutex sync.Mutex

				sync_mutex.Lock()
				squirrelGame.World.AddEntity(ecs.EntityNoCycleToEntity(&squirrelGame.World, customMessage.Data))
				log.Println("dazdazdazdaz")
				sync_mutex.Unlock()
			case "PING":
				latence, ok := message.Data.(float64)
				if !ok {
					log.Println("Could not parse message data to int")
					continue
				}
				var sync_mutex sync.Mutex

				sync_mutex.Lock()
				squirrelGame.Latence = latence
				sync_mutex.Unlock()

				bytes, err := engine.MessageToBytes(messages.CS_PONG())
				if err != nil {
					continue
				}

				_, err = tcpConnexion.Write(bytes)
				if err != nil {
					log.Println("Could not send PONG to server")
				}
			}
		}
	}()

	go func() {
		defer wait_group.Done()

		if squirrelGame.Token != nil {
			bytes, err := engine.MessageToBytes(messages.CS_CONNECT_TOKEN(squirrelGame.Token.String()))
			if err != nil {
				log.Println("Could not parse message to bytes")
			}
			remoteUdpConnexion.Write(bytes)
		}

		for i := 0; i < 10; i++ {
			go func() {
				for {
					buf := make([]byte, 20480)
					mutex.Lock()
					n, _, err := remoteUdpConnexion.ReadFromUDP(buf)
					mutex.Unlock()
					if err != nil {
						log.Println("Could not read message from remote udp connexion")
						continue
					}

					message, err := engine.BytesToMessage(buf[:n])
					if err != nil {
						continue
					}

					if message.MessageType == "UPDATE_ENTITY" {

						type CustomUpdateMessage struct {
							MessageType string                     `json:"message"`
							Data        messages.UpdateMessageData `json:"data"`
						}

						var customMessage CustomUpdateMessage

						if err := json.Unmarshal(buf[:n], &customMessage); err != nil {
							log.Println("Could not parse custom message", err)
							continue
						}

						entity := squirrelGame.World.GetEntity(customMessage.Data.Id)

						if entity != nil {
							entityLocalised := *entity

							entityLocalised.UpdateComponents(customMessage.Data.Components)
							// squirrelGame.World.UpdateEntityComponents(entityLocalised.GetId(), customMessage.Data.Components)
							// log.Println(entityLocalised.GetComponent("position"))
						}
						// sync_mutex.Lock()
						// squirrelGame.Entities = customMessage.Data
						// sync_mutex.Unlock()
					}
				}
			}()
		}
	}()

	go func() {
		defer wait_group.Done()

		scanner := bufio.NewScanner(os.Stdin)

		// Lisez l'entrée utilisateur ligne par ligne
		for {
			if !scanner.Scan() {
				break
			}

			mutex.Lock()
			if tcpConnexionDone != nil {
				break
			}
			mutex.Unlock()

			// Afficher la saisie de l'utilisateur
			fmt.Println("Vous avez tapé :", scanner.Text())

			scannerText := strings.Split(scanner.Text(), " ")
			if scannerText[0] == "DISCONNECT" {
				var message engine.Message = engine.Message{
					MessageType: "DISCONNECT",
				}

				if err := SendTCPMessage(&tcpConnexion, message); err != nil {
					log.Println(err.Error())
					break
				}
			} else if scannerText[0] == "UDP" {
				var message engine.Message = engine.Message{
					MessageType: "UDP",
				}

				if err := SendUDPMessage(remoteUdpConnexion, message); err != nil {
					log.Println(err.Error())
					break
				}
			} else if scannerText[0] == "CONNECT" {
				email := scannerText[1]
				password := scannerText[2]
				var message engine.Message = engine.Message{
					MessageType: "CONNECT",
					Data:        map[string]interface{}{"email": email, "password": password},
				}
				if err := SendTCPMessage(&tcpConnexion, message); err != nil {
					log.Println(err.Error())
					break
				}
				if err := SendUDPMessage(remoteUdpConnexion, message); err != nil {
					log.Println(err.Error())
					break
				}
			} else if scannerText[0] == "CREATE" {
			} else if scannerText[0] == "CONNECT_WORLD" {
				bytes, err := engine.MessageToBytes(messages.CS_CONNECT_WORLD())
				if err != nil {
					log.Println("Could parse message CS_CONNECT_WORLD")
				}
				_, err = tcpConnexion.Write(bytes)
				if err != nil {
					log.Println("Could not send message CS_CONNECT_WORLD")
				}
			} else if scannerText[0] == "DISCONNECT_WORLD" {
				bytes, err := engine.MessageToBytes(messages.CS_DISCONNECT_WORLD())
				if err != nil {
					log.Println("Could parse message CS_DISCONNECT_WORLD")
				}
				_, err = tcpConnexion.Write(bytes)
				if err != nil {
					log.Println("Could not send message CS_DISCONNECT_WORLD")
				}
			} else {
				var message engine.Message = engine.Message{
					MessageType: "MESSAGE",
					Data:        scannerText,
				}

				if err := SendTCPMessage(&tcpConnexion, message); err != nil {
					log.Println(err.Error())
					break
				}
			}
		}

		// Vérifiez les erreurs de scan
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "erreur de lecture de l'entrée standard:", err)
			os.Exit(1)
		}
	}()

	// Création de la fenêtre de l'application SquirrelGame
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Squirrel Game")
	ebiten.SetTPS(120)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	if err := ebiten.RunGame(squirrelGame); err != nil {
		panic(err)
	}

	wait_group.Wait()
}

func SendTCPMessage(tcpConnexionPointer *net.Conn, message engine.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message from user keyboard entry")
	}

	tcpConnexion := *tcpConnexionPointer

	_, err = tcpConnexion.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("could not write message to tcpConnexion")
	}

	return nil
}

func SendUDPMessage(udpConnexionPointer *net.UDPConn, message engine.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message from user keyboard entry")
	}

	udpConnexion := *udpConnexionPointer

	_, err = udpConnexion.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("could not write message to tcpConnexion")
	}

	return nil
}
