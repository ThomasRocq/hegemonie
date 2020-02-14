// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package hegemonie_auth_proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type None struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *None) Reset()         { *m = None{} }
func (m *None) String() string { return proto.CompactTextString(m) }
func (*None) ProtoMessage()    {}
func (*None) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *None) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_None.Unmarshal(m, b)
}
func (m *None) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_None.Marshal(b, m, deterministic)
}
func (m *None) XXX_Merge(src proto.Message) {
	xxx_messageInfo_None.Merge(m, src)
}
func (m *None) XXX_Size() int {
	return xxx_messageInfo_None.Size(m)
}
func (m *None) XXX_DiscardUnknown() {
	xxx_messageInfo_None.DiscardUnknown(m)
}

var xxx_messageInfo_None proto.InternalMessageInfo

type UserCreateReq struct {
	Mail                 string   `protobuf:"bytes,1,opt,name=mail,proto3" json:"mail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserCreateReq) Reset()         { *m = UserCreateReq{} }
func (m *UserCreateReq) String() string { return proto.CompactTextString(m) }
func (*UserCreateReq) ProtoMessage()    {}
func (*UserCreateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *UserCreateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserCreateReq.Unmarshal(m, b)
}
func (m *UserCreateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserCreateReq.Marshal(b, m, deterministic)
}
func (m *UserCreateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserCreateReq.Merge(m, src)
}
func (m *UserCreateReq) XXX_Size() int {
	return xxx_messageInfo_UserCreateReq.Size(m)
}
func (m *UserCreateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserCreateReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserCreateReq proto.InternalMessageInfo

func (m *UserCreateReq) GetMail() string {
	if m != nil {
		return m.Mail
	}
	return ""
}

type UserUpdateReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Pass                 string   `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserUpdateReq) Reset()         { *m = UserUpdateReq{} }
func (m *UserUpdateReq) String() string { return proto.CompactTextString(m) }
func (*UserUpdateReq) ProtoMessage()    {}
func (*UserUpdateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{2}
}

