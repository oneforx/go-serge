package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/oneforx/go-ecs"
)

type ServerEngine struct {
	TcpListener   *net.TCPListener
	TcpInit       bool
	TcpConnexions map[string]*net.TCPConn

	UdpListener   *net.UDPConn
	UdpInit       bool
	UdpConnexions map[string]*net.UDPAddr
	UdpMessages   map[string]chan []byte

	ServerMutex sync.RWMutex
}

type Client struct {
	Id uuid.UUID
}

func (se *ServerEngine) Start(address string) {
	var wait_mutex sync.Mutex
	// Listen for TCP connections
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	wait_mutex.Lock()
	se.TcpListener = tcpListener
	se.TcpInit = true
	se.TcpConnexions = map[string]*net.TCPConn{}
	wait_mutex.Unlock()

	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	wait_mutex.Lock()
	se.UdpListener = udpListener
	se.UdpInit = true
	se.UdpConnexions = make(map[string]*net.UDPAddr)
	se.UdpMessages = make(map[string]chan []byte, 1000)
	wait_mutex.Unlock()
}

func (se *ServerEngine) Close() {
	se.TcpListener.Close()
	se.UdpListener.Close()
}

func (se *ServerEngine) ListenTCP(connexionHandler func(connexion *net.TCPConn)) {
	for {
		if !se.TcpInit {
			continue
		}
		tcpConnexion, err := se.TcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		se.ServerMutex.Lock()
		if se.TcpConnexions[tcpConnexion.RemoteAddr().String()] == nil {
			se.TcpConnexions[tcpConnexion.RemoteAddr().String()] = tcpConnexion
		}
		se.ServerMutex.Unlock()

		se.ServerMutex.Lock()
		// Vérifiez si la connexion est fermée avant de l'utiliser
		if tcpConnexion == nil || tcpConnexion.RemoteAddr() == nil {
			continue
		}
		se.ServerMutex.Unlock()

		log.Printf("Got a client %v", tcpConnexion.RemoteAddr().String())

		go connexionHandler(tcpConnexion)
	}
}

func (se *ServerEngine) ListenMessages(messageHandler func(message Message, addr *net.UDPAddr)) {
	for {
		if !se.UdpInit {
			continue
		}
		buf := make([]byte, 1024)
		n, addr, err := se.UdpListener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}

		se.ServerMutex.Lock()
		_, ok := se.UdpConnexions[addr.String()]
		se.ServerMutex.Unlock()
		if !ok {
			se.ServerMutex.Lock()
			se.UdpConnexions[addr.String()] = addr
			se.ServerMutex.Unlock()
		}

		message, err := BytesToMessage(buf[:n])
		if err != nil {
			continue
		}

		messageHandler(message, addr)
	}
}

func (se *ServerEngine) ListenUDP(connexionHandler func(*net.UDPConn, *net.UDPAddr, chan []byte)) {
	for {
		if !se.UdpInit {
			continue
		}
		buf := make([]byte, 1024)
		n, addr, err := se.UdpListener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}

		se.ServerMutex.Lock()
		_, ok := se.UdpConnexions[addr.String()]
		se.ServerMutex.Unlock()
		if !ok {
			se.ServerMutex.Lock()
			se.UdpMessages[addr.String()] = make(chan []byte, 1000)
			se.ServerMutex.Unlock()
			go connexionHandler(se.UdpListener, addr, se.UdpMessages[addr.String()])
		}

		se.ServerMutex.Lock()
		if se.UdpConnexions[addr.String()] == nil {
			se.UdpConnexions[addr.String()] = addr
		}
		se.ServerMutex.Unlock()

		// Ici envois le message dans le channel
		// C'est mieux qu'un select
		go func() {
			se.ServerMutex.Lock()
			se.UdpMessages[addr.String()] <- buf[:n]
			se.ServerMutex.Unlock()
		}()

	}
}

func (se *ServerEngine) USendToAll(message Message) *ecs.FeedBack {
	bytes, err := MessageToBytes(message)
	if err != nil {
		return &ecs.FeedBack{
			Host:  "USendToAll",
			Job:   "MessageToBytes",
			Label: "COULD_NOT_CONVERT_MESSAGE",
			Data:  bytes,
		}
	}

	var writeToFeedback *ecs.FeedBack

	for _, addr := range se.UdpConnexions {
		_, err := se.UdpListener.WriteTo(bytes, addr)
		if err != nil {
			if writeToFeedback == nil {
				writeToFeedback = &ecs.FeedBack{
					Host:  "USendToAll",
					Job:   "WriteTo",
					Label: "COULD_NOT_WRITE_TO",
					Data:  []string{},
				}
			}

			// Vue que nous cherchons plusieurs connexions, il y aura plusieurs addresses
			// Avec le retour de la fonction USendToAll vous pourriez faire log.Print(feedBack.Label + ": ", Join(feedBack.Data, ", ")
			data := writeToFeedback.Data.([]string)
			writeToFeedback.Data = append(data, addr.String())
		}

	}

	if writeToFeedback != nil {
		return writeToFeedback
	}

	return &ecs.FeedBack{
		Host:  "USendToAll",
		Job:   "return",
		Label: "SUCCESS",
	}
}

