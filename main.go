package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"plugin"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oneforx/go-ecs"
	"github.com/oneforx/go-serge/engine"
	"github.com/oneforx/go-serge/library"
	"github.com/oneforx/go-serge/messages"
	"github.com/oneforx/go-serge/systems"
)

type Client struct {
	UDP        *net.UDPAddr
	TCP        *net.TCPConn
	Latence    int
	Token      uuid.UUID
	ID         uuid.UUID
	Email      string
	Password   string
	Characters []ecs.Entity
}

// Connexion Mode
type NET_MODE string

var (
	TCP NET_MODE = "TCP"
	UDP NET_MODE = "UDP"
)

var (
	LibraryManager = &library.LibraryManager{
		Libraries: map[ecs.Identifier]ecs.Library{},
	}
)

func init() {

	modDirectory := "./mods"

	err := filepath.Walk(modDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			p, err := plugin.Open(path)
			if err != nil {
				return err
			}

			loadPlugin, err := p.Lookup("Load")
			if err != nil {
				panic(err)
			}

			newLibrary := loadPlugin.(func() ecs.ILibrary)()

			log.Println("LOADING " + newLibrary.GetId().String())
			log.Println(newLibrary.GetSystems())
			LibraryManager.AddLibrary(newLibrary.GetId(), newLibrary.GetStruct())
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	// Récupérer la sauvegarde du monde
	// saveFichier -> &World{}
}

func main() {

	serverEngine := engine.ServerEngine{}

	var sync_mutex sync.Mutex
	var sync_group sync.WaitGroup
	// udpMessages := make(map[string](chan []byte), 100)

	var clients map[string]*Client = make(map[string]*Client)

	newUUID, err := uuid.Parse("b5a2aa76-984b-4e74-a5ae-99fef469e153")
	if err != nil {
		log.Println("COuld not parse string id to uuid")
	}
	clients[newUUID.String()] = &Client{
		Email:      "baw.developpement@gmail.com",
		Password:   "azerty",
		ID:         newUUID,
		Token:      newUUID,
		Characters: []ecs.Entity{},
	}

	loicId := uuid.New()
	clients[loicId.String()] = &Client{
		Email:    "loic@gmail.com",
		Password: "azerty",
		ID:       loicId,
		Token:    loicId,
	}

	var gameWorld ecs.IWorld = &ecs.World{
		Id: "test",
	}

	var gameEntities []*ecs.IEntity = []*ecs.IEntity{
		ecs.CEntityPossessed(&gameWorld, uuid.New(), newUUID, []*ecs.Component{
			LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"}, map[string]interface{}{"x": 0.0, "y": 0.0, "origin_x": 0.0, "origin_y": 0.0, "vel_x": 0.0, "vel_y": 0.0}),
			LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "dimension"}, map[string]interface{}{"width": 0, "height": 0}),
			LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "solid"}, true),
		}),
	}

	for _, entityLocation := range gameEntities {
		gameWorld.AddEntity(entityLocation)
	}

	var netUpdateSystem ecs.ISystem = &systems.NetUpdateSystem{
		Id: ecs.Identifier{
			Namespace: "oneforx",
			Path:      "update",
		},
		Name:         "net_update",
		ServerEngine: &serverEngine,
		World:        &gameWorld,
	}

	system, err := LibraryManager.InstantiateSystem(ecs.Identifier{Namespace: "oneforx", Path: "bullet"}, &gameWorld)
	if err != nil {
		panic(err)
	}
	gameWorld.AddSystem(system)
	gameWorld.AddSystem(&netUpdateSystem)

	// Run world.Update()
	go func() {
		ticker := time.NewTicker(time.Second / 120) // Crée un ticker qui envoie un signal toutes les 1/60ème de seconde
		defer ticker.Stop()                         // Arrête le ticker quand la fonction main se termine
		for range ticker.C {
			sync_mutex.Lock()
			gameWorld.Update()
			sync_mutex.Unlock()
		}
	}()

	// Ici nous avons une fonction Start qui va démarrer un serveur TCP et UDP
	// La fonction accepte 4 arguments, l'addresse, un canal de message pour udp
	// et deux handlers, chacune pour chaque type de serveur
	serverEngine.Start("localhost:8080")
	defer serverEngine.Close()

	sync_group.Add(3)
	go func() {
		defer sync_group.Done()
		serverEngine.ListenTCP(func(tcpConnexion *net.TCPConn) {
			startTime := time.Now()
			// bytes, err := engine.MessageToBytes(messages.SC_PING(0))
			// if err != nil {
			// 	log.Println("Could not convert message to bytes")
			// }

			// if err == nil {
			// 	_, err = tcpConnexion.Write(bytes)
			// 	if err != nil {
			// 		log.Printf("Could not send ping to %s", tcpConnexion.RemoteAddr().String())
			// 	}
			// }

			log.Println("TCP ADDRESS: ", tcpConnexion.RemoteAddr().String())

			for {
				buf := make([]byte, 1024)
				n, err := tcpConnexion.Read(buf)
				if err != nil {
					continue
				}

				log.Println("NEW_MESSAGE", string(buf[:n]))

				var message engine.Message

				if err := json.Unmarshal(buf[:n], &message); err != nil {
					log.Println("Could not read message for client id")
					continue
				}

				if message.MessageType == "DISCONNECT" {
					// Send back DISCONNECT -> to disconnect the user
					bytes, err := json.Marshal(message)
					if err != nil {
						return
					}

					_, err = tcpConnexion.Write(bytes)
					if err != nil {
						return
					}

					tcpConnexion.Close()
				} else if message.MessageType == "CONNECT" {
					messageData, ok := message.Data.(map[string]interface{})
					if !ok {
						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("data n'est pas un objet"))
						if err != nil {
							log.Println(err.Error())
						} else {
							tcpConnexion.Write(bytes)
						}
						continue
					}
					email, ok := messageData["email"].(string)
					if !ok {
						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
						if err != nil {
							log.Println(err.Error())
						} else {
							tcpConnexion.Write(bytes)
						}
						continue
					}

					password, ok := messageData["password"].(string)
					if !ok {
						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
						if err != nil {
							log.Println(err.Error())
						} else {
							tcpConnexion.Write(bytes)
						}
						continue
					}
					var foundId *string

					for id, c := range clients {
						if c.Email == email && c.Password == password {
							foundId = &id
							break
						}
					}

					if foundId != nil {
						foundClient := clients[*foundId]
						foundClient.TCP = tcpConnexion

						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS(foundClient.Token.String()))
						if err != nil {
							log.Println()
						}

						tcpConnexion.Write(bytes)
					} else {
						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("The compte " + email + " n'as pas été trouvé"))
						if err != nil {
							log.Println()
						}

						tcpConnexion.Write(bytes)
					}
				} else if message.MessageType == "CONNECT_TOKEN" {
					// We bind the connexion to the Client
					tokenString, ok := message.Data.(string)
					if !ok {
						log.Println("CONNECT_TOKEN Could not parse data to string")
					}

					token, err := uuid.Parse(tokenString)
					if err != nil {
						log.Println("CONNECT_TOKEN Could not parse tokenString to uuid")
					}

					var foundId *string

					for id, client := range clients {
						if client.Token.String() == token.String() {
							foundId = &id
							break
						}
					}

					if foundId != nil {
						foundClient := clients[*foundId]
						foundClient.TCP = tcpConnexion

						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS("Nous avons bien associé votre compte à votre connexion"))
						if err != nil {
							log.Println("Could not parse SC_CONNECT_SUCCESS to bytes")
						}
						tcpConnexion.Write(bytes)
					} else {
						bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("Le token que vous aviez était obsoléte"))
						if err != nil {
							log.Println("Could not parse SC_CONNECT_FAILED to bytes")
						}
						tcpConnexion.Write(bytes)
					}
				} else if message.MessageType == "PONG" {
					endTime := time.Now()
					latency := endTime.Sub(startTime)
					startTime = time.Now()
					bytes, err := engine.MessageToBytes(messages.SC_PING(int(latency.Milliseconds())))
					if err != nil {
						log.Println("Could not parse message to bytes")
						// Evitons de lancer le message en passant au prochain message
						continue
					}

					_, err = tcpConnexion.Write(bytes)
					if err != nil {
						log.Println("Could not send message ping")
					}
				} else if message.MessageType == "SHOOT" {
					var foundId *uuid.UUID

					for _, c := range clients {
						if c.TCP != nil {
							if c.TCP.RemoteAddr().String() == tcpConnexion.RemoteAddr().String() {
								foundId = &c.ID
								break
							}
						}
					}

					if foundId != nil {

						newBullet := ecs.CEntityWithOwner(&gameWorld, uuid.New(), *foundId, []*ecs.Component{
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"}, map[string]interface{}{"x": 0.0, "y": 0.0, "origin_x": 0.0, "origin_y": 0.0, "vel_x": 0.0, "vel_y": 0.0}),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "dimension"}, map[string]interface{}{"width": 0, "height": 0}),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "solid"}, true),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "force"}, 0),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "orientation"}, 0),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "speed"}, 10),
							LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "distance"}, 1000),
						})

						gameWorld.AddEntity(newBullet)

						bytes, err := engine.MessageToBytes(messages.SC_CREATE_ENTITY(ecs.EntityToNoCycle(newBullet)))
						if err != nil {
							log.Println("Could not parse SC_CREATE_ENTITY to bytes")
							continue
						}
						_, err = tcpConnexion.Write(bytes)
						if err != nil {
							log.Println("Could not send bytes of SC_CREATE_ENTITY to client: " + tcpConnexion.RemoteAddr().String())
						}
					}
				} else if message.MessageType == "CONNECT_WORLD" {
					var sync_mutex sync.Mutex
					var foundId *uuid.UUID

					sync_mutex.Lock()
					for _, c := range clients {
						if c.TCP != nil {
							if c.TCP.RemoteAddr().String() == tcpConnexion.RemoteAddr().String() {
								foundId = &c.ID
								break
							}
						}
					}
					sync_mutex.Unlock()

					if foundId != nil {
						entitiesPossessed := gameWorld.GetEntitiesPossessedBy(*foundId)
						log.Println(len(entitiesPossessed))
						if len(entitiesPossessed) == 0 {
							gameWorld.AddEntity(ecs.CEntityPossessed(&gameWorld, uuid.New(), *foundId, []*ecs.Component{
								LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"}, map[string]interface{}{"x": 0.0, "y": 0.0, "origin_x": 0.0, "origin_y": 0.0, "vel_x": 0.0, "vel_y": 0.0}),
								LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "dimension"}, map[string]interface{}{"width": 0, "height": 0}),
								LibraryManager.InstantiateComponent(ecs.Identifier{Namespace: "oneforx", Path: "solid"}, true),
							}))
						}

						bytes, err := engine.MessageToBytes(messages.SC_CREATE_WORLD(gameWorld.GetEntitiesNoCycle()))
						if err != nil {
							log.Println("Could not parse SC_CREATE_WORLD to bytes")
							continue
						}
						_, err = tcpConnexion.Write(bytes)
						if err != nil {
							log.Println("Could not send bytes of SC_CREATE_WORLD to client: " + tcpConnexion.RemoteAddr().String())
						}
					}
				} else if message.MessageType == "DISCONNECT_WORLD" {

				}
			}
		})

	}()

	go func() {
		defer sync_group.Done()

		serverEngine.ListenMessages(func(message engine.Message, addr *net.UDPAddr) {
			switch message.MessageType {
			case "CONNECT":
				messageData, ok := message.Data.(map[string]interface{})
				if !ok {
					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("data n'est pas un objet"))
					if err != nil {
						log.Println(err.Error())
					} else {
						serverEngine.UdpListener.WriteToUDP(bytes, addr)
					}
					return
				}
				email, ok := messageData["email"].(string)
				if !ok {
					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
					if err != nil {
						log.Println(err.Error())
					} else {
						serverEngine.UdpListener.WriteToUDP(bytes, addr)
					}
					return
				}

				password, ok := messageData["password"].(string)
				if !ok {
					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
					if err != nil {
						log.Println(err.Error())
					} else {
						serverEngine.UdpListener.WriteToUDP(bytes, addr)
					}
					return
				}
				var foundId *string

				for id, c := range clients {
					if c.Email == email && c.Password == password {
						foundId = &id
						break
					}
				}

				if foundId != nil {
					foundClient := clients[*foundId]
					foundClient.UDP = addr
					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS(foundClient.Token.String()))
					if err != nil {
						log.Println()
					}

					serverEngine.UdpListener.WriteToUDP(bytes, addr)
				} else {
					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("The compte " + email + " n'as pas été trouvé"))
					if err != nil {
						log.Println()
					}

					serverEngine.UdpListener.WriteToUDP(bytes, addr)
				}
			case "CONNECT_TOKEN":
				// We bind the connexion to the Client
				tokenString, ok := message.Data.(string)
				if !ok {
					log.Println("CONNECT_TOKEN Could not parse data to string")
				}

				token, err := uuid.Parse(tokenString)
				if err != nil {
					log.Println("CONNECT_TOKEN Could not parse tokenString to uuid")
				}

				var foundClient *Client

				for _, client := range clients {
					if client.Token == token {
						foundClient = client
						break
					}
				}

				if foundClient != nil {
					foundClient.UDP = addr

					bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS("Nous avons bien associé votre compte à votre connexion"))
					if err != nil {
						log.Println("Could not send SC_CONNECT_SUCCESS")
					}
					serverEngine.UdpListener.WriteToUDP(bytes, addr)
				}
			case "MOVE":
				go func() {
					direction := message.Data.(string)

					sync_mutex.Lock()
					var foundId *uuid.UUID
					for _, c := range clients {
						if c.UDP != nil {
							if c.UDP.String() == addr.String() {
								foundId = &c.ID
							}
						}
					}
					sync_mutex.Unlock()

					if foundId != nil {
						if entities := gameWorld.GetEntitiesPossessedBy(*foundId); len(entities) != 0 {
							for _, entity := range entities {
								sync_mutex.Lock()

								entityLocalised := *entity
								positionComponent := entityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"})
								if positionComponent == nil {
									log.Println("Could not get component position")
									continue
								}
								positionComponentLocalised := *positionComponent

								positionComponentDataParsed, ok := positionComponentLocalised.GetData().(map[string]interface{})
								if !ok {
									log.Println("Could not parse position component to map[string]float64")
								}

								x := positionComponentDataParsed["x"].(float64)
								y := positionComponentDataParsed["y"].(float64)

								switch direction {
								case "UP":
									positionComponentDataParsed["y"] = y - 2
								case "DOWN":
									positionComponentDataParsed["y"] = y + 2
								case "LEFT":
									positionComponentDataParsed["x"] = x - 2
								case "RIGHT":
									positionComponentDataParsed["x"] = x + 2
								}

								log.Println(positionComponentLocalised.GetData())
								sync_mutex.Unlock()
							}
						}
					}
				}()
			case "ORIENTATION":
				mousePosition, ok := message.Data.(map[string]interface{})
				if !ok {
					log.Println("Could not parse CS_ORIENTATION data to map[string]int")
				}

				log.Println(mousePosition)
			}
		})
	}()

	sync_group.Wait()

	go func() {
		defer sync_group.Done()
		serverEngine.ListenUDP(func(masterConnexion *net.UDPConn, addr *net.UDPAddr, messagesChannel chan []byte) {
			log.Println("UDP ADDRESS: ", addr.String())

			for {
				msg := <-messagesChannel
				message, err := engine.BytesToMessage(msg)
				if err != nil {
					log.Println(err.Error())
				}

				if err == nil {
					switch message.MessageType {
					case "CONNECT":
						messageData, ok := message.Data.(map[string]interface{})
						if !ok {
							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("data n'est pas un objet"))
							if err != nil {
								log.Println(err.Error())
							} else {
								masterConnexion.WriteToUDP(bytes, addr)
							}
							continue
						}
						email, ok := messageData["email"].(string)
						if !ok {
							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
							if err != nil {
								log.Println(err.Error())
							} else {
								masterConnexion.WriteToUDP(bytes, addr)
							}
							continue
						}

						password, ok := messageData["password"].(string)
						if !ok {
							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("l'objet manque l'email ou n'est pas une chaine de caractères"))
							if err != nil {
								log.Println(err.Error())
							} else {
								masterConnexion.WriteToUDP(bytes, addr)
							}
							continue
						}
						var foundId *string

						for id, c := range clients {
							if c.Email == email && c.Password == password {
								foundId = &id
								break
							}
						}

						if foundId != nil {
							foundClient := clients[*foundId]
							foundClient.UDP = addr
							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS(foundClient.Token.String()))
							if err != nil {
								log.Println()
							}

							masterConnexion.WriteToUDP(bytes, addr)
						} else {
							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_FAILED("The compte " + email + " n'as pas été trouvé"))
							if err != nil {
								log.Println()
							}

							masterConnexion.WriteToUDP(bytes, addr)
						}
					case "CONNECT_TOKEN":
						// We bind the connexion to the Client
						tokenString, ok := message.Data.(string)
						if !ok {
							log.Println("CONNECT_TOKEN Could not parse data to string")
						}

						token, err := uuid.Parse(tokenString)
						if err != nil {
							log.Println("CONNECT_TOKEN Could not parse tokenString to uuid")
						}

						var foundClient *Client

						for _, client := range clients {
							if client.Token == token {
								foundClient = client
								break
							}
						}

						if foundClient != nil {
							foundClient.UDP = addr

							bytes, err := engine.MessageToBytes(messages.SC_CONNECT_SUCCESS("Nous avons bien associé votre compte à votre connexion"))
							if err != nil {
								log.Println("Could not send SC_CONNECT_SUCCESS")
							}
							masterConnexion.WriteToUDP(bytes, addr)
						}
					case "MOVE":
						go func() {
							direction := message.Data.(string)

							sync_mutex.Lock()
							var foundId *uuid.UUID
							for _, c := range clients {
								if c.UDP != nil {
									if c.UDP.String() == addr.String() {
										foundId = &c.ID
									}
								}
							}
							sync_mutex.Unlock()

							if foundId != nil {
								if entities := gameWorld.GetEntitiesPossessedBy(*foundId); len(entities) != 0 {
									for _, entity := range entities {
										sync_mutex.Lock()

										entityLocalised := *entity
										positionComponent := entityLocalised.GetComponent(ecs.Identifier{Namespace: "oneforx", Path: "position"})
										if positionComponent == nil {
											log.Println("Could not get component position")
											continue
										}
										positionComponentLocalised := *positionComponent

										positionComponentDataParsed, ok := positionComponentLocalised.GetData().(map[string]interface{})
										if !ok {
											log.Println("Could not parse position component to map[string]float64")
										}

										x := positionComponentDataParsed["x"].(float64)
										y := positionComponentDataParsed["y"].(float64)

										switch direction {
										case "UP":
											positionComponentDataParsed["y"] = y - 2
										case "DOWN":
											positionComponentDataParsed["y"] = y + 2
										case "LEFT":
											positionComponentDataParsed["x"] = x - 2
										case "RIGHT":
											positionComponentDataParsed["x"] = x + 2
										}

										log.Println(positionComponentLocalised.GetData())
										sync_mutex.Unlock()
									}
								}
							}
						}()
					case "ORIENTATION":
						// mousePosition, ok := handleMessage.Data.(map[string]int)
						// if !ok {
						// 	log.Println("Could not parse CS_ORIENTATION data to map[string]int")
						// }
					}
				}

				<-messagesChannel
			}

		})

	}()
}

// func Old() {
// 	for {
// 		sync_mutex.Lock()
// 		message := <-serverEngine.UdpMessages[addr.String()]
// 		sync_mutex.Unlock()

// 		var handleMessage engine.Message
// 		if err := json.Unmarshal(message, &handleMessage); err != nil {
// 			log.Println("Failed to read message", message)
// 		}

// 		// fmt.Printf("Received parsed message %v\n", handleMessage)

// 		<-serverEngine.UdpMessages[addr.String()]
// 	}
// }
