package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	MatchmakingProtocolID  = 0x15 // the first matchmaking service protocol
	MatchmakingProtocolID2 = 0x6E // the second matchmaking service? if there is a joinable gathering, this returns it

	RegisterGathering  = 0x1  // creates an empty gathering, also CheckForGatherings
	TerminateGathering = 0x2  // ends a gathering
	UpdateGathering    = 0x4  // updates a gathering
	Participate        = 0xB  // unsure on this one, going off NintendoClients wiki for the name
	Unparticipate      = 0xC  // unsure on this one, going off NintendoClients wiki for the name
	LaunchSession      = 0x1A // unsure on this one, going off NintendoClients wiki for the name
)

// JsonProtocol handles the Json requests
type MatchmakingProtocol struct {
	server                    *nex.Server
	ConnectionIDCounter       *nex.Counter
	RegisterGatheringHandler  func(err error, client *nex.Client, callID uint32, gathering []byte)
	UpdateGatheringHandler    func(err error, client *nex.Client, callID uint32, gathering []byte, gatheringID uint32)
	ParticipateHandler        func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	UnparticipateHandler      func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	LaunchSessionHandler      func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	TerminateGatheringHandler func(err error, client *nex.Client, callID uint32, gatheringID uint32)

	// second handler
	CheckForGatheringsHandler func(err error, client *nex.Client, callID uint32, data []byte)
}

func (matchmakingProtocol *MatchmakingProtocol) Setup() {
	nexServer := matchmakingProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if MatchmakingProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go matchmakingProtocol.handleRegisterGathering(packet)
			case UpdateGathering:
				go matchmakingProtocol.handleUpdateGathering(packet)
			case Participate:
				go matchmakingProtocol.handleParticipate(packet)
			case Unparticipate:
				go matchmakingProtocol.handleUnparticipate(packet)
			case LaunchSession:
				go matchmakingProtocol.handleLaunchSession(packet)
			case TerminateGathering:
				go matchmakingProtocol.handleTerminateGathering(packet)
			default:
				log.Printf("Unsupported Matchmaking method ID: %#v\n", request.MethodID())
			}
		}

		if MatchmakingProtocolID2 == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go matchmakingProtocol.handleCheckForGatherings(packet)
			default:
				log.Printf("Unsupported Matchmaking2 method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (matchmakingProtocol *MatchmakingProtocol) RegisterGathering(handler func(err error, client *nex.Client, callID uint32, gathering []byte)) {
	matchmakingProtocol.RegisterGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) UpdateGathering(handler func(err error, client *nex.Client, callID uint32, gathering []byte, gatheringID uint32)) {
	matchmakingProtocol.UpdateGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) Participate(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.ParticipateHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) Unparticipate(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.UnparticipateHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) LaunchSession(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.LaunchSessionHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) TerminateGathering(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.TerminateGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) CheckForGatherings(handler func(err error, client *nex.Client, callID uint32, data []byte)) {
	matchmakingProtocol.CheckForGatheringsHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) handleRegisterGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::RegisterGatheringHandler not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	parametersStream.Read4ByteString()
	parametersStream.ReadUInt32LE()
	gathering, err := parametersStream.ReadBuffer()

	if err != nil {
		log.Println("Could not read gathering data")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	go matchmakingProtocol.RegisterGatheringHandler(nil, client, callID, gathering)
}

func (matchmakingProtocol *MatchmakingProtocol) handleUpdateGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::RegisterGatheringHandler not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	parametersStream.Read4ByteString()
	parametersStream.ReadUInt32LE()
	gathering, err := parametersStream.ReadBuffer()

	gatheringStream := NewStreamIn(gathering, matchmakingProtocol.server)

	gatheringID := gatheringStream.ReadUInt32LE()

	if err != nil {
		log.Println("Could not read gathering data")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	go matchmakingProtocol.UpdateGatheringHandler(nil, client, callID, gathering, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleParticipate(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::Participate not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.ParticipateHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleUnparticipate(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::Unparticipate not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.UnparticipateHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleLaunchSession(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::LaunchSession not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.LaunchSessionHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleTerminateGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::LaunchSession not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.TerminateGatheringHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleCheckForGatherings(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::LaunchSession not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	//parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	//gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.CheckForGatheringsHandler(nil, client, callID, parameters)
}

// NewSecureProtocol returns a new SecureProtocol
func NewMatchmakingProtocol(server *nex.Server) *MatchmakingProtocol {
	matchmakingProtocol := &MatchmakingProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	matchmakingProtocol.Setup()

	return matchmakingProtocol
}
