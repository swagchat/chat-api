package models

// ReceivedMessage: A message and its corresponding acknowledgment ID.
type ReceivedMessage struct {
	// AckId: This ID can be used to acknowledge the received message.
	AckId string `json:"ackId,omitempty"`

	// Message: The message.
	Message *PubsubMessage `json:"message,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AckId") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AckId") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

// PubsubMessage: A message data and its attributes. The message payload
// must not be empty; it must contain either a non-empty data field, or
// at least one attribute.
type PubsubMessage struct {
	// Attributes: Optional attributes for this message.
	Attributes map[string]string `json:"attributes,omitempty"`

	// Data: The message payload. For JSON requests, the value of this field
	// must be [base64-encoded](https://tools.ietf.org/html/rfc4648).
	Data string `json:"data,omitempty"`

	// MessageId: ID of this message, assigned by the server when the
	// message is published. Guaranteed to be unique within the topic. This
	// value may be read by a subscriber that receives a `PubsubMessage` via
	// a `Pull` call or a push delivery. It must not be populated by the
	// publisher in a `Publish` call.
	MessageId string `json:"messageId,omitempty"`

	// PublishTime: The time at which the message was published, populated
	// by the server when it receives the `Publish` call. It must not be
	// populated by the publisher in a `Publish` call.
	PublishTime string `json:"publishTime,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Attributes") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Attributes") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}
