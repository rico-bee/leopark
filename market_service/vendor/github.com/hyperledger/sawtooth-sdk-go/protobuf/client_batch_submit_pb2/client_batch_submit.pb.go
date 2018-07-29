// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/client_batch_submit_pb2/client_batch_submit.proto

package client_batch_submit_pb2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import batch_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ClientBatchStatus_Status int32

const (
	ClientBatchStatus_STATUS_UNSET ClientBatchStatus_Status = 0
	ClientBatchStatus_COMMITTED    ClientBatchStatus_Status = 1
	ClientBatchStatus_INVALID      ClientBatchStatus_Status = 2
	ClientBatchStatus_PENDING      ClientBatchStatus_Status = 3
	ClientBatchStatus_UNKNOWN      ClientBatchStatus_Status = 4
)

var ClientBatchStatus_Status_name = map[int32]string{
	0: "STATUS_UNSET",
	1: "COMMITTED",
	2: "INVALID",
	3: "PENDING",
	4: "UNKNOWN",
}
var ClientBatchStatus_Status_value = map[string]int32{
	"STATUS_UNSET": 0,
	"COMMITTED":    1,
	"INVALID":      2,
	"PENDING":      3,
	"UNKNOWN":      4,
}

func (x ClientBatchStatus_Status) String() string {
	return proto.EnumName(ClientBatchStatus_Status_name, int32(x))
}
func (ClientBatchStatus_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{0, 0}
}

type ClientBatchSubmitResponse_Status int32

const (
	ClientBatchSubmitResponse_STATUS_UNSET   ClientBatchSubmitResponse_Status = 0
	ClientBatchSubmitResponse_OK             ClientBatchSubmitResponse_Status = 1
	ClientBatchSubmitResponse_INTERNAL_ERROR ClientBatchSubmitResponse_Status = 2
	ClientBatchSubmitResponse_INVALID_BATCH  ClientBatchSubmitResponse_Status = 3
	ClientBatchSubmitResponse_QUEUE_FULL     ClientBatchSubmitResponse_Status = 4
)

var ClientBatchSubmitResponse_Status_name = map[int32]string{
	0: "STATUS_UNSET",
	1: "OK",
	2: "INTERNAL_ERROR",
	3: "INVALID_BATCH",
	4: "QUEUE_FULL",
}
var ClientBatchSubmitResponse_Status_value = map[string]int32{
	"STATUS_UNSET":   0,
	"OK":             1,
	"INTERNAL_ERROR": 2,
	"INVALID_BATCH":  3,
	"QUEUE_FULL":     4,
}

func (x ClientBatchSubmitResponse_Status) String() string {
	return proto.EnumName(ClientBatchSubmitResponse_Status_name, int32(x))
}
func (ClientBatchSubmitResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{2, 0}
}

type ClientBatchStatusResponse_Status int32

const (
	ClientBatchStatusResponse_STATUS_UNSET   ClientBatchStatusResponse_Status = 0
	ClientBatchStatusResponse_OK             ClientBatchStatusResponse_Status = 1
	ClientBatchStatusResponse_INTERNAL_ERROR ClientBatchStatusResponse_Status = 2
	ClientBatchStatusResponse_NO_RESOURCE    ClientBatchStatusResponse_Status = 5
	ClientBatchStatusResponse_INVALID_ID     ClientBatchStatusResponse_Status = 8
)

var ClientBatchStatusResponse_Status_name = map[int32]string{
	0: "STATUS_UNSET",
	1: "OK",
	2: "INTERNAL_ERROR",
	5: "NO_RESOURCE",
	8: "INVALID_ID",
}
var ClientBatchStatusResponse_Status_value = map[string]int32{
	"STATUS_UNSET":   0,
	"OK":             1,
	"INTERNAL_ERROR": 2,
	"NO_RESOURCE":    5,
	"INVALID_ID":     8,
}

