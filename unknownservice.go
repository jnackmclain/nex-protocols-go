package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	UnknownProtocolID = 0x17 // i have absolutely no clue what this does
	UnknownMethod     = 0x3  // same as above
)

type UnknownProtocol struct {
	server               *nex.Server
	ConnectionIDCounter  *nex.Counter
	UnknownMethodHandler func(err error, client *nex.Client, callID uint32, pid uint32)
}

func (unknownProtocol *UnknownProtocol) Setup() {
	nexServer := unknownProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if UnknownProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case UnknownMethod:
				go unknownProtocol.handleUnknownMethod(packet)
			default:
				log.Printf("Unsupported Unknown method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (unknownProtocol *UnknownProtocol) UnknownMethod(handler func(err error, client *nex.Client, callID uint32, pid uint32)) {
	unknownProtocol.UnknownMethodHandler = handler
}

func (unknownProtocol *UnknownProtocol) handleUnknownMethod(packet nex.PacketInterface) {
	if unknownProtocol.UnknownMethodHandler == nil {
		log.Println("[Warning] UnknownProtocol::UnknownMethodHandler not implemented")
		go respondNotImplemented(packet, UnknownProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, unknownProtocol.server)

	pid := parametersStream.ReadUInt32LE()
	parametersStream.ReadUInt32LE()
	parametersStream.ReadUInt32LE()
	parametersStream.ReadUInt32LE()

	go unknownProtocol.UnknownMethodHandler(nil, client, callID, pid)
}

// NewSecureProtocol returns a new SecureProtocol
func NewUnknownProtocol(server *nex.Server) *UnknownProtocol {
	unknownProtocol := &UnknownProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	unknownProtocol.Setup()

	return unknownProtocol
}
