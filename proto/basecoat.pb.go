// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: basecoat.proto

package proto

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

var File_basecoat_proto protoreflect.FileDescriptor

var file_basecoat_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x62, 0x61, 0x73, 0x65, 0x63, 0x6f, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x62, 0x61, 0x73, 0x65, 0x63, 0x6f, 0x61,
	0x74, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x32, 0xd3, 0x1a, 0x0a, 0x08, 0x42, 0x61, 0x73, 0x65, 0x63, 0x6f, 0x61, 0x74, 0x12, 0x4d,
	0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x50, 0x49, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x50, 0x49, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x50, 0x49,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x47, 0x65, 0x74,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0c,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a,
	0x12, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x67, 0x67,
	0x6c, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f,
	0x67, 0x67, 0x6c, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x46,
	0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x6d,
	0x75, 0x6c, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0c, 0x4c,
	0x69, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x6f,
	0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x68, 0x0a, 0x17, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72,
	0x6d, 0x75, 0x6c, 0x61, 0x57, 0x69, 0x74, 0x68, 0x4a, 0x6f, 0x62, 0x12, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72,
	0x6d, 0x75, 0x6c, 0x61, 0x57, 0x69, 0x74, 0x68, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x73, 0x73, 0x6f, 0x63,
	0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x57, 0x69, 0x74, 0x68, 0x4a,
	0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x71, 0x0a, 0x1a, 0x44, 0x69,
	0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c,
	0x61, 0x46, 0x72, 0x6f, 0x6d, 0x4a, 0x6f, 0x62, 0x12, 0x28, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72,
	0x6d, 0x75, 0x6c, 0x61, 0x46, 0x72, 0x6f, 0x6d, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x61, 0x73,
	0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x46, 0x72,
	0x6f, 0x6d, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a,
	0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x1b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72,
	0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x65,
	0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x73, 0x65, 0x73, 0x12, 0x17, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x61, 0x73, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x42, 0x61, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x41, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x12, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x6b, 0x0a, 0x18, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x42,
	0x61, 0x73, 0x65, 0x57, 0x69, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x26,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65,
	0x42, 0x61, 0x73, 0x65, 0x57, 0x69, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41,
	0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x57, 0x69, 0x74, 0x68,
	0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x74, 0x0a, 0x1b, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x42,
	0x61, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x29,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69,
	0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75,
	0x6c, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x42, 0x61,
	0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42,
	0x61, 0x73, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42,
	0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e,
	0x74, 0x73, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a,
	0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x12,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x77, 0x0a, 0x1c,
	0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e,
	0x74, 0x57, 0x69, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x2a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x57, 0x69, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61,
	0x6e, 0x74, 0x57, 0x69, 0x74, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x1f, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73,
	0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x46, 0x72,
	0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c,
	0x6f, 0x72, 0x61, 0x6e, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0c, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a,
	0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12,
	0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x50, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x10, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x1e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53,
	0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4a,
	0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x4c, 0x69,
	0x73, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4a, 0x6f, 0x62, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4a, 0x6f, 0x62, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4a, 0x6f, 0x62, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x69, 0x6e, 0x74, 0x6a, 0x65, 0x64, 0x77, 0x61,
	0x72, 0x64, 0x73, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x63, 0x6f, 0x61, 0x74, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_basecoat_proto_goTypes = []interface{}{
	(*CreateAPITokenRequest)(nil),                   // 0: proto.CreateAPITokenRequest
	(*GetSystemInfoRequest)(nil),                    // 1: proto.GetSystemInfoRequest
	(*GetAccountRequest)(nil),                       // 2: proto.GetAccountRequest
	(*ListAccountsRequest)(nil),                     // 3: proto.ListAccountsRequest
	(*CreateAccountRequest)(nil),                    // 4: proto.CreateAccountRequest
	(*UpdateAccountRequest)(nil),                    // 5: proto.UpdateAccountRequest
	(*ToggleAccountStateRequest)(nil),               // 6: proto.ToggleAccountStateRequest
	(*GetFormulaRequest)(nil),                       // 7: proto.GetFormulaRequest
	(*ListFormulasRequest)(nil),                     // 8: proto.ListFormulasRequest
	(*CreateFormulaRequest)(nil),                    // 9: proto.CreateFormulaRequest
	(*AssociateFormulaWithJobRequest)(nil),          // 10: proto.AssociateFormulaWithJobRequest
	(*DisassociateFormulaFromJobRequest)(nil),       // 11: proto.DisassociateFormulaFromJobRequest
	(*UpdateFormulaRequest)(nil),                    // 12: proto.UpdateFormulaRequest
	(*DeleteFormulaRequest)(nil),                    // 13: proto.DeleteFormulaRequest
	(*GetBaseRequest)(nil),                          // 14: proto.GetBaseRequest
	(*ListBasesRequest)(nil),                        // 15: proto.ListBasesRequest
	(*CreateBaseRequest)(nil),                       // 16: proto.CreateBaseRequest
	(*AssociateBaseWithFormulaRequest)(nil),         // 17: proto.AssociateBaseWithFormulaRequest
	(*DisassociateBaseFromFormulaRequest)(nil),      // 18: proto.DisassociateBaseFromFormulaRequest
	(*UpdateBaseRequest)(nil),                       // 19: proto.UpdateBaseRequest
	(*DeleteBaseRequest)(nil),                       // 20: proto.DeleteBaseRequest
	(*GetColorantRequest)(nil),                      // 21: proto.GetColorantRequest
	(*ListColorantsRequest)(nil),                    // 22: proto.ListColorantsRequest
	(*CreateColorantRequest)(nil),                   // 23: proto.CreateColorantRequest
	(*AssociateColorantWithFormulaRequest)(nil),     // 24: proto.AssociateColorantWithFormulaRequest
	(*DisassociateColorantFromFormulaRequest)(nil),  // 25: proto.DisassociateColorantFromFormulaRequest
	(*UpdateColorantRequest)(nil),                   // 26: proto.UpdateColorantRequest
	(*DeleteColorantRequest)(nil),                   // 27: proto.DeleteColorantRequest
	(*GetContactRequest)(nil),                       // 28: proto.GetContactRequest
	(*ListContactsRequest)(nil),                     // 29: proto.ListContactsRequest
	(*CreateContactRequest)(nil),                    // 30: proto.CreateContactRequest
	(*UpdateContactRequest)(nil),                    // 31: proto.UpdateContactRequest
	(*DeleteContactRequest)(nil),                    // 32: proto.DeleteContactRequest
	(*GetContractorRequest)(nil),                    // 33: proto.GetContractorRequest
	(*ListContractorsRequest)(nil),                  // 34: proto.ListContractorsRequest
	(*CreateContractorRequest)(nil),                 // 35: proto.CreateContractorRequest
	(*UpdateContractorRequest)(nil),                 // 36: proto.UpdateContractorRequest
	(*DeleteContractorRequest)(nil),                 // 37: proto.DeleteContractorRequest
	(*GetJobRequest)(nil),                           // 38: proto.GetJobRequest
	(*ListJobsRequest)(nil),                         // 39: proto.ListJobsRequest
	(*CreateJobRequest)(nil),                        // 40: proto.CreateJobRequest
	(*UpdateJobRequest)(nil),                        // 41: proto.UpdateJobRequest
	(*DeleteJobRequest)(nil),                        // 42: proto.DeleteJobRequest
	(*CreateAPITokenResponse)(nil),                  // 43: proto.CreateAPITokenResponse
	(*GetSystemInfoResponse)(nil),                   // 44: proto.GetSystemInfoResponse
	(*GetAccountResponse)(nil),                      // 45: proto.GetAccountResponse
	(*ListAccountsResponse)(nil),                    // 46: proto.ListAccountsResponse
	(*CreateAccountResponse)(nil),                   // 47: proto.CreateAccountResponse
	(*UpdateAccountResponse)(nil),                   // 48: proto.UpdateAccountResponse
	(*ToggleAccountStateResponse)(nil),              // 49: proto.ToggleAccountStateResponse
	(*GetFormulaResponse)(nil),                      // 50: proto.GetFormulaResponse
	(*ListFormulasResponse)(nil),                    // 51: proto.ListFormulasResponse
	(*CreateFormulaResponse)(nil),                   // 52: proto.CreateFormulaResponse
	(*AssociateFormulaWithJobResponse)(nil),         // 53: proto.AssociateFormulaWithJobResponse
	(*DisassociateFormulaFromJobResponse)(nil),      // 54: proto.DisassociateFormulaFromJobResponse
	(*UpdateFormulaResponse)(nil),                   // 55: proto.UpdateFormulaResponse
	(*DeleteFormulaResponse)(nil),                   // 56: proto.DeleteFormulaResponse
	(*GetBaseResponse)(nil),                         // 57: proto.GetBaseResponse
	(*ListBasesResponse)(nil),                       // 58: proto.ListBasesResponse
	(*CreateBaseResponse)(nil),                      // 59: proto.CreateBaseResponse
	(*AssociateBaseWithFormulaResponse)(nil),        // 60: proto.AssociateBaseWithFormulaResponse
	(*DisassociateBaseFromFormulaResponse)(nil),     // 61: proto.DisassociateBaseFromFormulaResponse
	(*UpdateBaseResponse)(nil),                      // 62: proto.UpdateBaseResponse
	(*DeleteBaseResponse)(nil),                      // 63: proto.DeleteBaseResponse
	(*GetColorantResponse)(nil),                     // 64: proto.GetColorantResponse
	(*ListColorantsResponse)(nil),                   // 65: proto.ListColorantsResponse
	(*CreateColorantResponse)(nil),                  // 66: proto.CreateColorantResponse
	(*AssociateColorantWithFormulaResponse)(nil),    // 67: proto.AssociateColorantWithFormulaResponse
	(*DisassociateColorantFromFormulaResponse)(nil), // 68: proto.DisassociateColorantFromFormulaResponse
	(*UpdateColorantResponse)(nil),                  // 69: proto.UpdateColorantResponse
	(*DeleteColorantResponse)(nil),                  // 70: proto.DeleteColorantResponse
	(*GetContactResponse)(nil),                      // 71: proto.GetContactResponse
	(*ListContactsResponse)(nil),                    // 72: proto.ListContactsResponse
	(*CreateContactResponse)(nil),                   // 73: proto.CreateContactResponse
	(*UpdateContactResponse)(nil),                   // 74: proto.UpdateContactResponse
	(*DeleteContactResponse)(nil),                   // 75: proto.DeleteContactResponse
	(*GetContractorResponse)(nil),                   // 76: proto.GetContractorResponse
	(*ListContractorsResponse)(nil),                 // 77: proto.ListContractorsResponse
	(*CreateContractorResponse)(nil),                // 78: proto.CreateContractorResponse
	(*UpdateContractorResponse)(nil),                // 79: proto.UpdateContractorResponse
	(*DeleteContractorResponse)(nil),                // 80: proto.DeleteContractorResponse
	(*GetJobResponse)(nil),                          // 81: proto.GetJobResponse
	(*ListJobsResponse)(nil),                        // 82: proto.ListJobsResponse
	(*CreateJobResponse)(nil),                       // 83: proto.CreateJobResponse
	(*UpdateJobResponse)(nil),                       // 84: proto.UpdateJobResponse
	(*DeleteJobResponse)(nil),                       // 85: proto.DeleteJobResponse
}
var file_basecoat_proto_depIdxs = []int32{
	0,  // 0: proto.Basecoat.CreateAPIToken:input_type -> proto.CreateAPITokenRequest
	1,  // 1: proto.Basecoat.GetSystemInfo:input_type -> proto.GetSystemInfoRequest
	2,  // 2: proto.Basecoat.GetAccount:input_type -> proto.GetAccountRequest
	3,  // 3: proto.Basecoat.ListAccounts:input_type -> proto.ListAccountsRequest
	4,  // 4: proto.Basecoat.CreateAccount:input_type -> proto.CreateAccountRequest
	5,  // 5: proto.Basecoat.UpdateAccount:input_type -> proto.UpdateAccountRequest
	6,  // 6: proto.Basecoat.ToggleAccountState:input_type -> proto.ToggleAccountStateRequest
	7,  // 7: proto.Basecoat.GetFormula:input_type -> proto.GetFormulaRequest
	8,  // 8: proto.Basecoat.ListFormulas:input_type -> proto.ListFormulasRequest
	9,  // 9: proto.Basecoat.CreateFormula:input_type -> proto.CreateFormulaRequest
	10, // 10: proto.Basecoat.AssociateFormulaWithJob:input_type -> proto.AssociateFormulaWithJobRequest
	11, // 11: proto.Basecoat.DisassociateFormulaFromJob:input_type -> proto.DisassociateFormulaFromJobRequest
	12, // 12: proto.Basecoat.UpdateFormula:input_type -> proto.UpdateFormulaRequest
	13, // 13: proto.Basecoat.DeleteFormula:input_type -> proto.DeleteFormulaRequest
	14, // 14: proto.Basecoat.GetBase:input_type -> proto.GetBaseRequest
	15, // 15: proto.Basecoat.ListBases:input_type -> proto.ListBasesRequest
	16, // 16: proto.Basecoat.CreateBase:input_type -> proto.CreateBaseRequest
	17, // 17: proto.Basecoat.AssociateBaseWithFormula:input_type -> proto.AssociateBaseWithFormulaRequest
	18, // 18: proto.Basecoat.DisassociateBaseFromFormula:input_type -> proto.DisassociateBaseFromFormulaRequest
	19, // 19: proto.Basecoat.UpdateBase:input_type -> proto.UpdateBaseRequest
	20, // 20: proto.Basecoat.DeleteBase:input_type -> proto.DeleteBaseRequest
	21, // 21: proto.Basecoat.GetColorant:input_type -> proto.GetColorantRequest
	22, // 22: proto.Basecoat.ListColorants:input_type -> proto.ListColorantsRequest
	23, // 23: proto.Basecoat.CreateColorant:input_type -> proto.CreateColorantRequest
	24, // 24: proto.Basecoat.AssociateColorantWithFormula:input_type -> proto.AssociateColorantWithFormulaRequest
	25, // 25: proto.Basecoat.DisassociateColorantFromFormula:input_type -> proto.DisassociateColorantFromFormulaRequest
	26, // 26: proto.Basecoat.UpdateColorant:input_type -> proto.UpdateColorantRequest
	27, // 27: proto.Basecoat.DeleteColorant:input_type -> proto.DeleteColorantRequest
	28, // 28: proto.Basecoat.GetContact:input_type -> proto.GetContactRequest
	29, // 29: proto.Basecoat.ListContacts:input_type -> proto.ListContactsRequest
	30, // 30: proto.Basecoat.CreateContact:input_type -> proto.CreateContactRequest
	31, // 31: proto.Basecoat.UpdateContact:input_type -> proto.UpdateContactRequest
	32, // 32: proto.Basecoat.DeleteContact:input_type -> proto.DeleteContactRequest
	33, // 33: proto.Basecoat.GetContractor:input_type -> proto.GetContractorRequest
	34, // 34: proto.Basecoat.ListContractors:input_type -> proto.ListContractorsRequest
	35, // 35: proto.Basecoat.CreateContractor:input_type -> proto.CreateContractorRequest
	36, // 36: proto.Basecoat.UpdateContractor:input_type -> proto.UpdateContractorRequest
	37, // 37: proto.Basecoat.DeleteContractor:input_type -> proto.DeleteContractorRequest
	38, // 38: proto.Basecoat.GetJob:input_type -> proto.GetJobRequest
	39, // 39: proto.Basecoat.ListJobs:input_type -> proto.ListJobsRequest
	40, // 40: proto.Basecoat.CreateJob:input_type -> proto.CreateJobRequest
	41, // 41: proto.Basecoat.UpdateJob:input_type -> proto.UpdateJobRequest
	42, // 42: proto.Basecoat.DeleteJob:input_type -> proto.DeleteJobRequest
	43, // 43: proto.Basecoat.CreateAPIToken:output_type -> proto.CreateAPITokenResponse
	44, // 44: proto.Basecoat.GetSystemInfo:output_type -> proto.GetSystemInfoResponse
	45, // 45: proto.Basecoat.GetAccount:output_type -> proto.GetAccountResponse
	46, // 46: proto.Basecoat.ListAccounts:output_type -> proto.ListAccountsResponse
	47, // 47: proto.Basecoat.CreateAccount:output_type -> proto.CreateAccountResponse
	48, // 48: proto.Basecoat.UpdateAccount:output_type -> proto.UpdateAccountResponse
	49, // 49: proto.Basecoat.ToggleAccountState:output_type -> proto.ToggleAccountStateResponse
	50, // 50: proto.Basecoat.GetFormula:output_type -> proto.GetFormulaResponse
	51, // 51: proto.Basecoat.ListFormulas:output_type -> proto.ListFormulasResponse
	52, // 52: proto.Basecoat.CreateFormula:output_type -> proto.CreateFormulaResponse
	53, // 53: proto.Basecoat.AssociateFormulaWithJob:output_type -> proto.AssociateFormulaWithJobResponse
	54, // 54: proto.Basecoat.DisassociateFormulaFromJob:output_type -> proto.DisassociateFormulaFromJobResponse
	55, // 55: proto.Basecoat.UpdateFormula:output_type -> proto.UpdateFormulaResponse
	56, // 56: proto.Basecoat.DeleteFormula:output_type -> proto.DeleteFormulaResponse
	57, // 57: proto.Basecoat.GetBase:output_type -> proto.GetBaseResponse
	58, // 58: proto.Basecoat.ListBases:output_type -> proto.ListBasesResponse
	59, // 59: proto.Basecoat.CreateBase:output_type -> proto.CreateBaseResponse
	60, // 60: proto.Basecoat.AssociateBaseWithFormula:output_type -> proto.AssociateBaseWithFormulaResponse
	61, // 61: proto.Basecoat.DisassociateBaseFromFormula:output_type -> proto.DisassociateBaseFromFormulaResponse
	62, // 62: proto.Basecoat.UpdateBase:output_type -> proto.UpdateBaseResponse
	63, // 63: proto.Basecoat.DeleteBase:output_type -> proto.DeleteBaseResponse
	64, // 64: proto.Basecoat.GetColorant:output_type -> proto.GetColorantResponse
	65, // 65: proto.Basecoat.ListColorants:output_type -> proto.ListColorantsResponse
	66, // 66: proto.Basecoat.CreateColorant:output_type -> proto.CreateColorantResponse
	67, // 67: proto.Basecoat.AssociateColorantWithFormula:output_type -> proto.AssociateColorantWithFormulaResponse
	68, // 68: proto.Basecoat.DisassociateColorantFromFormula:output_type -> proto.DisassociateColorantFromFormulaResponse
	69, // 69: proto.Basecoat.UpdateColorant:output_type -> proto.UpdateColorantResponse
	70, // 70: proto.Basecoat.DeleteColorant:output_type -> proto.DeleteColorantResponse
	71, // 71: proto.Basecoat.GetContact:output_type -> proto.GetContactResponse
	72, // 72: proto.Basecoat.ListContacts:output_type -> proto.ListContactsResponse
	73, // 73: proto.Basecoat.CreateContact:output_type -> proto.CreateContactResponse
	74, // 74: proto.Basecoat.UpdateContact:output_type -> proto.UpdateContactResponse
	75, // 75: proto.Basecoat.DeleteContact:output_type -> proto.DeleteContactResponse
	76, // 76: proto.Basecoat.GetContractor:output_type -> proto.GetContractorResponse
	77, // 77: proto.Basecoat.ListContractors:output_type -> proto.ListContractorsResponse
	78, // 78: proto.Basecoat.CreateContractor:output_type -> proto.CreateContractorResponse
	79, // 79: proto.Basecoat.UpdateContractor:output_type -> proto.UpdateContractorResponse
	80, // 80: proto.Basecoat.DeleteContractor:output_type -> proto.DeleteContractorResponse
	81, // 81: proto.Basecoat.GetJob:output_type -> proto.GetJobResponse
	82, // 82: proto.Basecoat.ListJobs:output_type -> proto.ListJobsResponse
	83, // 83: proto.Basecoat.CreateJob:output_type -> proto.CreateJobResponse
	84, // 84: proto.Basecoat.UpdateJob:output_type -> proto.UpdateJobResponse
	85, // 85: proto.Basecoat.DeleteJob:output_type -> proto.DeleteJobResponse
	43, // [43:86] is the sub-list for method output_type
	0,  // [0:43] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_basecoat_proto_init() }
func file_basecoat_proto_init() {
	if File_basecoat_proto != nil {
		return
	}
	file_basecoat_transport_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_basecoat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_basecoat_proto_goTypes,
		DependencyIndexes: file_basecoat_proto_depIdxs,
	}.Build()
	File_basecoat_proto = out.File
	file_basecoat_proto_rawDesc = nil
	file_basecoat_proto_goTypes = nil
	file_basecoat_proto_depIdxs = nil
}