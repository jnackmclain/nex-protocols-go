package nexproto

import (
	"errors"
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// AccountManagementProtocolID is the protocol ID for the Account Management protocol
	AccountManagementProtocolID = 0x19

	// AccountManagementMethodNintendoCreateAccount is the method ID for the method NintendoCreateAccount
	AccountManagementMethodNintendoCreateAccount = 0x1B
)

// AccountManagementProtocol handles the Account Management nex protocol
type AccountManagementProtocol struct {
	server                       *nex.Server
	NintendoCreateAccountHandler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)
}

// Setup initializes the protocol
func (accountManagementProtocol *AccountManagementProtocol) Setup() {
	nexServer := accountManagementProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if AccountManagementProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case AccountManagementMethodNintendoCreateAccount:
				go accountManagementProtocol.handleNintendoCreateAccountHandler(packet)
			default:
				log.Printf("Unsupported AccountManagement method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// NintendoCreateAccount sets the NintendoCreateAccount handler function
func (accountManagementProtocol *AccountManagementProtocol) NintendoCreateAccount(handler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)) {
	accountManagementProtocol.NintendoCreateAccountHandler = handler
}

func (accountManagementProtocol *AccountManagementProtocol) handleNintendoCreateAccountHandler(packet nex.PacketInterface) {
	if accountManagementProtocol.NintendoCreateAccountHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::NintendoCreateAccount not implemented")
		go respondNotImplemented(packet, AccountManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, accountManagementProtocol.server)

	username, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	key, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	groups := parametersStream.ReadUInt32LE()
	email, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	dataHolderName, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	if dataHolderName != "NintendoToken" {
		err := errors.New("[AccountManagementProtocol::NintendoCreateAccount] Data holder name does not match")
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	go accountManagementProtocol.NintendoCreateAccountHandler(nil, client, callID, username, key, groups, email)
}

// NewAccountManagementProtocol returns a new AccountManagementProtocol
func NewAccountManagementProtocol(server *nex.Server) *AccountManagementProtocol {
	accountManagementProtocol := &AccountManagementProtocol{server: server}

	accountManagementProtocol.Setup()

	return accountManagementProtocol
}
