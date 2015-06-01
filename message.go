package dzpk

type MessageType int

const (
	HoleCards MessageType = 20
	Flop      MessageType = 20 + 1
	Turn      MessageType = 20 + 2
	River     MessageType = 20 + 3

	Check MessageType = 30
	Bet   MessageType = 30 + 1
	Raise MessageType = 30 + 2
	Fold  MessageType = 30 + 3
)

type RequestMessage struct {
	Id          int    `json:"id"`
	MessageType int    `json:"msgtype"`
	Payload     []byte `json:"payload"`
}

type ResponseMessage struct {
	RequestId int    `json:"requestid"`
	Result    int    `json:"result"`
	Payload   []byte `json:"payload"`
}
