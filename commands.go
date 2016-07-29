package main

import (
	"errors"
	pb "github.com/rajivnavada/cryptz_pb"
	"strings"
)

var (
	ErrInvalidProjectName     = errors.New("provided project name is invalid")
	ErrInvalidProjectId       = errors.New("project id is invalid")
	ErrInvalidCredentialKey   = errors.New("credential key is invalid")
	ErrInvalidCredentialValue = errors.New("credential value is invalid")
	ErrInvalidMemberEmail     = errors.New("member email is invalid")
	ErrInvalidMemberId        = errors.New("member id is invalid")
)

func NewProjectOperation() *pb.ProjectOperation {
	return &pb.ProjectOperation{}
}

func NewProjectListOperation() *pb.Operation {
	o := NewProjectOperation()
	o.Command = pb.ProjectOperation_LIST
	return &pb.Operation{ProjectOp: o}
}

func NewProjectCreateOperation(name, env string) (*pb.Operation, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidProjectName
	}

	o := NewProjectOperation()
	o.Command = pb.ProjectOperation_CREATE
	o.Name = name
	o.Environment = env
	return &pb.Operation{ProjectOp: o}, nil
}

//----------------------------------------
// Credential related operations
//----------------------------------------

type Project int

func (p Project) validate() error {
	if p <= 0 {
		return ErrInvalidProjectId
	}
	return nil
}

func (p Project) newOperation() (*pb.ProjectOperation, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}
	o := NewProjectOperation()
	o.ProjectId = int32(p)
	return o, nil
}

func (p Project) NewDeleteOperation() (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}

	o.Command = pb.ProjectOperation_DELETE
	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewListCredentialsOperation() (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}

	o.Command = pb.ProjectOperation_LIST_CREDENTIALS
	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewAddCredentialOperation(key, value string) (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}
	o.Command = pb.ProjectOperation_ADD_CREDENTIAL

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidCredentialKey
	}
	o.Key = key

	value = strings.TrimSpace(value)
	if value == "" {
		return nil, ErrInvalidCredentialValue
	}
	o.Value = value

	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewGetCredentialOperation(key string) (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}
	o.Command = pb.ProjectOperation_GET_CREDENTIAL

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidCredentialKey
	}
	o.Key = key

	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewDeleteCredentialOperation(key string) (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}
	o.Command = pb.ProjectOperation_DELETE_CREDENTIAL

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidCredentialKey
	}
	o.Key = key

	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewAddMemberOperation(email string) (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}
	o.Command = pb.ProjectOperation_ADD_MEMBER

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, ErrInvalidMemberEmail
	}
	o.MemberEmail = email

	return &pb.Operation{ProjectOp: o}, nil
}

func (p Project) NewDeleteMemberOperation(memberId int) (*pb.Operation, error) {
	o, err := p.newOperation()
	if err != nil {
		return nil, err
	}
	o.Command = pb.ProjectOperation_ADD_MEMBER

	if memberId <= 0 {
		return nil, ErrInvalidMemberId
	}
	o.MemberId = int32(memberId)

	return &pb.Operation{ProjectOp: o}, nil
}
