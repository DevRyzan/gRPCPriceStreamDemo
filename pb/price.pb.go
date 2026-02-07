package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetCurrentPriceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetCurrentPriceRequest) Reset() {
	*x = GetCurrentPriceRequest{}
	mi := &file_proto_price_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCurrentPriceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCurrentPriceRequest) ProtoMessage() {}

func (x *GetCurrentPriceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_price_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetCurrentPriceRequest) Descriptor() ([]byte, []int) {
	return file_proto_price_proto_rawDescGZIP(), []int{0}
}

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	mi := &file_proto_price_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_price_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *SubscribeRequest) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

type PriceUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string  `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Price  float64 `protobuf:"fixed64,2,opt,name=price,proto3" json:"price,omitempty"`
	AtTs   int64   `protobuf:"varint,3,opt,name=at_ts,json=atTs,proto3" json:"at_ts,omitempty"`
}

func (x *PriceUpdate) Reset() {
	*x = PriceUpdate{}
	mi := &file_proto_price_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PriceUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PriceUpdate) ProtoMessage() {}

func (x *PriceUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_price_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PriceUpdate) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *PriceUpdate) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *PriceUpdate) GetAtTs() int64 {
	if x != nil {
		return x.AtTs
	}
	return 0
}

var File_proto_price_proto protoreflect.FileDescriptor

var file_proto_price_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x2a, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x22, 0x4c, 0x0a, 0x0b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x13, 0x0a,
	0x05, 0x61, 0x74, 0x5f, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x61, 0x74,
	0x54, 0x73, 0x32, 0xa3, 0x01, 0x0a, 0x0c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x69, 0x63, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x14, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_price_proto_rawDescOnce sync.Once
	file_proto_price_proto_rawDescData = file_proto_price_proto_rawDesc
)

func file_proto_price_proto_rawDescGZIP() []byte {
	file_proto_price_proto_rawDescOnce.Do(func() {
		file_proto_price_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_price_proto_rawDescData)
	})
	return file_proto_price_proto_rawDescData
}

var file_proto_price_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_price_proto_goTypes = []any{
	(*GetCurrentPriceRequest)(nil), // 0: price.GetCurrentPriceRequest
	(*SubscribeRequest)(nil),       // 1: price.SubscribeRequest
	(*PriceUpdate)(nil),            // 2: price.PriceUpdate
}
var file_proto_price_proto_depIdxs = []int32{
	0, // 0: price.PriceService.GetCurrentPrice:input_type -> price.GetCurrentPriceRequest
	1, // 1: price.PriceService.SubscribePriceUpdates:input_type -> price.SubscribeRequest
	2, // 2: price.PriceService.GetCurrentPrice:output_type -> price.PriceUpdate
	2, // 3: price.PriceService.SubscribePriceUpdates:output_type -> price.PriceUpdate
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_price_proto_init() }
func file_proto_price_proto_init() {
	if File_proto_price_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_price_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetCurrentPriceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_price_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_price_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PriceUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_price_proto_rawDesc,
			NumEnums:      0,
			NumMessages:  3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_price_proto_goTypes,
		DependencyIndexes: file_proto_price_proto_depIdxs,
		MessageInfos:      file_proto_price_proto_msgTypes,
	}.Build()
	File_proto_price_proto = out.File
	file_proto_price_proto_rawDesc = nil
	file_proto_price_proto_goTypes = nil
	file_proto_price_proto_depIdxs = nil
}
