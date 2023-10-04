// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: network/networkservice.proto

// Network gRPC service

package networkservice

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NicSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vendors []string `protobuf:"bytes,1,rep,name=vendors,proto3" json:"vendors,omitempty"`
	Drivers []string `protobuf:"bytes,2,rep,name=drivers,proto3" json:"drivers,omitempty"`
	Devices []string `protobuf:"bytes,3,rep,name=devices,proto3" json:"devices,omitempty"`
	PfNames []string `protobuf:"bytes,4,rep,name=pfNames,proto3" json:"pfNames,omitempty"`
}

func (x *NicSelector) Reset() {
	*x = NicSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NicSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NicSelector) ProtoMessage() {}

func (x *NicSelector) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NicSelector.ProtoReflect.Descriptor instead.
func (*NicSelector) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{0}
}

func (x *NicSelector) GetVendors() []string {
	if x != nil {
		return x.Vendors
	}
	return nil
}

func (x *NicSelector) GetDrivers() []string {
	if x != nil {
		return x.Drivers
	}
	return nil
}

func (x *NicSelector) GetDevices() []string {
	if x != nil {
		return x.Devices
	}
	return nil
}

func (x *NicSelector) GetPfNames() []string {
	if x != nil {
		return x.PfNames
	}
	return nil
}

type ResourceConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Mtu         uint32       `protobuf:"varint,2,opt,name=mtu,proto3" json:"mtu,omitempty"`
	NumVfs      uint32       `protobuf:"varint,3,opt,name=numVfs,proto3" json:"numVfs,omitempty"`
	NicSelector *NicSelector `protobuf:"bytes,4,opt,name=nicSelector,proto3" json:"nicSelector,omitempty"`
	DeviceType  string       `protobuf:"bytes,5,opt,name=deviceType,proto3" json:"deviceType,omitempty"`
}

func (x *ResourceConfig) Reset() {
	*x = ResourceConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceConfig) ProtoMessage() {}

func (x *ResourceConfig) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceConfig.ProtoReflect.Descriptor instead.
func (*ResourceConfig) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{1}
}

func (x *ResourceConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ResourceConfig) GetMtu() uint32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *ResourceConfig) GetNumVfs() uint32 {
	if x != nil {
		return x.NumVfs
	}
	return 0
}

func (x *ResourceConfig) GetNicSelector() *NicSelector {
	if x != nil {
		return x.NicSelector
	}
	return nil
}

func (x *ResourceConfig) GetDeviceType() string {
	if x != nil {
		return x.DeviceType
	}
	return ""
}

type VFResourceStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Vendor string `protobuf:"bytes,2,opt,name=vendor,proto3" json:"vendor,omitempty"`
	Driver string `protobuf:"bytes,3,opt,name=driver,proto3" json:"driver,omitempty"`
	Device string `protobuf:"bytes,4,opt,name=device,proto3" json:"device,omitempty"`
}

func (x *VFResourceStatus) Reset() {
	*x = VFResourceStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VFResourceStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VFResourceStatus) ProtoMessage() {}

func (x *VFResourceStatus) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VFResourceStatus.ProtoReflect.Descriptor instead.
func (*VFResourceStatus) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{2}
}

func (x *VFResourceStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *VFResourceStatus) GetVendor() string {
	if x != nil {
		return x.Vendor
	}
	return ""
}

func (x *VFResourceStatus) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *VFResourceStatus) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

type ResourceStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Mtu    uint32              `protobuf:"varint,2,opt,name=mtu,proto3" json:"mtu,omitempty"`
	NumVfs uint32              `protobuf:"varint,3,opt,name=numVfs,proto3" json:"numVfs,omitempty"`
	Mac    string              `protobuf:"bytes,4,opt,name=mac,proto3" json:"mac,omitempty"`
	Vendor string              `protobuf:"bytes,5,opt,name=vendor,proto3" json:"vendor,omitempty"`
	Driver string              `protobuf:"bytes,6,opt,name=driver,proto3" json:"driver,omitempty"`
	Device string              `protobuf:"bytes,7,opt,name=device,proto3" json:"device,omitempty"`
	Vfs    []*VFResourceStatus `protobuf:"bytes,8,rep,name=vfs,proto3" json:"vfs,omitempty"`
}

func (x *ResourceStatus) Reset() {
	*x = ResourceStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceStatus) ProtoMessage() {}

func (x *ResourceStatus) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceStatus.ProtoReflect.Descriptor instead.
func (*ResourceStatus) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{3}
}

func (x *ResourceStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ResourceStatus) GetMtu() uint32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *ResourceStatus) GetNumVfs() uint32 {
	if x != nil {
		return x.NumVfs
	}
	return 0
}

func (x *ResourceStatus) GetMac() string {
	if x != nil {
		return x.Mac
	}
	return ""
}

func (x *ResourceStatus) GetVendor() string {
	if x != nil {
		return x.Vendor
	}
	return ""
}

func (x *ResourceStatus) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *ResourceStatus) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

func (x *ResourceStatus) GetVfs() []*VFResourceStatus {
	if x != nil {
		return x.Vfs
	}
	return nil
}

type ResourceSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Mtu     uint32   `protobuf:"varint,2,opt,name=mtu,proto3" json:"mtu,omitempty"`
	NumVfs  uint32   `protobuf:"varint,3,opt,name=numVfs,proto3" json:"numVfs,omitempty"`
	Devices []string `protobuf:"bytes,4,rep,name=devices,proto3" json:"devices,omitempty"`
}