func (m *UserUpdateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserUpdateReq.Unmarshal(m, b)
}
func (m *UserUpdateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserUpdateReq.Marshal(b, m, deterministic)
}
func (m *UserUpdateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserUpdateReq.Merge(m, src)
}
func (m *UserUpdateReq) XXX_Size() int {
	return xxx_messageInfo_UserUpdateReq.Size(m)
}
func (m *UserUpdateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserUpdateReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserUpdateReq proto.InternalMessageInfo

func (m *UserUpdateReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserUpdateReq) GetPass() string {
	if m != nil {
		return m.Pass
	}
	return ""
}

func (m *UserUpdateReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type UserAuthReq struct {
	Mail                 string   `protobuf:"bytes,1,opt,name=mail,proto3" json:"mail,omitempty"`
	Pass                 string   `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserAuthReq) Reset()         { *m = UserAuthReq{} }
func (m *UserAuthReq) String() string { return proto.CompactTextString(m) }
func (*UserAuthReq) ProtoMessage()    {}
func (*UserAuthReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{3}
}

func (m *UserAuthReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserAuthReq.Unmarshal(m, b)
}
func (m *UserAuthReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserAuthReq.Marshal(b, m, deterministic)
}
func (m *UserAuthReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserAuthReq.Merge(m, src)
}
func (m *UserAuthReq) XXX_Size() int {
	return xxx_messageInfo_UserAuthReq.Size(m)
}
func (m *UserAuthReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserAuthReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserAuthReq proto.InternalMessageInfo

func (m *UserAuthReq) GetMail() string {
	if m != nil {
		return m.Mail
	}
	return ""
}

func (m *UserAuthReq) GetPass() string {
	if m != nil {
		return m.Pass
	}
	return ""
}

type UserShowReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Mail                 string   `protobuf:"bytes,2,opt,name=mail,proto3" json:"mail,omitempty"`
	Active               bool     `protobuf:"varint,3,opt,name=active,proto3" json:"active,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserShowReq) Reset()         { *m = UserShowReq{} }
func (m *UserShowReq) String() string { return proto.CompactTextString(m) }
func (*UserShowReq) ProtoMessage()    {}
func (*UserShowReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{4}
}

func (m *UserShowReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserShowReq.Unmarshal(m, b)
}
func (m *UserShowReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserShowReq.Marshal(b, m, deterministic)
}
func (m *UserShowReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserShowReq.Merge(m, src)
}
func (m *UserShowReq) XXX_Size() int {
	return xxx_messageInfo_UserShowReq.Size(m)
}
func (m *UserShowReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserShowReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserShowReq proto.InternalMessageInfo

func (m *UserShowReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserShowReq) GetMail() string {
	if m != nil {
		return m.Mail
	}
	return ""
}

func (m *UserShowReq) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

type NamedItem struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NamedItem) Reset()         { *m = NamedItem{} }
func (m *NamedItem) String() string { return proto.CompactTextString(m) }
func (*NamedItem) ProtoMessage()    {}
func (*NamedItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{5}
}

func (m *NamedItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NamedItem.Unmarshal(m, b)
}
func (m *NamedItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NamedItem.Marshal(b, m, deterministic)
}
func (m *NamedItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NamedItem.Merge(m, src)
}
func (m *NamedItem) XXX_Size() int {
	return xxx_messageInfo_NamedItem.Size(m)
}
func (m *NamedItem) XXX_DiscardUnknown() {
	xxx_messageInfo_NamedItem.DiscardUnknown(m)
}

var xxx_messageInfo_NamedItem proto.InternalMessageInfo

func (m *NamedItem) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NamedItem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CharacterView struct {
	Id                   uint64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Region               string       `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Name                 string       `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Off                  bool         `protobuf:"varint,4,opt,name=off,proto3" json:"off,omitempty"`
	Cities               []*NamedItem `protobuf:"bytes,5,rep,name=cities,proto3" json:"cities,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CharacterView) Reset()         { *m = CharacterView{} }
func (m *CharacterView) String() string { return proto.CompactTextString(m) }
func (*CharacterView) ProtoMessage()    {}
func (*CharacterView) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{6}
}

func (m *CharacterView) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterView.Unmarshal(m, b)
}
func (m *CharacterView) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterView.Marshal(b, m, deterministic)
}
func (m *CharacterView) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterView.Merge(m, src)
}
func (m *CharacterView) XXX_Size() int {
	return xxx_messageInfo_CharacterView.Size(m)
}
func (m *CharacterView) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterView.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterView proto.InternalMessageInfo

func (m *CharacterView) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CharacterView) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *CharacterView) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CharacterView) GetOff() bool {
	if m != nil {
		return m.Off
	}
	return false
}

func (m *CharacterView) GetCities() []*NamedItem {
	if m != nil {
		return m.Cities
	}
	return nil
}

type UserView struct {
	Id                   uint64           `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Mail                 string           `protobuf:"bytes,2,opt,name=mail,proto3" json:"mail,omitempty"`
	Name                 string           `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Inactive             bool             `protobuf:"varint,4,opt,name=inactive,proto3" json:"inactive,omitempty"`
	Suspended            bool             `protobuf:"varint,5,opt,name=suspended,proto3" json:"suspended,omitempty"`
	Admin                bool             `protobuf:"varint,6,opt,name=admin,proto3" json:"admin,omitempty"`
	Characters           []*CharacterView `protobuf:"bytes,7,rep,name=characters,proto3" json:"characters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *UserView) Reset()         { *m = UserView{} }
func (m *UserView) String() string { return proto.CompactTextString(m) }
func (*UserView) ProtoMessage()    {}
func (*UserView) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{7}
}

func (m *UserView) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserView.Unmarshal(m, b)
}
func (m *UserView) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserView.Marshal(b, m, deterministic)
}
func (m *UserView) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserView.Merge(m, src)
}
func (m *UserView) XXX_Size() int {
	return xxx_messageInfo_UserView.Size(m)
}
func (m *UserView) XXX_DiscardUnknown() {
	xxx_messageInfo_UserView.DiscardUnknown(m)
}

var xxx_messageInfo_UserView proto.InternalMessageInfo

func (m *UserView) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserView) GetMail() string {
	if m != nil {
		return m.Mail
	}
	return ""
}

func (m *UserView) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UserView) GetInactive() bool {
	if m != nil {
		return m.Inactive
	}
	return false
}

func (m *UserView) GetSuspended() bool {
	if m != nil {
		return m.Suspended
	}
	return false
}

func (m *UserView) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

func (m *UserView) GetCharacters() []*CharacterView {
	if m != nil {
		return m.Characters
	}
	return nil
}

type UserSuspendReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserSuspendReq) Reset()         { *m = UserSuspendReq{} }
func (m *UserSuspendReq) String() string { return proto.CompactTextString(m) }
func (*UserSuspendReq) ProtoMessage()    {}
func (*UserSuspendReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{8}
}

func (m *UserSuspendReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserSuspendReq.Unmarshal(m, b)
}
func (m *UserSuspendReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserSuspendReq.Marshal(b, m, deterministic)
}
func (m *UserSuspendReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserSuspendReq.Merge(m, src)
}
func (m *UserSuspendReq) XXX_Size() int {
	return xxx_messageInfo_UserSuspendReq.Size(m)
}
func (m *UserSuspendReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserSuspendReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserSuspendReq proto.InternalMessageInfo

func (m *UserSuspendReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserSuspendReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type UserListReq struct {
	Marker               uint64   `protobuf:"varint,1,opt,name=marker,proto3" json:"marker,omitempty"`
	Limit                uint64   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserListReq) Reset()         { *m = UserListReq{} }
func (m *UserListReq) String() string { return proto.CompactTextString(m) }
func (*UserListReq) ProtoMessage()    {}
func (*UserListReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{9}
}

func (m *UserListReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserListReq.Unmarshal(m, b)
}
func (m *UserListReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserListReq.Marshal(b, m, deterministic)
}
func (m *UserListReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserListReq.Merge(m, src)
}
func (m *UserListReq) XXX_Size() int {
	return xxx_messageInfo_UserListReq.Size(m)
}
func (m *UserListReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserListReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserListReq proto.InternalMessageInfo

func (m *UserListReq) GetMarker() uint64 {
	if m != nil {
		return m.Marker
	}
	return 0
}

func (m *UserListReq) GetLimit() uint64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type UserListRep struct {
	Items                []*UserView `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *UserListRep) Reset()         { *m = UserListRep{} }
func (m *UserListRep) String() string { return proto.CompactTextString(m) }
func (*UserListRep) ProtoMessage()    {}
func (*UserListRep) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{10}
}

func (m *UserListRep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserListRep.Unmarshal(m, b)
}
func (m *UserListRep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserListRep.Marshal(b, m, deterministic)
}
func (m *UserListRep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserListRep.Merge(m, src)
}
func (m *UserListRep) XXX_Size() int {
	return xxx_messageInfo_UserListRep.Size(m)
}
func (m *UserListRep) XXX_DiscardUnknown() {
	xxx_messageInfo_UserListRep.DiscardUnknown(m)
}

var xxx_messageInfo_UserListRep proto.InternalMessageInfo

func (m *UserListRep) GetItems() []*UserView {
	if m != nil {
		return m.Items
	}
	return nil
}

type CharacterShowReq struct {
	User                 uint64   `protobuf:"varint,1,opt,name=user,proto3" json:"user,omitempty"`
	Character            uint64   `protobuf:"varint,2,opt,name=character,proto3" json:"character,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CharacterShowReq) Reset()         { *m = CharacterShowReq{} }
func (m *CharacterShowReq) String() string { return proto.CompactTextString(m) }
func (*CharacterShowReq) ProtoMessage()    {}
func (*CharacterShowReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{11}
}

func (m *CharacterShowReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterShowReq.Unmarshal(m, b)
}
func (m *CharacterShowReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterShowReq.Marshal(b, m, deterministic)
}
func (m *CharacterShowReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterShowReq.Merge(m, src)
}
func (m *CharacterShowReq) XXX_Size() int {
	return xxx_messageInfo_CharacterShowReq.Size(m)
}
func (m *CharacterShowReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterShowReq.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterShowReq proto.InternalMessageInfo

func (m *CharacterShowReq) GetUser() uint64 {
	if m != nil {
		return m.User
	}
	return 0
}

func (m *CharacterShowReq) GetCharacter() uint64 {
	if m != nil {
		return m.Character
	}
	return 0
}

func init() {
	proto.RegisterType((*None)(nil), "hegemonie.auth.proto.None")
	proto.RegisterType((*UserCreateReq)(nil), "hegemonie.auth.proto.UserCreateReq")
	proto.RegisterType((*UserUpdateReq)(nil), "hegemonie.auth.proto.UserUpdateReq")
	proto.RegisterType((*UserAuthReq)(nil), "hegemonie.auth.proto.UserAuthReq")
	proto.RegisterType((*UserShowReq)(nil), "hegemonie.auth.proto.UserShowReq")
	proto.RegisterType((*NamedItem)(nil), "hegemonie.auth.proto.NamedItem")
	proto.RegisterType((*CharacterView)(nil), "hegemonie.auth.proto.CharacterView")
	proto.RegisterType((*UserView)(nil), "hegemonie.auth.proto.UserView")
	proto.RegisterType((*UserSuspendReq)(nil), "hegemonie.auth.proto.UserSuspendReq")
	proto.RegisterType((*UserListReq)(nil), "hegemonie.auth.proto.UserListReq")
	proto.RegisterType((*UserListRep)(nil), "hegemonie.auth.proto.UserListRep")
	proto.RegisterType((*CharacterShowReq)(nil), "hegemonie.auth.proto.CharacterShowReq")
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x1d, 0xc7, 0x24, 0x53, 0xb5, 0xaa, 0x56, 0x55, 0x64, 0x59, 0x08, 0xca, 0x06, 0xa1,
	0x9e, 0x82, 0x54, 0x8a, 0x38, 0x70, 0x42, 0x41, 0x42, 0x95, 0x50, 0x2b, 0x5c, 0x15, 0x89, 0xe3,
	0x12, 0x4f, 0xeb, 0x15, 0x5d, 0xdb, 0x78, 0x37, 0xf4, 0x9f, 0xf0, 0xbb, 0xb8, 0xf1, 0x77, 0xd0,
	0x7e, 0xc4, 0x31, 0xea, 0x3a, 0xc9, 0x6d, 0x67, 0x3c, 0xf3, 0xf6, 0xbd, 0x37, 0xeb, 0x01, 0x60,
	0x4b, 0x55, 0xcc, 0xea, 0xa6, 0x52, 0x15, 0x39, 0x2e, 0xf0, 0x0e, 0x45, 0x55, 0x72, 0x9c, 0xad,
	0xb3, 0x34, 0x86, 0xe8, 0xb2, 0x2a, 0x91, 0x4e, 0xe1, 0xe0, 0x46, 0x62, 0x33, 0x6f, 0x90, 0x29,
	0xcc, 0xf0, 0x27, 0x21, 0x10, 0x09, 0xc6, 0xef, 0x93, 0xe0, 0x24, 0x38, 0x1d, 0x67, 0xe6, 0x4c,
	0x3f, 0xd9, 0xa2, 0x9b, 0x3a, 0x77, 0x45, 0x87, 0x10, 0xf2, 0xdc, 0x94, 0x44, 0x59, 0xc8, 0x73,
	0xdd, 0x54, 0x33, 0x29, 0x93, 0xd0, 0x36, 0xe9, 0xb3, 0xce, 0x95, 0x4c, 0x60, 0x32, 0xb0, 0x39,
	0x7d, 0xa6, 0x6f, 0x61, 0x5f, 0x03, 0x7d, 0x58, 0xaa, 0xa2, 0xe7, 0x2e, 0x1f, 0x14, 0xbd, 0xb0,
	0x6d, 0xd7, 0x45, 0xf5, 0xd0, 0x73, 0xbb, 0x81, 0x09, 0x3b, 0x30, 0x13, 0x88, 0xd9, 0x42, 0xf1,
	0x5f, 0xf6, 0xfe, 0x51, 0xe6, 0x22, 0xfa, 0x1a, 0xc6, 0x97, 0x4c, 0x60, 0x7e, 0xa1, 0x50, 0xf8,
	0x80, 0x0c, 0xe5, 0xb0, 0x43, 0xf9, 0x77, 0x00, 0x07, 0xf3, 0x82, 0x35, 0x6c, 0xa1, 0xb0, 0xf9,
	0xca, 0xf1, 0xe1, 0x51, 0xd7, 0x04, 0xe2, 0x06, 0xef, 0x78, 0x55, 0xba, 0x3e, 0x17, 0xf9, 0x0c,
	0x20, 0x47, 0x30, 0xa8, 0x6e, 0x6f, 0x93, 0xc8, 0x70, 0xd2, 0x47, 0xf2, 0x0e, 0xe2, 0x05, 0x57,
	0x1c, 0x65, 0x32, 0x3c, 0x19, 0x9c, 0xee, 0x9f, 0x3d, 0x9f, 0xf9, 0xe6, 0x35, 0x6b, 0x49, 0x67,
	0xae, 0x9c, 0xfe, 0x0d, 0x60, 0xa4, 0x5d, 0xf1, 0x72, 0xf2, 0x59, 0xe2, 0xe3, 0x93, 0xc2, 0x88,
	0x97, 0xce, 0x28, 0x4b, 0xaa, 0x8d, 0xc9, 0x53, 0x18, 0xcb, 0xa5, 0xac, 0xb1, 0xcc, 0x31, 0x4f,
	0x86, 0xe6, 0xe3, 0x3a, 0x41, 0x8e, 0x61, 0xc8, 0x72, 0xc1, 0xcb, 0x24, 0x36, 0x5f, 0x6c, 0x40,
	0xe6, 0x00, 0x8b, 0x95, 0x59, 0x32, 0x79, 0x62, 0x14, 0x4d, 0xfd, 0x8a, 0xfe, 0x33, 0x35, 0xeb,
	0xb4, 0xd1, 0x73, 0x38, 0x34, 0xe3, 0xb6, 0x77, 0xf5, 0x4c, 0xfc, 0xd1, 0xa0, 0xde, 0xdb, 0x47,
	0xf2, 0x99, 0x4b, 0xa5, 0x5b, 0x26, 0x10, 0x0b, 0xd6, 0xfc, 0xc0, 0xc6, 0xb5, 0xb9, 0x48, 0xf3,
	0xbe, 0xe7, 0x82, 0x2b, 0xd3, 0x1b, 0x65, 0x36, 0xa0, 0xf3, 0x6e, 0x73, 0x4d, 0xce, 0x61, 0xc8,
	0x15, 0x0a, 0x99, 0x04, 0x46, 0xc1, 0x33, 0xbf, 0x82, 0x95, 0xfb, 0x99, 0x2d, 0xa6, 0x1f, 0xe1,
	0xa8, 0x15, 0xb5, 0x7a, 0xab, 0x04, 0xa2, 0xa5, 0x6c, 0x49, 0x98, 0xb3, 0x36, 0xb6, 0x55, 0xeb,
	0x68, 0xac, 0x13, 0x67, 0x7f, 0x22, 0x88, 0xf4, 0x0f, 0x42, 0x32, 0x3b, 0x5f, 0xcd, 0x89, 0xbc,
	0xe8, 0x67, 0xe0, 0x04, 0xa7, 0x5b, 0x4b, 0x6a, 0xba, 0x47, 0xae, 0x2c, 0xa6, 0x66, 0xb7, 0x09,
	0xd3, 0xb1, 0x4f, 0xb7, 0x08, 0xa7, 0x7b, 0xe4, 0x1a, 0x60, 0xbd, 0x3f, 0xc8, 0xb4, 0xbf, 0xbe,
	0xdd, 0x30, 0x3b, 0x80, 0x5e, 0x59, 0x50, 0xbb, 0x6f, 0x36, 0x81, 0xb6, 0x1b, 0x29, 0x4d, 0x7b,
	0x7e, 0x1b, 0xbd, 0xe3, 0xf6, 0xc8, 0x17, 0xb7, 0x40, 0xec, 0x8b, 0x22, 0x2f, 0x37, 0x28, 0x6f,
	0x1f, 0xdd, 0x16, 0x48, 0xe7, 0xa4, 0x99, 0xd4, 0x06, 0x27, 0xdd, 0xaa, 0xdb, 0x41, 0xf4, 0xb7,
	0xce, 0x9e, 0x31, 0xf3, 0x79, 0xb5, 0xe5, 0xbf, 0xd9, 0x79, 0x48, 0xdf, 0x63, 0x93, 0x79, 0xf3,
	0x2f, 0x00, 0x00, 0xff, 0xff, 0x9a, 0x6d, 0x06, 0x00, 0x17, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthClient interface {
	UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListRep, error)
	UserShow(ctx context.Context, in *UserShowReq, opts ...grpc.CallOption) (*UserView, error)
	UserCreate(ctx context.Context, in *UserCreateReq, opts ...grpc.CallOption) (*UserView, error)
	UserUpdate(ctx context.Context, in *UserUpdateReq, opts ...grpc.CallOption) (*None, error)
	UserSuspend(ctx context.Context, in *UserSuspendReq, opts ...grpc.CallOption) (*None, error)
	UserAuth(ctx context.Context, in *UserAuthReq, opts ...grpc.CallOption) (*UserView, error)
	// Check the given Character can be managed by the given User and return
	// an abstract of the Character information. The user information is also
	// returned to save calls from the main service.
	CharacterShow(ctx context.Context, in *CharacterShowReq, opts ...grpc.CallOption) (*UserView, error)
}

type authClient struct {
	cc *grpc.ClientConn
}

func NewAuthClient(cc *grpc.ClientConn) AuthClient {
	return &authClient{cc}
}

func (c *authClient) UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListRep, error) {
	out := new(UserListRep)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UserShow(ctx context.Context, in *UserShowReq, opts ...grpc.CallOption) (*UserView, error) {
	out := new(UserView)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserShow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UserCreate(ctx context.Context, in *UserCreateReq, opts ...grpc.CallOption) (*UserView, error) {
	out := new(UserView)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UserUpdate(ctx context.Context, in *UserUpdateReq, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UserSuspend(ctx context.Context, in *UserSuspendReq, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserSuspend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UserAuth(ctx context.Context, in *UserAuthReq, opts ...grpc.CallOption) (*UserView, error) {
	out := new(UserView)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/UserAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) CharacterShow(ctx context.Context, in *CharacterShowReq, opts ...grpc.CallOption) (*UserView, error) {
	out := new(UserView)
	err := c.cc.Invoke(ctx, "/hegemonie.auth.proto.Auth/CharacterShow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
type AuthServer interface {
	UserList(context.Context, *UserListReq) (*UserListRep, error)
	UserShow(context.Context, *UserShowReq) (*UserView, error)
	UserCreate(context.Context, *UserCreateReq) (*UserView, error)
	UserUpdate(context.Context, *UserUpdateReq) (*None, error)
	UserSuspend(context.Context, *UserSuspendReq) (*None, error)
	UserAuth(context.Context, *UserAuthReq) (*UserView, error)
	// Check the given Character can be managed by the given User and return
	// an abstract of the Character information. The user information is also
	// returned to save calls from the main service.
	CharacterShow(context.Context, *CharacterShowReq) (*UserView, error)
}

// UnimplementedAuthServer can be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (*UnimplementedAuthServer) UserList(ctx context.Context, req *UserListReq) (*UserListRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserList not implemented")
}
func (*UnimplementedAuthServer) UserShow(ctx context.Context, req *UserShowReq) (*UserView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserShow not implemented")
}
func (*UnimplementedAuthServer) UserCreate(ctx context.Context, req *UserCreateReq) (*UserView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreate not implemented")
}
func (*UnimplementedAuthServer) UserUpdate(ctx context.Context, req *UserUpdateReq) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUpdate not implemented")
}
func (*UnimplementedAuthServer) UserSuspend(ctx context.Context, req *UserSuspendReq) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserSuspend not implemented")
}
func (*UnimplementedAuthServer) UserAuth(ctx context.Context, req *UserAuthReq) (*UserView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserAuth not implemented")
}
func (*UnimplementedAuthServer) CharacterShow(ctx context.Context, req *CharacterShowReq) (*UserView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CharacterShow not implemented")
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_UserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserList(ctx, req.(*UserListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UserShow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserShowReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserShow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserShow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserShow(ctx, req.(*UserShowReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UserCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserCreate(ctx, req.(*UserCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UserUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserUpdate(ctx, req.(*UserUpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UserSuspend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserSuspendReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserSuspend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserSuspend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserSuspend(ctx, req.(*UserSuspendReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UserAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserAuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UserAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/UserAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UserAuth(ctx, req.(*UserAuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_CharacterShow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CharacterShowReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).CharacterShow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hegemonie.auth.proto.Auth/CharacterShow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).CharacterShow(ctx, req.(*CharacterShowReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hegemonie.auth.proto.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserList",
			Handler:    _Auth_UserList_Handler,
		},
		{
			MethodName: "UserShow",
			Handler:    _Auth_UserShow_Handler,
		},
		{
			MethodName: "UserCreate",
			Handler:    _Auth_UserCreate_Handler,
		},
		{
			MethodName: "UserUpdate",
			Handler:    _Auth_UserUpdate_Handler,
		},
		{
			MethodName: "UserSuspend",
			Handler:    _Auth_UserSuspend_Handler,
		},
		{
			MethodName: "UserAuth",
			Handler:    _Auth_UserAuth_Handler,
		},
		{
			MethodName: "CharacterShow",
			Handler:    _Auth_CharacterShow_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}