func (x ClientBatchStatusResponse_Status) String() string {
	return proto.EnumName(ClientBatchStatusResponse_Status_name, int32(x))
}
func (ClientBatchStatusResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{4, 0}
}

// Information about the status of a batch submitted to the validator.
//
// Attributes:
//     batch_id: The id (header_signature) of the batch
//     status: The committed status of the batch
//     invalid_transactions: Info for transactions that failed, if any
//
// Statuses:
//     COMMITTED - the batch was accepted and has been committed to the chain
//     INVALID - the batch failed validation, it should be resubmitted
//     PENDING - the batch is still being processed
//     UNKNOWN - no status for the batch could be found (possibly invalid)
type ClientBatchStatus struct {
	BatchId              string                                  `protobuf:"bytes,1,opt,name=batch_id,json=batchId" json:"batch_id,omitempty"`
	Status               ClientBatchStatus_Status                `protobuf:"varint,2,opt,name=status,enum=ClientBatchStatus_Status" json:"status,omitempty"`
	InvalidTransactions  []*ClientBatchStatus_InvalidTransaction `protobuf:"bytes,3,rep,name=invalid_transactions,json=invalidTransactions" json:"invalid_transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *ClientBatchStatus) Reset()         { *m = ClientBatchStatus{} }
func (m *ClientBatchStatus) String() string { return proto.CompactTextString(m) }
func (*ClientBatchStatus) ProtoMessage()    {}
func (*ClientBatchStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{0}
}
func (m *ClientBatchStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchStatus.Unmarshal(m, b)
}
func (m *ClientBatchStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchStatus.Marshal(b, m, deterministic)
}
func (dst *ClientBatchStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchStatus.Merge(dst, src)
}
func (m *ClientBatchStatus) XXX_Size() int {
	return xxx_messageInfo_ClientBatchStatus.Size(m)
}
func (m *ClientBatchStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchStatus proto.InternalMessageInfo

func (m *ClientBatchStatus) GetBatchId() string {
	if m != nil {
		return m.BatchId
	}
	return ""
}

func (m *ClientBatchStatus) GetStatus() ClientBatchStatus_Status {
	if m != nil {
		return m.Status
	}
	return ClientBatchStatus_STATUS_UNSET
}

func (m *ClientBatchStatus) GetInvalidTransactions() []*ClientBatchStatus_InvalidTransaction {
	if m != nil {
		return m.InvalidTransactions
	}
	return nil
}

type ClientBatchStatus_InvalidTransaction struct {
	TransactionId        string   `protobuf:"bytes,1,opt,name=transaction_id,json=transactionId" json:"transaction_id,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	ExtendedData         []byte   `protobuf:"bytes,3,opt,name=extended_data,json=extendedData,proto3" json:"extended_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClientBatchStatus_InvalidTransaction) Reset()         { *m = ClientBatchStatus_InvalidTransaction{} }
func (m *ClientBatchStatus_InvalidTransaction) String() string { return proto.CompactTextString(m) }
func (*ClientBatchStatus_InvalidTransaction) ProtoMessage()    {}
func (*ClientBatchStatus_InvalidTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{0, 0}
}
func (m *ClientBatchStatus_InvalidTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchStatus_InvalidTransaction.Unmarshal(m, b)
}
func (m *ClientBatchStatus_InvalidTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchStatus_InvalidTransaction.Marshal(b, m, deterministic)
}
func (dst *ClientBatchStatus_InvalidTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchStatus_InvalidTransaction.Merge(dst, src)
}
func (m *ClientBatchStatus_InvalidTransaction) XXX_Size() int {
	return xxx_messageInfo_ClientBatchStatus_InvalidTransaction.Size(m)
}
func (m *ClientBatchStatus_InvalidTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchStatus_InvalidTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchStatus_InvalidTransaction proto.InternalMessageInfo

func (m *ClientBatchStatus_InvalidTransaction) GetTransactionId() string {
	if m != nil {
		return m.TransactionId
	}
	return ""
}

func (m *ClientBatchStatus_InvalidTransaction) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ClientBatchStatus_InvalidTransaction) GetExtendedData() []byte {
	if m != nil {
		return m.ExtendedData
	}
	return nil
}

