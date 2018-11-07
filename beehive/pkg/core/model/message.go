package model

import (
	"time"

	"github.com/satori/go.uuid"
)

const (
	InsertOperation        = "insert"
	DeleteOperation        = "delete"
	QueryOperation         = "query"
	UpdateOperation        = "update"
	ResponseOperation      = "response"
	ResponseErrorOperation = "error"

	ResourceTypePod        = "pod"
	ResourceTypeConfigmap  = "configmap"
	ResourceTypeSecret     = "secret"
	ResourceTypeNode       = "node"
	ResourceTypePodlist    = "podlist"
	ResourceTypePodStatus  = "podstatus"
	ResourceTypeNodeStatus = "nodestatus"
)

// Message struct
type Message struct {
	Header  MessageHeader `json:"header"`
	Router  MessageRoute  `json:"route, omitempty"`
	Content interface{}   `json:"content"`
}

type MessageRoute struct {
	// where the message come from
	Source string `json:"source,omitempty"`
	// where the message will broadcasted to
	Group string `json:"group, omitempty"`

	// what's the operation on resource
	Operation string `json:"operation,omitempty"`
	// what's the resource want to operate
	Resource string `json:"resource,omitempty"`
}

type MessageHeader struct {
	// the messsage uuid
	ID string `json:"msg_id"`
	// the response message parentid must be same with message received
	// please use NewRespByMessage to new response message
	ParentID string `json:"parent_msg_id, omitempty"`
	// the time of creating
	Timestamp int64 `json:"timestamp"`
	// the flag will be set in sendsync
	Sync bool `json:"sync, omitempty"`
}

func (msg *Message) BuildRouter(source, group, res, opr string) *Message {
	msg.SetRoute(source, group)
	msg.SetResourceOperation(res, opr)
	return msg
}

func (msg *Message) SetResourceOperation(res, opr string) *Message {
	msg.Router.Resource = res
	msg.Router.Operation = opr
	return msg
}

func (msg *Message) SetRoute(source, group string) *Message {
	msg.Router.Source = source
	msg.Router.Group = group
	return msg
}

// msg.Header.Sync will be set in sendsync
func (msg *Message) IsSync() bool {
	return msg.Header.Sync
}

func (msg *Message) GetResource() string {
	return msg.Router.Resource
}

func (msg *Message) GetOperation() string {
	return msg.Router.Operation
}

func (msg *Message) GetSource() string {
	return msg.Router.Source
}

func (msg *Message) GetGroup() string {
	return msg.Router.Group
}

func (msg *Message) GetID() string {
	return msg.Header.ID
}

func (msg *Message) GetParentID() string {
	return msg.Header.ParentID
}

func (msg *Message) GetTimestamp() int64 {
	return msg.Header.Timestamp
}

func (msg *Message) GetContent() interface{} {
	return msg.Content
}

func (msg *Message) UpdateID() *Message {
	msg.Header.ID = uuid.NewV4().String()
	return msg
}

// you can also use for updating message header
func (msg *Message) BuildHeader(ID, parentID string, timestamp int64) *Message {
	msg.Header.ID = ID
	msg.Header.ParentID = parentID
	msg.Header.Timestamp = timestamp
	return msg
}

// the content that you want to send
func (msg *Message) FillBody(content interface{}) *Message {
	msg.Content = content
	return msg
}

// new a raw message:
// model.NewRawMessage().BildHeader().BuildRouter().FillBody()
func NewRawMessage() *Message {
	return &Message{}
}

// new a basic message:
// model.NewMessage().BuildRouter().FillBody()
func NewMessage(parentId string) *Message {
	msg := &Message{}
	msg.Header.ID = uuid.NewV4().String()
	msg.Header.ParentID = parentId
	msg.Header.Timestamp = time.Now().UnixNano() / 1e6
	return msg
}

// clone a message
// only update message id
func (msg *Message) Clone(message *Message) *Message {
	msgID := uuid.NewV4().String()
	return NewRawMessage().BuildHeader(msgID, message.GetParentID(), message.GetTimestamp()).
		BuildRouter(message.GetSource(), message.GetGroup(), message.GetResource(), message.GetOperation()).
		FillBody(message.GetContent())
}

// new a response message by a message received
func (msg *Message) NewRespByMessage(message *Message, content interface{}) *Message {
	return NewMessage(message.GetID()).SetRoute(message.GetSource(), message.GetGroup()).
		SetResourceOperation(message.GetResource(), ResponseOperation).
		FillBody(content)
}

// new a error message by a message received
func NewErrorMessage(message *Message, errContent string) *Message {
	return NewMessage(message.Header.ParentID).
		SetResourceOperation(message.Router.Resource, ResponseErrorOperation).
		FillBody(errContent)
}