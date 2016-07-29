// Code generated by protoc-gen-go.
// source: project.proto
// DO NOT EDIT!

/*
Package crypto_pb is a generated protocol buffer package.

It is generated from these files:
	project.proto

It has these top-level messages:
	ProjectOperation
	Operation
	Credential
	Project
	ProjectOperationResponse
	Response
*/
package crypto_pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ProjectOperation_Command int32

const (
	ProjectOperation_LIST              ProjectOperation_Command = 0
	ProjectOperation_CREATE            ProjectOperation_Command = 1
	ProjectOperation_UPDATE            ProjectOperation_Command = 2
	ProjectOperation_DELETE            ProjectOperation_Command = 3
	ProjectOperation_LIST_CREDENTIALS  ProjectOperation_Command = 4
	ProjectOperation_ADD_MEMBER        ProjectOperation_Command = 5
	ProjectOperation_DELETE_MEMBER     ProjectOperation_Command = 6
	ProjectOperation_ADD_CREDENTIAL    ProjectOperation_Command = 7
	ProjectOperation_DELETE_CREDENTIAL ProjectOperation_Command = 8
	ProjectOperation_GET_CREDENTIAL    ProjectOperation_Command = 9
)

var ProjectOperation_Command_name = map[int32]string{
	0: "LIST",
	1: "CREATE",
	2: "UPDATE",
	3: "DELETE",
	4: "LIST_CREDENTIALS",
	5: "ADD_MEMBER",
	6: "DELETE_MEMBER",
	7: "ADD_CREDENTIAL",
	8: "DELETE_CREDENTIAL",
	9: "GET_CREDENTIAL",
}
var ProjectOperation_Command_value = map[string]int32{
	"LIST":              0,
	"CREATE":            1,
	"UPDATE":            2,
	"DELETE":            3,
	"LIST_CREDENTIALS":  4,
	"ADD_MEMBER":        5,
	"DELETE_MEMBER":     6,
	"ADD_CREDENTIAL":    7,
	"DELETE_CREDENTIAL": 8,
	"GET_CREDENTIAL":    9,
}

func (x ProjectOperation_Command) String() string {
	return proto.EnumName(ProjectOperation_Command_name, int32(x))
}
func (ProjectOperation_Command) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Response_Status int32

const (
	Response_ERROR   Response_Status = 0
	Response_SUCCESS Response_Status = 1
)

var Response_Status_name = map[int32]string{
	0: "ERROR",
	1: "SUCCESS",
}
var Response_Status_value = map[string]int32{
	"ERROR":   0,
	"SUCCESS": 1,
}

