package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	NATTraversalID = 0x3 // the first matchmaking service protocol

	InitiateProbe = 0x1
)

// JsonProtocol handles the Json requests
type NATTraversalProtocol struct {
	server               *nex.Server
	ConnectionIDCounter  *nex.Counter
	InitiateProbeHandler func(err error, client *nex.Client, callID uint32, stationURL string)
}

func (natTraversalProtocol *NATTraversalProtocol) Setup() {
	nexServer := natTraversalProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if NATTraversalID == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go natTraversalProtocol.handleInitiateProbe(packet)
			default:
				log.Printf("Unsupported NAT traversal method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (natTraversalProtocol *NATTraversalProtocol) InitiateProbe(handler func(err error, client *nex.Client, callID uint32, stationURL string)) {
	natTraversalProtocol.InitiateProbeHandler = handler
}

func (natTraversalProtocol *NATTraversalProtocol) handleInitiateProbe(packet nex.PacketInterface) {
	if natTraversalProtocol.InitiateProbeHandler == nil {
		log.Println("[Warning] NATTraversal::InitiateProbe not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, natTraversalProtocol.server)

	parametersStream.ReadUInt32LE()
	stationURL, err := parametersStream.Read4ByteString()

	if err != nil {
		go natTraversalProtocol.InitiateProbeHandler(nil, client, callID, "")
	}

	go natTraversalProtocol.InitiateProbeHandler(nil, client, callID, stationURL)
}

// NewSecureProtocol returns a new SecureProtocol
func NewNATTraversalProtocol(server *nex.Server) *NATTraversalProtocol {
	natTraversalProtocol := &NATTraversalProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	natTraversalProtocol.Setup()

	return natTraversalProtocol
}
