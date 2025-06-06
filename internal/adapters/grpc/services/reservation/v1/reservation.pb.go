// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: reservation/v1/reservation.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BaseReservationParams struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BaseReservationParams) Reset() {
	*x = BaseReservationParams{}
	mi := &file_reservation_v1_reservation_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BaseReservationParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseReservationParams) ProtoMessage() {}

func (x *BaseReservationParams) ProtoReflect() protoreflect.Message {
	mi := &file_reservation_v1_reservation_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseReservationParams.ProtoReflect.Descriptor instead.
func (*BaseReservationParams) Descriptor() ([]byte, []int) {
	return file_reservation_v1_reservation_proto_rawDescGZIP(), []int{0}
}

func (x *BaseReservationParams) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CommitRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Params        *BaseReservationParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	mi := &file_reservation_v1_reservation_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reservation_v1_reservation_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_reservation_v1_reservation_proto_rawDescGZIP(), []int{1}
}

func (x *CommitRequest) GetParams() *BaseReservationParams {
	if x != nil {
		return x.Params
	}
	return nil
}

type RollbackRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Params        *BaseReservationParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RollbackRequest) Reset() {
	*x = RollbackRequest{}
	mi := &file_reservation_v1_reservation_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RollbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RollbackRequest) ProtoMessage() {}

func (x *RollbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reservation_v1_reservation_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RollbackRequest.ProtoReflect.Descriptor instead.
func (*RollbackRequest) Descriptor() ([]byte, []int) {
	return file_reservation_v1_reservation_proto_rawDescGZIP(), []int{2}
}

func (x *RollbackRequest) GetParams() *BaseReservationParams {
	if x != nil {
		return x.Params
	}
	return nil
}

var File_reservation_v1_reservation_proto protoreflect.FileDescriptor

const file_reservation_v1_reservation_proto_rawDesc = "" +
	"\n" +
	" reservation/v1/reservation.proto\x12\x0ereservation.v1\x1a\x1bgoogle/protobuf/empty.proto\"'\n" +
	"\x15BaseReservationParams\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"N\n" +
	"\rCommitRequest\x12=\n" +
	"\x06params\x18\x01 \x01(\v2%.reservation.v1.BaseReservationParamsR\x06params\"P\n" +
	"\x0fRollbackRequest\x12=\n" +
	"\x06params\x18\x01 \x01(\v2%.reservation.v1.BaseReservationParamsR\x06params2\x9a\x01\n" +
	"\x12ReservationService\x12?\n" +
	"\x06Commit\x12\x1d.reservation.v1.CommitRequest\x1a\x16.google.protobuf.Empty\x12C\n" +
	"\bRollback\x12\x1f.reservation.v1.RollbackRequest\x1a\x16.google.protobuf.EmptyB\x10Z\x0ereservation/pbb\x06proto3"

var (
	file_reservation_v1_reservation_proto_rawDescOnce sync.Once
	file_reservation_v1_reservation_proto_rawDescData []byte
)

func file_reservation_v1_reservation_proto_rawDescGZIP() []byte {
	file_reservation_v1_reservation_proto_rawDescOnce.Do(func() {
		file_reservation_v1_reservation_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_reservation_v1_reservation_proto_rawDesc), len(file_reservation_v1_reservation_proto_rawDesc)))
	})
	return file_reservation_v1_reservation_proto_rawDescData
}

var file_reservation_v1_reservation_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_reservation_v1_reservation_proto_goTypes = []any{
	(*BaseReservationParams)(nil), // 0: reservation.v1.BaseReservationParams
	(*CommitRequest)(nil),         // 1: reservation.v1.CommitRequest
	(*RollbackRequest)(nil),       // 2: reservation.v1.RollbackRequest
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_reservation_v1_reservation_proto_depIdxs = []int32{
	0, // 0: reservation.v1.CommitRequest.params:type_name -> reservation.v1.BaseReservationParams
	0, // 1: reservation.v1.RollbackRequest.params:type_name -> reservation.v1.BaseReservationParams
	1, // 2: reservation.v1.ReservationService.Commit:input_type -> reservation.v1.CommitRequest
	2, // 3: reservation.v1.ReservationService.Rollback:input_type -> reservation.v1.RollbackRequest
	3, // 4: reservation.v1.ReservationService.Commit:output_type -> google.protobuf.Empty
	3, // 5: reservation.v1.ReservationService.Rollback:output_type -> google.protobuf.Empty
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_reservation_v1_reservation_proto_init() }
func file_reservation_v1_reservation_proto_init() {
	if File_reservation_v1_reservation_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_reservation_v1_reservation_proto_rawDesc), len(file_reservation_v1_reservation_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_reservation_v1_reservation_proto_goTypes,
		DependencyIndexes: file_reservation_v1_reservation_proto_depIdxs,
		MessageInfos:      file_reservation_v1_reservation_proto_msgTypes,
	}.Build()
	File_reservation_v1_reservation_proto = out.File
	file_reservation_v1_reservation_proto_goTypes = nil
	file_reservation_v1_reservation_proto_depIdxs = nil
}
