package nexproto

import nex "github.com/ihatecompvir/nex-go"

func respondNotImplemented(packet nex.PacketInterface, protocolID uint8) {
	client := packet.Sender()
	request := packet.RMCRequest()

	rmcResponse := nex.NewRMCResponse(protocolID, request.CallID())
	rmcResponse.SetError(0x80010002)

	rmcResponseBytes := rmcResponse.Bytes()

	var responsePacket nex.PacketInterface

	responsePacket, _ = nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(packet.Version())
	responsePacket.SetSource(packet.Destination())
	responsePacket.SetDestination(packet.Source())
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	client.Server().Send(responsePacket)
}