// Submits a list of Batches to be added to the blockchain.
type ClientBatchSubmitRequest struct {
	Batches              []*batch_pb2.Batch `protobuf:"bytes,1,rep,name=batches" json:"batches,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ClientBatchSubmitRequest) Reset()         { *m = ClientBatchSubmitRequest{} }
func (m *ClientBatchSubmitRequest) String() string { return proto.CompactTextString(m) }
func (*ClientBatchSubmitRequest) ProtoMessage()    {}
func (*ClientBatchSubmitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{1}
}
func (m *ClientBatchSubmitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchSubmitRequest.Unmarshal(m, b)
}
func (m *ClientBatchSubmitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchSubmitRequest.Marshal(b, m, deterministic)
}
func (dst *ClientBatchSubmitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchSubmitRequest.Merge(dst, src)
}
func (m *ClientBatchSubmitRequest) XXX_Size() int {
	return xxx_messageInfo_ClientBatchSubmitRequest.Size(m)
}
func (m *ClientBatchSubmitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchSubmitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchSubmitRequest proto.InternalMessageInfo

func (m *ClientBatchSubmitRequest) GetBatches() []*batch_pb2.Batch {
	if m != nil {
		return m.Batches
	}
	return nil
}

// This is a response to a submission of one or more Batches.
// Statuses:
//   * OK - everything with the request worked as expected
//   * INTERNAL_ERROR - general error, such as protobuf failing to deserialize
//   * INVALID_BATCH - the batch failed validation, likely due to a bad signature
//   * QUEUE_FULL - the batch is unable to be queued for processing, due to
//        a full processing queue.  The batch may be submitted again.
type ClientBatchSubmitResponse struct {
	Status               ClientBatchSubmitResponse_Status `protobuf:"varint,1,opt,name=status,enum=ClientBatchSubmitResponse_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ClientBatchSubmitResponse) Reset()         { *m = ClientBatchSubmitResponse{} }
func (m *ClientBatchSubmitResponse) String() string { return proto.CompactTextString(m) }
func (*ClientBatchSubmitResponse) ProtoMessage()    {}
func (*ClientBatchSubmitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{2}
}
func (m *ClientBatchSubmitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchSubmitResponse.Unmarshal(m, b)
}
func (m *ClientBatchSubmitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchSubmitResponse.Marshal(b, m, deterministic)
}
func (dst *ClientBatchSubmitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchSubmitResponse.Merge(dst, src)
}
func (m *ClientBatchSubmitResponse) XXX_Size() int {
	return xxx_messageInfo_ClientBatchSubmitResponse.Size(m)
}
func (m *ClientBatchSubmitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchSubmitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchSubmitResponse proto.InternalMessageInfo

func (m *ClientBatchSubmitResponse) GetStatus() ClientBatchSubmitResponse_Status {
	if m != nil {
		return m.Status
	}
	return ClientBatchSubmitResponse_STATUS_UNSET
}

// A request for the status of one or more batches, specified by id.
// If `wait` is set to true, the validator will wait to respond until all
// batches are committed, or until the specified `timeout` in seconds has
// elapsed. Defaults to 300.
type ClientBatchStatusRequest struct {
	BatchIds             []string `protobuf:"bytes,1,rep,name=batch_ids,json=batchIds" json:"batch_ids,omitempty"`
	Wait                 bool     `protobuf:"varint,2,opt,name=wait" json:"wait,omitempty"`
	Timeout              uint32   `protobuf:"varint,3,opt,name=timeout" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClientBatchStatusRequest) Reset()         { *m = ClientBatchStatusRequest{} }
func (m *ClientBatchStatusRequest) String() string { return proto.CompactTextString(m) }
func (*ClientBatchStatusRequest) ProtoMessage()    {}
func (*ClientBatchStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{3}
}
func (m *ClientBatchStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchStatusRequest.Unmarshal(m, b)
}
func (m *ClientBatchStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchStatusRequest.Marshal(b, m, deterministic)
}
func (dst *ClientBatchStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchStatusRequest.Merge(dst, src)
}
func (m *ClientBatchStatusRequest) XXX_Size() int {
	return xxx_messageInfo_ClientBatchStatusRequest.Size(m)
}
func (m *ClientBatchStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchStatusRequest proto.InternalMessageInfo

func (m *ClientBatchStatusRequest) GetBatchIds() []string {
	if m != nil {
		return m.BatchIds
	}
	return nil
}

func (m *ClientBatchStatusRequest) GetWait() bool {
	if m != nil {
		return m.Wait
	}
	return false
}

func (m *ClientBatchStatusRequest) GetTimeout() uint32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

// This is a response to a request for the status of specific batches.
// Statuses:
//   * OK - everything with the request worked as expected
//   * INTERNAL_ERROR - general error, such as protobuf failing to deserialize
//   * NO_RESOURCE - the response contains no data, likely because
//     no ids were specified in the request
type ClientBatchStatusResponse struct {
	Status               ClientBatchStatusResponse_Status `protobuf:"varint,1,opt,name=status,enum=ClientBatchStatusResponse_Status" json:"status,omitempty"`
	BatchStatuses        []*ClientBatchStatus             `protobuf:"bytes,2,rep,name=batch_statuses,json=batchStatuses" json:"batch_statuses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ClientBatchStatusResponse) Reset()         { *m = ClientBatchStatusResponse{} }
func (m *ClientBatchStatusResponse) String() string { return proto.CompactTextString(m) }
func (*ClientBatchStatusResponse) ProtoMessage()    {}
func (*ClientBatchStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_batch_submit_bf71acee74810670, []int{4}
}
func (m *ClientBatchStatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientBatchStatusResponse.Unmarshal(m, b)
}
func (m *ClientBatchStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientBatchStatusResponse.Marshal(b, m, deterministic)
}
func (dst *ClientBatchStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientBatchStatusResponse.Merge(dst, src)
}
func (m *ClientBatchStatusResponse) XXX_Size() int {
	return xxx_messageInfo_ClientBatchStatusResponse.Size(m)
}
func (m *ClientBatchStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientBatchStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClientBatchStatusResponse proto.InternalMessageInfo

func (m *ClientBatchStatusResponse) GetStatus() ClientBatchStatusResponse_Status {
	if m != nil {
		return m.Status
	}
	return ClientBatchStatusResponse_STATUS_UNSET
}

func (m *ClientBatchStatusResponse) GetBatchStatuses() []*ClientBatchStatus {
	if m != nil {
		return m.BatchStatuses
	}
	return nil
}

func init() {
	proto.RegisterType((*ClientBatchStatus)(nil), "ClientBatchStatus")
	proto.RegisterType((*ClientBatchStatus_InvalidTransaction)(nil), "ClientBatchStatus.InvalidTransaction")
	proto.RegisterType((*ClientBatchSubmitRequest)(nil), "ClientBatchSubmitRequest")
	proto.RegisterType((*ClientBatchSubmitResponse)(nil), "ClientBatchSubmitResponse")
	proto.RegisterType((*ClientBatchStatusRequest)(nil), "ClientBatchStatusRequest")
	proto.RegisterType((*ClientBatchStatusResponse)(nil), "ClientBatchStatusResponse")
	proto.RegisterEnum("ClientBatchStatus_Status", ClientBatchStatus_Status_name, ClientBatchStatus_Status_value)
	proto.RegisterEnum("ClientBatchSubmitResponse_Status", ClientBatchSubmitResponse_Status_name, ClientBatchSubmitResponse_Status_value)
	proto.RegisterEnum("ClientBatchStatusResponse_Status", ClientBatchStatusResponse_Status_name, ClientBatchStatusResponse_Status_value)
}

func init() {
	proto.RegisterFile("github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2/client_batch_submit.proto", fileDescriptor_client_batch_submit_bf71acee74810670)
}

var fileDescriptor_client_batch_submit_bf71acee74810670 = []byte{
	// 555 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdd, 0x6e, 0xd3, 0x30,
	0x18, 0x25, 0xcd, 0xe8, 0xda, 0x6f, 0x4d, 0xf1, 0x0c, 0x88, 0x74, 0x48, 0xa8, 0x04, 0x4d, 0xea,
	0x55, 0x10, 0xe5, 0x6a, 0x88, 0x9b, 0xfe, 0x04, 0x88, 0xd6, 0x25, 0xc5, 0x4d, 0x18, 0x70, 0x13,
	0x39, 0x8d, 0x61, 0x11, 0x6b, 0x52, 0x66, 0x97, 0x21, 0xde, 0x80, 0xc7, 0xe1, 0xb5, 0x78, 0x0a,
	0x14, 0x27, 0x59, 0xbb, 0xb6, 0x43, 0x88, 0xab, 0xe4, 0x7c, 0x3e, 0x9f, 0xad, 0x73, 0xce, 0x67,
	0xc3, 0x8b, 0xf9, 0x45, 0x2a, 0xd2, 0x70, 0xf1, 0xe9, 0xe9, 0xf4, 0x3c, 0x66, 0x89, 0x08, 0x42,
	0x2a, 0xa6, 0x67, 0x01, 0x5f, 0x84, 0xb3, 0x58, 0x04, 0xf3, 0xb0, 0xbb, 0xad, 0x6e, 0xca, 0xa6,
	0x83, 0x47, 0x57, 0xbd, 0xf9, 0x62, 0xc6, 0x96, 0x7f, 0xf9, 0xba, 0xf1, 0x53, 0x85, 0xfd, 0x81,
	0xec, 0xee, 0x67, 0xd5, 0x89, 0xa0, 0x62, 0xc1, 0x71, 0x0b, 0x6a, 0x39, 0x3d, 0x8e, 0x74, 0xa5,
	0xad, 0x74, 0xea, 0x64, 0x57, 0x62, 0x3b, 0xc2, 0xcf, 0xa0, 0xca, 0x25, 0x49, 0xaf, 0xb4, 0x95,
	0x4e, 0xb3, 0xdb, 0x32, 0x37, 0xda, 0xcd, 0xfc, 0x43, 0x0a, 0x22, 0x7e, 0x0f, 0xf7, 0xe2, 0xe4,
	0x1b, 0x3d, 0x8f, 0xa3, 0x40, 0x5c, 0xd0, 0x84, 0xd3, 0xa9, 0x88, 0xd3, 0x84, 0xeb, 0x6a, 0x5b,
	0xed, 0xec, 0x75, 0x0f, 0xb7, 0x6c, 0x60, 0xe7, 0x74, 0x6f, 0xc9, 0x26, 0x77, 0xe3, 0x8d, 0x1a,
	0x3f, 0xf8, 0x01, 0x78, 0x93, 0x8a, 0x0f, 0xa1, 0xb9, 0x72, 0xce, 0x52, 0x83, 0xb6, 0x52, 0xb5,
	0x23, 0xac, 0xc3, 0xee, 0x8c, 0x71, 0x4e, 0x3f, 0x33, 0x29, 0xa5, 0x4e, 0x4a, 0x88, 0x9f, 0x80,
	0xc6, 0xbe, 0x0b, 0x96, 0x44, 0x2c, 0x0a, 0x22, 0x2a, 0xa8, 0xae, 0xb6, 0x95, 0x4e, 0x83, 0x34,
	0xca, 0xe2, 0x90, 0x0a, 0x6a, 0x8c, 0xa1, 0x5a, 0xb8, 0x85, 0xa0, 0x31, 0xf1, 0x7a, 0x9e, 0x3f,
	0x09, 0x7c, 0x67, 0x62, 0x79, 0xe8, 0x16, 0xd6, 0xa0, 0x3e, 0x70, 0x4f, 0x4e, 0x6c, 0xcf, 0xb3,
	0x86, 0x48, 0xc1, 0x7b, 0xb0, 0x6b, 0x3b, 0xef, 0x7a, 0x23, 0x7b, 0x88, 0x2a, 0x19, 0x18, 0x5b,
	0xce, 0xd0, 0x76, 0x5e, 0x23, 0x35, 0x03, 0xbe, 0x73, 0xec, 0xb8, 0xa7, 0x0e, 0xda, 0x31, 0x5e,
	0x82, 0xbe, 0x6a, 0x85, 0x8c, 0x91, 0xb0, 0xaf, 0x0b, 0xc6, 0x05, 0x6e, 0x43, 0x9e, 0x00, 0xe3,
	0xba, 0x22, 0x6d, 0xab, 0x9a, 0x92, 0x45, 0xca, 0xb2, 0xf1, 0x4b, 0x81, 0xd6, 0x96, 0x76, 0x3e,
	0x4f, 0x13, 0xce, 0xf0, 0xd1, 0x55, 0x6c, 0x8a, 0x8c, 0xed, 0xb1, 0x79, 0x23, 0x77, 0x2d, 0x3e,
	0xe3, 0xc3, 0x5f, 0x84, 0x56, 0xa1, 0xe2, 0x1e, 0x23, 0x05, 0x63, 0x68, 0xda, 0x8e, 0x67, 0x11,
	0xa7, 0x37, 0x0a, 0x2c, 0x42, 0x5c, 0x82, 0x2a, 0x78, 0x1f, 0xb4, 0x42, 0x75, 0xd0, 0xef, 0x79,
	0x83, 0x37, 0x48, 0xc5, 0x4d, 0x80, 0xb7, 0xbe, 0xe5, 0x5b, 0xc1, 0x2b, 0x7f, 0x34, 0x42, 0x3b,
	0x06, 0xbb, 0xae, 0x38, 0x3f, 0xb7, 0x50, 0xfc, 0x10, 0xea, 0xe5, 0x0c, 0xe6, 0x9a, 0xeb, 0xa4,
	0x56, 0x0c, 0x21, 0xc7, 0x18, 0x76, 0x2e, 0x69, 0x2c, 0x64, 0x70, 0x35, 0x22, 0xff, 0xb3, 0x3c,
	0x45, 0x3c, 0x63, 0xe9, 0x42, 0xc8, 0xbc, 0x34, 0x52, 0x42, 0xe3, 0xf7, 0x9a, 0x35, 0xc5, 0x39,
	0xff, 0x64, 0xcd, 0x35, 0xee, 0xfa, 0x64, 0x1f, 0x41, 0xb3, 0xb8, 0x73, 0x12, 0xb3, 0xec, 0x52,
	0x64, 0xe1, 0xe0, 0x2d, 0x5b, 0x68, 0xe1, 0x12, 0x30, 0x6e, 0x9c, 0xfe, 0xa7, 0xab, 0x77, 0x60,
	0xcf, 0x71, 0x03, 0x62, 0x4d, 0x5c, 0x9f, 0x0c, 0x2c, 0x74, 0x3b, 0xf3, 0xb4, 0xb4, 0xd9, 0x1e,
	0xa2, 0x5a, 0xbf, 0x0b, 0xf7, 0x39, 0xbd, 0x14, 0x69, 0x2a, 0xce, 0x4c, 0x1e, 0x7d, 0x31, 0xcb,
	0x07, 0x60, 0xac, 0x7c, 0x7c, 0x70, 0xc3, 0xfb, 0x11, 0x56, 0x25, 0xe9, 0xf9, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x25, 0x7b, 0x7e, 0xa1, 0x6a, 0x04, 0x00, 0x00,
}