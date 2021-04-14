package nexproto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// SecureProtocolID is the protocol ID for the Secure Connection protocol
	JsonProtocolID = 0x75

	// SecureMethodRegister is the method ID for the method Register
	JsonRequest = 0x1
)

// JsonProtocol handles the Json requests
type JsonProtocol struct {
	server              *nex.Server
	ConnectionIDCounter *nex.Counter
	JSONRequestHandler  func(err error, client *nex.Client, callID uint32, rawJson string)
}

// Setup initializes the protocol
func (jsonProtocol *JsonProtocol) Setup() {
	nexServer := jsonProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if JsonProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case JsonRequest:
				go jsonProtocol.handleRequest(packet)
			default:
				fmt.Printf("Unsupported Secure method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// Register sets the Register handler function
func (jsonProtocol *JsonProtocol) JSONRequest(handler func(err error, client *nex.Client, callID uint32, rawJson string)) {
	jsonProtocol.JSONRequestHandler = handler
}

func (jsonProtocol *JsonProtocol) handleRequest(packet nex.PacketInterface) {
	if jsonProtocol.JSONRequestHandler == nil {
		fmt.Println("[Warning] JsonProtocol::JSONRequest not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, jsonProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[SecureProtocol::Register] Json missing length")
		go jsonProtocol.JSONRequestHandler(err, client, callID, "")
		return
	}

	rawJson, err := parametersStream.Read4ByteString()

	if err != nil {
		go jsonProtocol.JSONRequestHandler(err, client, callID, "")
		return
	}

	go jsonProtocol.JSONRequestHandler(nil, client, callID, rawJson)
}

func (jsonProtocol *JsonProtocol) RouteJSONRequest(rawJson string) string {

	var decodedJson [][]interface{}
	json.Unmarshal([]byte(rawJson), &decodedJson)

	// structure is just an array of arrays for most if not all methods and the first is always the method name
	methodName := decodedJson[0][0]

	switch methodName {
	case "config/get":
		content, jsonErr := ioutil.ReadFile("E:/GitHub Projects/GoCentral/motd.json")
		if jsonErr != nil {
			panic(jsonErr)
		}

		return string(content)
	case "misc/get_accounts_web_linked_status":
		return "[[\"misc/get_accounts_web_linked_status\", \"dd\", [\"pid\", \"linked\"], [[12345, 1]]]]"
	case "misc/get_accounts_setlist_creation_status":
		return "[[\"songlists/setlist_creator_status\", \"dd\", [\"pid\", \"creator\"], [[12345, 0]]]]"
	}

	return "A"

}

// NewSecureProtocol returns a new SecureProtocol
func NewJsonProtocol(server *nex.Server) *JsonProtocol {
	jsonProtocol := &JsonProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	jsonProtocol.Setup()

	return jsonProtocol
}