func (se *ServerEngine) USendToAddr(address string, message Message) *ecs.FeedBack {
	bytes, err := MessageToBytes(message)
	if err != nil {
		return &ecs.FeedBack{
			Host:  "USendToAddr",
			Job:   "MessageToBytes",
			Label: "COULD_NOT_CONVERT_MESSAGE",
			Data:  message,
		}
	}

	var writeToFeedback *ecs.FeedBack

	for addrString, addr := range se.UdpConnexions {
		if addrString == address {
			_, err := se.UdpListener.WriteToUDP(bytes, addr)
			if err != nil {
				if writeToFeedback == nil {
					writeToFeedback = &ecs.FeedBack{
						Host:  "USendToAddr",
						Job:   "WriteToUDP",
						Label: "COULD_NOT_WRITE_TO",
					}
				}
				writeToFeedback.Data = addr
			}

			// Dès qu'on trouve la  connexion cible, on arrête de chercher
			break
		}
	}

	if writeToFeedback != nil {
		return writeToFeedback
	}

	return &ecs.FeedBack{
		Host:  "USendToAddr",
		Job:   "return",
		Label: "SUCCESS",
	}
}

func (se *ServerEngine) TSendToAll(message Message) *ecs.FeedBack {
	bytes, err := MessageToBytes(message)
	if err != nil {
		return &ecs.FeedBack{
			Host:  "TSendToAll",
			Job:   "MessageToBytes",
			Label: "COULD_NOT_CONVERT_MESSAGE",
			Data:  bytes,
		}
	}

	var writeToFeedback *ecs.FeedBack

	for _, connexion := range se.TcpConnexions {
		_, err := connexion.Write(bytes)
		if err != nil {
			if writeToFeedback == nil {
				writeToFeedback = &ecs.FeedBack{
					Host:  "TSendToAll",
					Job:   "MessageToBytes",
					Label: "COULD_NOT_WRITE_TO",
					Data:  []string{},
				}
			}

			// Vue que nous cherchons plusieurs connexions, il y aura plusieurs addresses
			// Avec le retour de la fonction TSendToAll vous pourriez faire log.Print(Label + ": ", ConcatData)
			data := writeToFeedback.Data.([]string)
			writeToFeedback.Data = append(data, connexion.RemoteAddr().String())
		}
	}

	if writeToFeedback != nil {
		return writeToFeedback
	}

	return &ecs.FeedBack{
		Host:  "TsendToAll",
		Job:   "return",
		Label: "SUCCESS",
	}
}

func (se *ServerEngine) TSendToAddr(address string, message Message) *ecs.FeedBack {
	bytes, err := MessageToBytes(message)
	if err != nil {
		return &ecs.FeedBack{
			Host:  "TSendToAddr",
			Job:   "MessageToBytes",
			Label: "COULD_NOT_CONVERT_MESSAGE",
			Data:  message,
		}
	}

	var writeToFeedback *ecs.FeedBack

	for addr, connexion := range se.TcpConnexions {
		if addr == address {
			_, err := connexion.Write(bytes)
			if err != nil {
				if writeToFeedback == nil {
					writeToFeedback = &ecs.FeedBack{
						Host:  "TSendToAddr",
						Job:   "connexion.Write()",
						Label: "COULD_NOT_WRITE_TO",
					}
				}
				writeToFeedback.Data = addr
			}

			// Dès qu'on trouve la  connexion cible, on arrête de chercher
			break
		}
	}

	if writeToFeedback != nil {
		return writeToFeedback
	}

	return &ecs.FeedBack{
		Host:  "TSendToAddr",
		Job:   "return",
		Label: "SUCCESS",
	}
}

func MessageToBytes(msg Message) ([]byte, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func BytesToMessage(bytes []byte) (Message, error) {
	var message Message
	if err := json.Unmarshal(bytes, &message); err != nil {
		return Message{}, err
	}
	return message, nil
}

func LogFeedBack(fb ecs.FeedBack) {
	type FeedbackData struct {
		Data interface{} `json:"data"`
	}

	feedbackData := FeedbackData{
		Data: fb.Data,
	}

	_, err := json.Marshal(feedbackData)
	if err != nil {
		return
	}

	log.Println("["+fb.Host+"]"+"["+fb.Job+"]"+"["+fb.Label+"]: ", fb.Data)
}
