package nexproto

import nex "github.com/ihatecompvir/nex-go"

/*
	NEX and Rendez-Vous have multiple protocols for match making
	These protocols all share the same types
	In an effort to keep this library organized, each type used in all match making protocols is defined here
*/

// Gathering holds information about a matchmake gathering
type Gathering struct {
	ID                  uint32
	OwnerPID            uint32
	HostPID             uint32
	MinimumParticipants uint16
	MaximumParticipants uint16
	ParticipationPolicy uint32
	PolicyArgument      uint32
	Flags               uint32
	State               uint32
	Description         string

	*nex.Structure
}

// ExtractFromStream extracts a Gathering structure from a stream
func (gathering *Gathering) ExtractFromStream(stream *StreamIn) error {
	var err error

	gathering.ID = stream.ReadUInt32LE()
	gathering.OwnerPID = stream.ReadUInt32LE()
	gathering.HostPID = stream.ReadUInt32LE()
	gathering.MinimumParticipants = stream.ReadUInt16LE()
	gathering.MaximumParticipants = stream.ReadUInt16LE()
	gathering.ParticipationPolicy = stream.ReadUInt32LE()
	gathering.PolicyArgument = stream.ReadUInt32LE()
	gathering.Flags = stream.ReadUInt32LE()
	gathering.State = stream.ReadUInt32LE()
	gathering.Description, err = stream.ReadString()

	if err != nil {
		return err
	}

	return nil
}

// NewGathering returns a new Gathering
func NewGathering() *Gathering {
	return &Gathering{}
}