func (x *ResourceSpec) Reset() {
	*x = ResourceSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceSpec) ProtoMessage() {}

func (x *ResourceSpec) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceSpec.ProtoReflect.Descriptor instead.
func (*ResourceSpec) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{4}
}

func (x *ResourceSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ResourceSpec) GetMtu() uint32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *ResourceSpec) GetNumVfs() uint32 {
	if x != nil {
		return x.NumVfs
	}
	return 0
}

func (x *ResourceSpec) GetDevices() []string {
	if x != nil {
		return x.Devices
	}
	return nil
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Spec   *ResourceSpec     `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
	Status []*ResourceStatus `protobuf:"bytes,2,rep,name=status,proto3" json:"status,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_network_networkservice_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_network_networkservice_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_network_networkservice_proto_rawDescGZIP(), []int{5}
}

func (x *Resource) GetSpec() *ResourceSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *Resource) GetStatus() []*ResourceStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_network_networkservice_proto protoreflect.FileDescriptor

var file_network_networkservice_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x75, 0x0a, 0x0b, 0x4e, 0x69, 0x63,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x6e, 0x64,
	0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x6e, 0x64, 0x6f,
	0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x07, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x66, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x70, 0x66, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x22, 0xad, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6d, 0x74, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x56, 0x66, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x56, 0x66,
	0x73, 0x12, 0x3d, 0x0a, 0x0b, 0x6e, 0x69, 0x63, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4e, 0x69, 0x63, 0x53, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x52, 0x0b, 0x6e, 0x69, 0x63, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x6e, 0x0a, 0x10, 0x56, 0x46, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x6e, 0x64,
	0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x22, 0xdc, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6d, 0x74, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x56, 0x66, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x56, 0x66,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x61, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x72, 0x69,
	0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x03, 0x76,
	0x66, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x46, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x03, 0x76, 0x66, 0x73, 0x22,
	0x66, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x6d, 0x74, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x56, 0x66, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x56, 0x66, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x74, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x65, 0x63, 0x52,
	0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x36, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0xe0, 0x02,
	0x0a, 0x0e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x73, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x1e, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12, 0x18, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x30, 0x01, 0x12, 0x77, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1e, 0x2e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x18, 0x2e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x3a,
	0x01, 0x2a, 0x22, 0x18, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x30, 0x01, 0x12, 0x60,
	0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x30, 0x01,
	0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x72, 0x69, 0x72, 0x61, 0x6d, 0x79, 0x2f, 0x76, 0x66, 0x2d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x3b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_network_networkservice_proto_rawDescOnce sync.Once
	file_network_networkservice_proto_rawDescData = file_network_networkservice_proto_rawDesc
)

func file_network_networkservice_proto_rawDescGZIP() []byte {
	file_network_networkservice_proto_rawDescOnce.Do(func() {
		file_network_networkservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_network_networkservice_proto_rawDescData)
	})
	return file_network_networkservice_proto_rawDescData
}

var file_network_networkservice_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_network_networkservice_proto_goTypes = []interface{}{
	(*NicSelector)(nil),      // 0: networkservice.NicSelector
	(*ResourceConfig)(nil),   // 1: networkservice.ResourceConfig
	(*VFResourceStatus)(nil), // 2: networkservice.VFResourceStatus
	(*ResourceStatus)(nil),   // 3: networkservice.ResourceStatus
	(*ResourceSpec)(nil),     // 4: networkservice.ResourceSpec
	(*Resource)(nil),         // 5: networkservice.Resource
	(*empty.Empty)(nil),      // 6: google.protobuf.Empty
}
var file_network_networkservice_proto_depIdxs = []int32{
	0, // 0: networkservice.ResourceConfig.nicSelector:type_name -> networkservice.NicSelector
	2, // 1: networkservice.ResourceStatus.vfs:type_name -> networkservice.VFResourceStatus
	4, // 2: networkservice.Resource.spec:type_name -> networkservice.ResourceSpec
	3, // 3: networkservice.Resource.status:type_name -> networkservice.ResourceStatus
	6, // 4: networkservice.NetworkService.GetAllResourceConfigs:input_type -> google.protobuf.Empty
	1, // 5: networkservice.NetworkService.CreateResourceConfig:input_type -> networkservice.ResourceConfig
	6, // 6: networkservice.NetworkService.GetAllResources:input_type -> google.protobuf.Empty
	1, // 7: networkservice.NetworkService.GetAllResourceConfigs:output_type -> networkservice.ResourceConfig
	5, // 8: networkservice.NetworkService.CreateResourceConfig:output_type -> networkservice.Resource
	5, // 9: networkservice.NetworkService.GetAllResources:output_type -> networkservice.Resource
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_network_networkservice_proto_init() }
func file_network_networkservice_proto_init() {
	if File_network_networkservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_network_networkservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NicSelector); i {
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
		file_network_networkservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceConfig); i {
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
		file_network_networkservice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VFResourceStatus); i {
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
		file_network_networkservice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceStatus); i {
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
		file_network_networkservice_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceSpec); i {
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
		file_network_networkservice_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
			RawDescriptor: file_network_networkservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_network_networkservice_proto_goTypes,
		DependencyIndexes: file_network_networkservice_proto_depIdxs,
		MessageInfos:      file_network_networkservice_proto_msgTypes,
	}.Build()
	File_network_networkservice_proto = out.File
	file_network_networkservice_proto_rawDesc = nil
	file_network_networkservice_proto_goTypes = nil
	file_network_networkservice_proto_depIdxs = nil
}
