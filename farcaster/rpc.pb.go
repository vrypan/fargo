// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: rpc.proto

package farcaster

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x68, 0x75, 0x62, 0x5f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72,
	0x6f, 0x6f, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x6f, 0x6e, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xdc,
	0x14, 0x0a, 0x0a, 0x48, 0x75, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x08,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x30, 0x0a, 0x0f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x13, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x12, 0x11, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x48, 0x75, 0x62, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x30,
	0x01, 0x12, 0x24, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0d, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x48,
	0x75, 0x62, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x43, 0x61,
	0x73, 0x74, 0x12, 0x07, 0x2e, 0x43, 0x61, 0x73, 0x74, 0x49, 0x64, 0x1a, 0x08, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43, 0x61, 0x73, 0x74,
	0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x61, 0x73,
	0x74, 0x73, 0x42, 0x79, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x15, 0x2e, 0x43, 0x61, 0x73,
	0x74, 0x73, 0x42, 0x79, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x61, 0x73, 0x74, 0x73,
	0x42, 0x79, 0x4d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x2e, 0x52, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x3e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x16, 0x2e, 0x52, 0x65, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x43, 0x61, 0x73, 0x74, 0x12, 0x19, 0x2e, 0x52, 0x65, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x12, 0x19, 0x2e, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x0b, 0x2e,
	0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x12, 0x15, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x3e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x73, 0x42, 0x79, 0x46,
	0x69, 0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x2e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42,
	0x79, 0x46, 0x69, 0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x12, 0x0e, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x4f, 0x6e,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x73, 0x42, 0x79, 0x46, 0x69,
	0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x2e, 0x4f, 0x6e, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x49, 0x64, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0d, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x55, 0x0a, 0x22, 0x47, 0x65, 0x74, 0x49, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x20, 0x2e, 0x49, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69,
	0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x43, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x64, 0x73, 0x12, 0x0c, 0x2e, 0x46, 0x69, 0x64, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x46, 0x69, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x0c,
	0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e,
	0x6b, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x12, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x42,
	0x79, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x42, 0x79, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x12, 0x15, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x42, 0x79, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x17,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x43, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x14, 0x2e, 0x46, 0x69, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12,
	0x14, 0x2e, 0x46, 0x69, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x41,
	0x6c, 0x6c, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x14, 0x2e, 0x46, 0x69,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x73,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x42, 0x79,
	0x46, 0x69, 0x64, 0x12, 0x14, 0x2e, 0x46, 0x69, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x17,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x4c, 0x69, 0x6e, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x42, 0x79, 0x46, 0x69, 0x64, 0x12, 0x14, 0x2e, 0x46, 0x69, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x41, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x70, 0x61,
	0x63, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x79,
	0x46, 0x69, 0x64, 0x12, 0x0b, 0x2e, 0x46, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x42, 0x75, 0x6c,
	0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x42, 0x75, 0x6c, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x42, 0x75,
	0x6c, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0f, 0x2e,
	0x48, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10,
	0x2e, 0x48, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x63, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x27, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x06, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x09, 0x46, 0x6f,
	0x72, 0x63, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x12, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x79,
	0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x38, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x12, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x64, 0x73, 0x42, 0x79, 0x50, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x12, 0x0f, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x1a, 0x08, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x64, 0x73, 0x12, 0x36,
	0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x42, 0x79, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x64, 0x73, 0x12, 0x08, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x49, 0x64, 0x73, 0x1a, 0x11, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53, 0x79, 0x6e,
	0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x42, 0x79, 0x50, 0x72, 0x65, 0x66, 0x69,
	0x78, 0x12, 0x0f, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x65, 0x66,
	0x69, 0x78, 0x1a, 0x19, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a,
	0x17, 0x47, 0x65, 0x74, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x42, 0x79, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x0f, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e,
	0x6f, 0x64, 0x65, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x1a, 0x19, 0x2e, 0x54, 0x72, 0x69, 0x65,
	0x4e, 0x6f, 0x64, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x79,
	0x6e, 0x63, 0x12, 0x12, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x79, 0x6e, 0x63, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53,
	0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12,
	0x3c, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x13,
	0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x32, 0x90, 0x01,
	0x0a, 0x0c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x21,
	0x0a, 0x0f, 0x52, 0x65, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x79, 0x6e, 0x63, 0x54, 0x72, 0x69,
	0x65, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x29, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x46, 0x72, 0x6f, 0x6d, 0x44, 0x62, 0x12, 0x06, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x12,
	0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x0d, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x1a, 0x0d, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_rpc_proto_goTypes = []any{
	(*Message)(nil),                         // 0: Message
	(*SubscribeRequest)(nil),                // 1: SubscribeRequest
	(*EventRequest)(nil),                    // 2: EventRequest
	(*CastId)(nil),                          // 3: CastId
	(*FidRequest)(nil),                      // 4: FidRequest
	(*CastsByParentRequest)(nil),            // 5: CastsByParentRequest
	(*ReactionRequest)(nil),                 // 6: ReactionRequest
	(*ReactionsByFidRequest)(nil),           // 7: ReactionsByFidRequest
	(*ReactionsByTargetRequest)(nil),        // 8: ReactionsByTargetRequest
	(*UserDataRequest)(nil),                 // 9: UserDataRequest
	(*UsernameProofRequest)(nil),            // 10: UsernameProofRequest
	(*VerificationRequest)(nil),             // 11: VerificationRequest
	(*SignerRequest)(nil),                   // 12: SignerRequest
	(*OnChainEventRequest)(nil),             // 13: OnChainEventRequest
	(*IdRegistryEventByAddressRequest)(nil), // 14: IdRegistryEventByAddressRequest
	(*FidsRequest)(nil),                     // 15: FidsRequest
	(*LinkRequest)(nil),                     // 16: LinkRequest
	(*LinksByFidRequest)(nil),               // 17: LinksByFidRequest
	(*LinksByTargetRequest)(nil),            // 18: LinksByTargetRequest
	(*FidTimestampRequest)(nil),             // 19: FidTimestampRequest
	(*SubmitBulkMessagesRequest)(nil),       // 20: SubmitBulkMessagesRequest
	(*HubInfoRequest)(nil),                  // 21: HubInfoRequest
	(*Empty)(nil),                           // 22: Empty
	(*SyncStatusRequest)(nil),               // 23: SyncStatusRequest
	(*TrieNodePrefix)(nil),                  // 24: TrieNodePrefix
	(*SyncIds)(nil),                         // 25: SyncIds
	(*StreamSyncRequest)(nil),               // 26: StreamSyncRequest
	(*StreamFetchRequest)(nil),              // 27: StreamFetchRequest
	(*OnChainEvent)(nil),                    // 28: OnChainEvent
	(*ValidationResponse)(nil),              // 29: ValidationResponse
	(*HubEvent)(nil),                        // 30: HubEvent
	(*MessagesResponse)(nil),                // 31: MessagesResponse
	(*UserNameProof)(nil),                   // 32: UserNameProof
	(*UsernameProofsResponse)(nil),          // 33: UsernameProofsResponse
	(*OnChainEventResponse)(nil),            // 34: OnChainEventResponse
	(*StorageLimitsResponse)(nil),           // 35: StorageLimitsResponse
	(*FidsResponse)(nil),                    // 36: FidsResponse
	(*SubmitBulkMessagesResponse)(nil),      // 37: SubmitBulkMessagesResponse
	(*HubInfoResponse)(nil),                 // 38: HubInfoResponse
	(*ContactInfoResponse)(nil),             // 39: ContactInfoResponse
	(*SyncStatusResponse)(nil),              // 40: SyncStatusResponse
	(*TrieNodeMetadataResponse)(nil),        // 41: TrieNodeMetadataResponse
	(*TrieNodeSnapshotResponse)(nil),        // 42: TrieNodeSnapshotResponse
	(*StreamSyncResponse)(nil),              // 43: StreamSyncResponse
	(*StreamFetchResponse)(nil),             // 44: StreamFetchResponse
}
var file_rpc_proto_depIdxs = []int32{
	0,  // 0: HubService.SubmitMessage:input_type -> Message
	0,  // 1: HubService.ValidateMessage:input_type -> Message
	1,  // 2: HubService.Subscribe:input_type -> SubscribeRequest
	2,  // 3: HubService.GetEvent:input_type -> EventRequest
	3,  // 4: HubService.GetCast:input_type -> CastId
	4,  // 5: HubService.GetCastsByFid:input_type -> FidRequest
	5,  // 6: HubService.GetCastsByParent:input_type -> CastsByParentRequest
	4,  // 7: HubService.GetCastsByMention:input_type -> FidRequest
	6,  // 8: HubService.GetReaction:input_type -> ReactionRequest
	7,  // 9: HubService.GetReactionsByFid:input_type -> ReactionsByFidRequest
	8,  // 10: HubService.GetReactionsByCast:input_type -> ReactionsByTargetRequest
	8,  // 11: HubService.GetReactionsByTarget:input_type -> ReactionsByTargetRequest
	9,  // 12: HubService.GetUserData:input_type -> UserDataRequest
	4,  // 13: HubService.GetUserDataByFid:input_type -> FidRequest
	10, // 14: HubService.GetUsernameProof:input_type -> UsernameProofRequest
	4,  // 15: HubService.GetUserNameProofsByFid:input_type -> FidRequest
	11, // 16: HubService.GetVerification:input_type -> VerificationRequest
	4,  // 17: HubService.GetVerificationsByFid:input_type -> FidRequest
	12, // 18: HubService.GetOnChainSigner:input_type -> SignerRequest
	4,  // 19: HubService.GetOnChainSignersByFid:input_type -> FidRequest
	13, // 20: HubService.GetOnChainEvents:input_type -> OnChainEventRequest
	4,  // 21: HubService.GetIdRegistryOnChainEvent:input_type -> FidRequest
	14, // 22: HubService.GetIdRegistryOnChainEventByAddress:input_type -> IdRegistryEventByAddressRequest
	4,  // 23: HubService.GetCurrentStorageLimitsByFid:input_type -> FidRequest
	15, // 24: HubService.GetFids:input_type -> FidsRequest
	16, // 25: HubService.GetLink:input_type -> LinkRequest
	17, // 26: HubService.GetLinksByFid:input_type -> LinksByFidRequest
	18, // 27: HubService.GetLinksByTarget:input_type -> LinksByTargetRequest
	19, // 28: HubService.GetAllCastMessagesByFid:input_type -> FidTimestampRequest
	19, // 29: HubService.GetAllReactionMessagesByFid:input_type -> FidTimestampRequest
	19, // 30: HubService.GetAllVerificationMessagesByFid:input_type -> FidTimestampRequest
	19, // 31: HubService.GetAllUserDataMessagesByFid:input_type -> FidTimestampRequest
	19, // 32: HubService.GetAllLinkMessagesByFid:input_type -> FidTimestampRequest
	4,  // 33: HubService.GetLinkCompactStateMessageByFid:input_type -> FidRequest
	20, // 34: HubService.SubmitBulkMessages:input_type -> SubmitBulkMessagesRequest
	21, // 35: HubService.GetInfo:input_type -> HubInfoRequest
	22, // 36: HubService.GetCurrentPeers:input_type -> Empty
	22, // 37: HubService.StopSync:input_type -> Empty
	23, // 38: HubService.ForceSync:input_type -> SyncStatusRequest
	23, // 39: HubService.GetSyncStatus:input_type -> SyncStatusRequest
	24, // 40: HubService.GetAllSyncIdsByPrefix:input_type -> TrieNodePrefix
	25, // 41: HubService.GetAllMessagesBySyncIds:input_type -> SyncIds
	24, // 42: HubService.GetSyncMetadataByPrefix:input_type -> TrieNodePrefix
	24, // 43: HubService.GetSyncSnapshotByPrefix:input_type -> TrieNodePrefix
	26, // 44: HubService.StreamSync:input_type -> StreamSyncRequest
	27, // 45: HubService.StreamFetch:input_type -> StreamFetchRequest
	22, // 46: AdminService.RebuildSyncTrie:input_type -> Empty
	22, // 47: AdminService.DeleteAllMessagesFromDb:input_type -> Empty
	28, // 48: AdminService.SubmitOnChainEvent:input_type -> OnChainEvent
	0,  // 49: HubService.SubmitMessage:output_type -> Message
	29, // 50: HubService.ValidateMessage:output_type -> ValidationResponse
	30, // 51: HubService.Subscribe:output_type -> HubEvent
	30, // 52: HubService.GetEvent:output_type -> HubEvent
	0,  // 53: HubService.GetCast:output_type -> Message
	31, // 54: HubService.GetCastsByFid:output_type -> MessagesResponse
	31, // 55: HubService.GetCastsByParent:output_type -> MessagesResponse
	31, // 56: HubService.GetCastsByMention:output_type -> MessagesResponse
	0,  // 57: HubService.GetReaction:output_type -> Message
	31, // 58: HubService.GetReactionsByFid:output_type -> MessagesResponse
	31, // 59: HubService.GetReactionsByCast:output_type -> MessagesResponse
	31, // 60: HubService.GetReactionsByTarget:output_type -> MessagesResponse
	0,  // 61: HubService.GetUserData:output_type -> Message
	31, // 62: HubService.GetUserDataByFid:output_type -> MessagesResponse
	32, // 63: HubService.GetUsernameProof:output_type -> UserNameProof
	33, // 64: HubService.GetUserNameProofsByFid:output_type -> UsernameProofsResponse
	0,  // 65: HubService.GetVerification:output_type -> Message
	31, // 66: HubService.GetVerificationsByFid:output_type -> MessagesResponse
	28, // 67: HubService.GetOnChainSigner:output_type -> OnChainEvent
	34, // 68: HubService.GetOnChainSignersByFid:output_type -> OnChainEventResponse
	34, // 69: HubService.GetOnChainEvents:output_type -> OnChainEventResponse
	28, // 70: HubService.GetIdRegistryOnChainEvent:output_type -> OnChainEvent
	28, // 71: HubService.GetIdRegistryOnChainEventByAddress:output_type -> OnChainEvent
	35, // 72: HubService.GetCurrentStorageLimitsByFid:output_type -> StorageLimitsResponse
	36, // 73: HubService.GetFids:output_type -> FidsResponse
	0,  // 74: HubService.GetLink:output_type -> Message
	31, // 75: HubService.GetLinksByFid:output_type -> MessagesResponse
	31, // 76: HubService.GetLinksByTarget:output_type -> MessagesResponse
	31, // 77: HubService.GetAllCastMessagesByFid:output_type -> MessagesResponse
	31, // 78: HubService.GetAllReactionMessagesByFid:output_type -> MessagesResponse
	31, // 79: HubService.GetAllVerificationMessagesByFid:output_type -> MessagesResponse
	31, // 80: HubService.GetAllUserDataMessagesByFid:output_type -> MessagesResponse
	31, // 81: HubService.GetAllLinkMessagesByFid:output_type -> MessagesResponse
	31, // 82: HubService.GetLinkCompactStateMessageByFid:output_type -> MessagesResponse
	37, // 83: HubService.SubmitBulkMessages:output_type -> SubmitBulkMessagesResponse
	38, // 84: HubService.GetInfo:output_type -> HubInfoResponse
	39, // 85: HubService.GetCurrentPeers:output_type -> ContactInfoResponse
	40, // 86: HubService.StopSync:output_type -> SyncStatusResponse
	40, // 87: HubService.ForceSync:output_type -> SyncStatusResponse
	40, // 88: HubService.GetSyncStatus:output_type -> SyncStatusResponse
	25, // 89: HubService.GetAllSyncIdsByPrefix:output_type -> SyncIds
	31, // 90: HubService.GetAllMessagesBySyncIds:output_type -> MessagesResponse
	41, // 91: HubService.GetSyncMetadataByPrefix:output_type -> TrieNodeMetadataResponse
	42, // 92: HubService.GetSyncSnapshotByPrefix:output_type -> TrieNodeSnapshotResponse
	43, // 93: HubService.StreamSync:output_type -> StreamSyncResponse
	44, // 94: HubService.StreamFetch:output_type -> StreamFetchResponse
	22, // 95: AdminService.RebuildSyncTrie:output_type -> Empty
	22, // 96: AdminService.DeleteAllMessagesFromDb:output_type -> Empty
	28, // 97: AdminService.SubmitOnChainEvent:output_type -> OnChainEvent
	49, // [49:98] is the sub-list for method output_type
	0,  // [0:49] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_proto_init() }
func file_rpc_proto_init() {
	if File_rpc_proto != nil {
		return
	}
	file_message_proto_init()
	file_hub_event_proto_init()
	file_request_response_proto_init()
	file_username_proof_proto_init()
	file_onchain_event_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_depIdxs,
	}.Build()
	File_rpc_proto = out.File
	file_rpc_proto_rawDesc = nil
	file_rpc_proto_goTypes = nil
	file_rpc_proto_depIdxs = nil
}