func (x Response_Status) String() string {
	return proto.EnumName(Response_Status_name, int32(x))
}
func (Response_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

type ProjectOperation struct {
	Command     ProjectOperation_Command `protobuf:"varint,1,opt,name=command,enum=crypto_pb.ProjectOperation_Command" json:"command,omitempty"`
	Name        string                   `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Environment string                   `protobuf:"bytes,3,opt,name=environment" json:"environment,omitempty"`
	ProjectId   int32                    `protobuf:"varint,4,opt,name=projectId" json:"projectId,omitempty"`
	MemberId    int32                    `protobuf:"varint,5,opt,name=memberId" json:"memberId,omitempty"`
	UserId      int32                    `protobuf:"varint,6,opt,name=userId" json:"userId,omitempty"`
	AccessLevel string                   `protobuf:"bytes,7,opt,name=accessLevel" json:"accessLevel,omitempty"`
	MemberEmail string                   `protobuf:"bytes,8,opt,name=memberEmail" json:"memberEmail,omitempty"`
	Key         string                   `protobuf:"bytes,9,opt,name=key" json:"key,omitempty"`
	Value       string                   `protobuf:"bytes,10,opt,name=value" json:"value,omitempty"`
}

func (m *ProjectOperation) Reset()                    { *m = ProjectOperation{} }
func (m *ProjectOperation) String() string            { return proto.CompactTextString(m) }
func (*ProjectOperation) ProtoMessage()               {}
func (*ProjectOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Operation struct {
	OpId      int32             `protobuf:"varint,1,opt,name=opId" json:"opId,omitempty"`
	ProjectOp *ProjectOperation `protobuf:"bytes,2,opt,name=projectOp" json:"projectOp,omitempty"`
}

func (m *Operation) Reset()                    { *m = Operation{} }
func (m *Operation) String() string            { return proto.CompactTextString(m) }
func (*Operation) ProtoMessage()               {}
func (*Operation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Operation) GetProjectOp() *ProjectOperation {
	if m != nil {
		return m.ProjectOp
	}
	return nil
}

type Credential struct {
	Id     int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Cipher string `protobuf:"bytes,3,opt,name=cipher" json:"cipher,omitempty"`
}

func (m *Credential) Reset()                    { *m = Credential{} }
func (m *Credential) String() string            { return proto.CompactTextString(m) }
func (*Credential) ProtoMessage()               {}
func (*Credential) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Project struct {
	Id          int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Environment string `protobuf:"bytes,3,opt,name=environment" json:"environment,omitempty"`
}

func (m *Project) Reset()                    { *m = Project{} }
func (m *Project) String() string            { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()               {}
func (*Project) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ProjectOperationResponse struct {
	Command     ProjectOperation_Command `protobuf:"varint,1,opt,name=command,enum=crypto_pb.ProjectOperation_Command" json:"command,omitempty"`
	MemberId    int32                    `protobuf:"varint,3,opt,name=memberId" json:"memberId,omitempty"`
	Project     *Project                 `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
	Credential  *Credential              `protobuf:"bytes,6,opt,name=credential" json:"credential,omitempty"`
	Credentials []*Credential            `protobuf:"bytes,4,rep,name=credentials" json:"credentials,omitempty"`
	Projects    []*Project               `protobuf:"bytes,5,rep,name=projects" json:"projects,omitempty"`
}

func (m *ProjectOperationResponse) Reset()                    { *m = ProjectOperationResponse{} }
func (m *ProjectOperationResponse) String() string            { return proto.CompactTextString(m) }
func (*ProjectOperationResponse) ProtoMessage()               {}
func (*ProjectOperationResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ProjectOperationResponse) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

func (m *ProjectOperationResponse) GetCredential() *Credential {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (m *ProjectOperationResponse) GetCredentials() []*Credential {
	if m != nil {
		return m.Credentials
	}
	return nil
}

func (m *ProjectOperationResponse) GetProjects() []*Project {
	if m != nil {
		return m.Projects
	}
	return nil
}

type Response struct {
	Status            Response_Status           `protobuf:"varint,1,opt,name=status,enum=crypto_pb.Response_Status" json:"status,omitempty"`
	Error             string                    `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	Info              string                    `protobuf:"bytes,3,opt,name=info" json:"info,omitempty"`
	OpId              int32                     `protobuf:"varint,4,opt,name=opId" json:"opId,omitempty"`
	ProjectOpResponse *ProjectOperationResponse `protobuf:"bytes,5,opt,name=projectOpResponse" json:"projectOpResponse,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Response) GetProjectOpResponse() *ProjectOperationResponse {
	if m != nil {
		return m.ProjectOpResponse
	}
	return nil
}

func init() {
	proto.RegisterType((*ProjectOperation)(nil), "crypto_pb.ProjectOperation")
	proto.RegisterType((*Operation)(nil), "crypto_pb.Operation")
	proto.RegisterType((*Credential)(nil), "crypto_pb.Credential")
	proto.RegisterType((*Project)(nil), "crypto_pb.Project")
	proto.RegisterType((*ProjectOperationResponse)(nil), "crypto_pb.ProjectOperationResponse")
	proto.RegisterType((*Response)(nil), "crypto_pb.Response")
	proto.RegisterEnum("crypto_pb.ProjectOperation_Command", ProjectOperation_Command_name, ProjectOperation_Command_value)
	proto.RegisterEnum("crypto_pb.Response_Status", Response_Status_name, Response_Status_value)
}

func init() { proto.RegisterFile("project.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 590 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x54, 0xbb, 0x8e, 0xd3, 0x40,
	0x14, 0xc5, 0x71, 0xfc, 0xba, 0xd6, 0xae, 0x9c, 0x11, 0x8b, 0xac, 0x85, 0x22, 0x32, 0x0d, 0x05,
	0x72, 0x11, 0x84, 0x10, 0x05, 0x45, 0x48, 0x0c, 0x8a, 0x94, 0x25, 0xcb, 0x38, 0xdb, 0xd0, 0x44,
	0x8e, 0x33, 0x08, 0x43, 0xfc, 0x90, 0xed, 0x44, 0xda, 0xaf, 0xe2, 0x23, 0xf8, 0x22, 0x2a, 0x5a,
	0x66, 0xc6, 0x63, 0x7b, 0xb4, 0xd1, 0x6e, 0x01, 0x55, 0xee, 0xe3, 0xdc, 0xb9, 0x93, 0x73, 0xce,
	0x18, 0xce, 0x8a, 0x32, 0xff, 0x4e, 0xe2, 0xda, 0xa7, 0xbf, 0x75, 0x8e, 0xac, 0xb8, 0xbc, 0x2d,
	0xea, 0x7c, 0x53, 0x6c, 0xbd, 0x3f, 0x2a, 0x38, 0xd7, 0x4d, 0x73, 0x55, 0x90, 0x32, 0xaa, 0x93,
	0x3c, 0x43, 0xef, 0xc0, 0x88, 0xf3, 0x34, 0x8d, 0xb2, 0x9d, 0xab, 0x8c, 0x95, 0x17, 0xe7, 0x93,
	0xe7, 0x7e, 0x37, 0xe1, 0xdf, 0x45, 0xfb, 0xb3, 0x06, 0x8a, 0xdb, 0x19, 0x84, 0x60, 0x98, 0x45,
	0x29, 0x71, 0x07, 0x74, 0xd6, 0xc2, 0x3c, 0x46, 0x63, 0xb0, 0x49, 0x76, 0x4c, 0xca, 0x3c, 0x4b,
	0x49, 0x56, 0xbb, 0x2a, 0x6f, 0xc9, 0x25, 0xf4, 0x0c, 0x2c, 0x71, 0xcb, 0xc5, 0xce, 0x1d, 0xd2,
	0xbe, 0x86, 0xfb, 0x02, 0xba, 0x04, 0x33, 0x25, 0xe9, 0x96, 0x94, 0xb4, 0xa9, 0xf1, 0x66, 0x97,
	0xa3, 0x27, 0xa0, 0x1f, 0x2a, 0xde, 0xd1, 0x79, 0x47, 0x64, 0x6c, 0x67, 0x14, 0xc7, 0xa4, 0xaa,
	0x96, 0xe4, 0x48, 0xf6, 0xae, 0xd1, 0xec, 0x94, 0x4a, 0x0c, 0xd1, 0x9c, 0x12, 0xa4, 0x51, 0xb2,
	0x77, 0xcd, 0x06, 0x21, 0x95, 0x90, 0x03, 0xea, 0x0f, 0x72, 0xeb, 0x5a, 0xbc, 0xc3, 0x42, 0xf4,
	0x18, 0xb4, 0x63, 0xb4, 0x3f, 0x10, 0x17, 0x78, 0xad, 0x49, 0xbc, 0x9f, 0x0a, 0x18, 0x82, 0x08,
	0x64, 0xc2, 0x70, 0xb9, 0x08, 0xd7, 0xce, 0x23, 0x04, 0xa0, 0xcf, 0x70, 0x30, 0x5d, 0x07, 0x8e,
	0xc2, 0xe2, 0x9b, 0xeb, 0x39, 0x8b, 0x07, 0x2c, 0x9e, 0x07, 0xcb, 0x80, 0xc6, 0x2a, 0x3d, 0xcf,
	0x61, 0xe8, 0x0d, 0x05, 0xce, 0x83, 0x4f, 0xeb, 0xc5, 0x74, 0x19, 0x3a, 0x43, 0x74, 0x0e, 0x30,
	0x9d, 0xcf, 0x37, 0x57, 0xc1, 0xd5, 0xfb, 0x00, 0x3b, 0x1a, 0x1a, 0xc1, 0x59, 0x33, 0xd1, 0x96,
	0x74, 0x4a, 0xf3, 0x39, 0x83, 0xf4, 0x73, 0x8e, 0x81, 0x2e, 0x60, 0x24, 0x60, 0x52, 0xd9, 0x64,
	0xd0, 0x8f, 0x81, 0xbc, 0xc2, 0xb1, 0xbc, 0x2f, 0x60, 0xf5, 0x8a, 0x53, 0xc9, 0xf2, 0x62, 0xd1,
	0xc8, 0xad, 0x61, 0x1e, 0xa3, 0xb7, 0x9d, 0x20, 0xab, 0x82, 0x6b, 0x69, 0x4f, 0x9e, 0x3e, 0xe0,
	0x03, 0xdc, 0xa3, 0xbd, 0x0f, 0x00, 0xb3, 0x92, 0xec, 0xa8, 0xac, 0x49, 0xb4, 0xa7, 0xff, 0x65,
	0x90, 0xb4, 0x47, 0xd3, 0xa8, 0xe5, 0x74, 0xd0, 0x73, 0x4a, 0x15, 0x8c, 0x93, 0xe2, 0x1b, 0x29,
	0x85, 0x31, 0x44, 0xe6, 0xad, 0xc0, 0x10, 0x6b, 0x4e, 0x0e, 0xf9, 0x27, 0x93, 0x79, 0xbf, 0x06,
	0xe0, 0x9e, 0x5c, 0x9c, 0x54, 0x45, 0x9e, 0x55, 0xe4, 0x7f, 0x6d, 0x2f, 0x5b, 0x54, 0xbd, 0x63,
	0xd1, 0x97, 0x60, 0x08, 0x76, 0x04, 0x93, 0xe8, 0xf4, 0x68, 0xdc, 0x42, 0xd0, 0x6b, 0x80, 0xb8,
	0xa3, 0x8f, 0x9b, 0xda, 0x9e, 0x5c, 0x48, 0x03, 0x3d, 0xb7, 0x58, 0x02, 0xa2, 0x37, 0x60, 0xf7,
	0x59, 0x45, 0xdf, 0x90, 0x7a, 0xff, 0x9c, 0x8c, 0x44, 0x3e, 0x98, 0x62, 0x75, 0x45, 0x1f, 0x97,
	0x7a, 0xcf, 0xf5, 0x3a, 0x8c, 0xf7, 0x5b, 0x01, 0xb3, 0x63, 0x6d, 0x02, 0x7a, 0x55, 0x47, 0xf5,
	0xa1, 0x12, 0xa4, 0x5d, 0x4a, 0xa3, 0x2d, 0xc8, 0x0f, 0x39, 0x02, 0x0b, 0x24, 0x7b, 0x43, 0xa4,
	0x2c, 0xf3, 0x52, 0xa8, 0xd7, 0x24, 0x4c, 0xd2, 0x24, 0xfb, 0x9a, 0x0b, 0xdd, 0x78, 0xdc, 0x19,
	0x73, 0x28, 0x19, 0xf3, 0x33, 0x8c, 0x3a, 0xab, 0xb5, 0x1b, 0xf8, 0x47, 0xc1, 0x7e, 0x50, 0xb1,
	0x16, 0x8a, 0x4f, 0xa7, 0xbd, 0x31, 0xe8, 0xcd, 0x15, 0x91, 0x05, 0x5a, 0x80, 0xf1, 0x0a, 0xd3,
	0xd7, 0x6b, 0x83, 0x11, 0xde, 0xcc, 0x66, 0x41, 0x18, 0x3a, 0xca, 0x56, 0xe7, 0x9f, 0xce, 0x57,
	0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x25, 0xe9, 0x87, 0xe0, 0x4b, 0x05, 0x00, 0x00,
}
