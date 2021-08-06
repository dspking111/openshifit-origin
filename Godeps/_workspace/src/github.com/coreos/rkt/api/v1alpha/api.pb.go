// Code generated by protoc-gen-go.
// source: api.proto
// DO NOT EDIT!

/*
Package v1alpha is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	ImageFormat
	Image
	Network
	App
	Pod
	KeyValue
	PodFilter
	ImageFilter
	Info
	Event
	EventFilter
	GetInfoRequest
	GetInfoResponse
	ListPodsRequest
	ListPodsResponse
	InspectPodRequest
	InspectPodResponse
	ListImagesRequest
	ListImagesResponse
	InspectImageRequest
	InspectImageResponse
	ListenEventsRequest
	ListenEventsResponse
	GetLogsRequest
	GetLogsResponse
*/
package v1alpha

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// ImageType defines the supported image type.
type ImageType int32

const (
	ImageType_IMAGE_TYPE_UNDEFINED ImageType = 0
	ImageType_IMAGE_TYPE_APPC      ImageType = 1
	ImageType_IMAGE_TYPE_DOCKER    ImageType = 2
	ImageType_IMAGE_TYPE_OCI       ImageType = 3
)

var ImageType_name = map[int32]string{
	0: "IMAGE_TYPE_UNDEFINED",
	1: "IMAGE_TYPE_APPC",
	2: "IMAGE_TYPE_DOCKER",
	3: "IMAGE_TYPE_OCI",
}
var ImageType_value = map[string]int32{
	"IMAGE_TYPE_UNDEFINED": 0,
	"IMAGE_TYPE_APPC":      1,
	"IMAGE_TYPE_DOCKER":    2,
	"IMAGE_TYPE_OCI":       3,
}

func (x ImageType) String() string {
	return proto.EnumName(ImageType_name, int32(x))
}

// AppState defines the possible states of the app.
type AppState int32

const (
	AppState_APP_STATE_UNDEFINED AppState = 0
	AppState_APP_STATE_RUNNING   AppState = 1
	AppState_APP_STATE_EXITED    AppState = 2
)

var AppState_name = map[int32]string{
	0: "APP_STATE_UNDEFINED",
	1: "APP_STATE_RUNNING",
	2: "APP_STATE_EXITED",
}
var AppState_value = map[string]int32{
	"APP_STATE_UNDEFINED": 0,
	"APP_STATE_RUNNING":   1,
	"APP_STATE_EXITED":    2,
}

func (x AppState) String() string {
	return proto.EnumName(AppState_name, int32(x))
}

// PodState defines the possible states of the pod.
// See https://github.com/coreos/rkt/blob/master/Documentation/devel/pod-lifecycle.md for a detailed
// explanation of each state.
type PodState int32

const (
	PodState_POD_STATE_UNDEFINED PodState = 0
	// States before the pod is running.
	PodState_POD_STATE_EMBRYO    PodState = 1
	PodState_POD_STATE_PREPARING PodState = 2
	PodState_POD_STATE_PREPARED  PodState = 3
	// State that indicates the pod is running.
	PodState_POD_STATE_RUNNING PodState = 4
	// States that indicates the pod is exited, and will never run.
	PodState_POD_STATE_ABORTED_PREPARE PodState = 5
	PodState_POD_STATE_EXITED          PodState = 6
	PodState_POD_STATE_DELETING        PodState = 7
	PodState_POD_STATE_GARBAGE         PodState = 8
)

var PodState_name = map[int32]string{
	0: "POD_STATE_UNDEFINED",
	1: "POD_STATE_EMBRYO",
	2: "POD_STATE_PREPARING",
	3: "POD_STATE_PREPARED",
	4: "POD_STATE_RUNNING",
	5: "POD_STATE_ABORTED_PREPARE",
	6: "POD_STATE_EXITED",
	7: "POD_STATE_DELETING",
	8: "POD_STATE_GARBAGE",
}
var PodState_value = map[string]int32{
	"POD_STATE_UNDEFINED":       0,
	"POD_STATE_EMBRYO":          1,
	"POD_STATE_PREPARING":       2,
	"POD_STATE_PREPARED":        3,
	"POD_STATE_RUNNING":         4,
	"POD_STATE_ABORTED_PREPARE": 5,
	"POD_STATE_EXITED":          6,
	"POD_STATE_DELETING":        7,
	"POD_STATE_GARBAGE":         8,
}

func (x PodState) String() string {
	return proto.EnumName(PodState_name, int32(x))
}

// EventType defines the type of the events that will be received via ListenEvents().
type EventType int32

const (
	EventType_EVENT_TYPE_UNDEFINED EventType = 0
	// Pod events.
	EventType_EVENT_TYPE_POD_PREPARED          EventType = 1
	EventType_EVENT_TYPE_POD_PREPARE_ABORTED   EventType = 2
	EventType_EVENT_TYPE_POD_STARTED           EventType = 3
	EventType_EVENT_TYPE_POD_EXITED            EventType = 4
	EventType_EVENT_TYPE_POD_GARBAGE_COLLECTED EventType = 5
	// App events.
	EventType_EVENT_TYPE_APP_STARTED EventType = 6
	EventType_EVENT_TYPE_APP_EXITED  EventType = 7
	// Image events.
	EventType_EVENT_TYPE_IMAGE_IMPORTED EventType = 8
	EventType_EVENT_TYPE_IMAGE_REMOVED  EventType = 9
)

var EventType_name = map[int32]string{
	0: "EVENT_TYPE_UNDEFINED",
	1: "EVENT_TYPE_POD_PREPARED",
	2: "EVENT_TYPE_POD_PREPARE_ABORTED",
	3: "EVENT_TYPE_POD_STARTED",
	4: "EVENT_TYPE_POD_EXITED",
	5: "EVENT_TYPE_POD_GARBAGE_COLLECTED",
	6: "EVENT_TYPE_APP_STARTED",
	7: "EVENT_TYPE_APP_EXITED",
	8: "EVENT_TYPE_IMAGE_IMPORTED",
	9: "EVENT_TYPE_IMAGE_REMOVED",
}
var EventType_value = map[string]int32{
	"EVENT_TYPE_UNDEFINED":             0,
	"EVENT_TYPE_POD_PREPARED":          1,
	"EVENT_TYPE_POD_PREPARE_ABORTED":   2,
	"EVENT_TYPE_POD_STARTED":           3,
	"EVENT_TYPE_POD_EXITED":            4,
	"EVENT_TYPE_POD_GARBAGE_COLLECTED": 5,
	"EVENT_TYPE_APP_STARTED":           6,
	"EVENT_TYPE_APP_EXITED":            7,
	"EVENT_TYPE_IMAGE_IMPORTED":        8,
	"EVENT_TYPE_IMAGE_REMOVED":         9,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}

// ImageFormat defines the format of the image.
type ImageFormat struct {
	// Type of the image, required.
	Type ImageType `protobuf:"varint,1,opt,name=type,enum=v1alpha.ImageType" json:"type,omitempty"`
	// Version of the image format, required.
	Version string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
}

func (m *ImageFormat) Reset()         { *m = ImageFormat{} }
func (m *ImageFormat) String() string { return proto.CompactTextString(m) }
func (*ImageFormat) ProtoMessage()    {}

// Image describes the image's information.
type Image struct {
	// Base format of the image, required. This indicates the original format
	// for the image as nowadays all the image formats will be transformed to
	// ACI.
	BaseFormat *ImageFormat `protobuf:"bytes,1,opt,name=base_format" json:"base_format,omitempty"`
	// ID of the image, a string that can be used to uniquely identify the image,
	// e.g. sha512 hash of the ACIs, required.
	Id string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	// Name of the image in the image manifest, e.g. 'coreos.com/etcd', optional.
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	// Version of the image, e.g. 'latest', '2.0.10', optional.
	Version string `protobuf:"bytes,4,opt,name=version" json:"version,omitempty"`
	// Timestamp of when the image is imported, it is the seconds since epoch, optional.
	ImportTimestamp int64 `protobuf:"varint,5,opt,name=import_timestamp" json:"import_timestamp,omitempty"`
	// JSON-encoded byte array that represents the image manifest, optional.
	Manifest []byte `protobuf:"bytes,6,opt,name=manifest,proto3" json:"manifest,omitempty"`
}

func (m *Image) Reset()         { *m = Image{} }
func (m *Image) String() string { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()    {}

func (m *Image) GetBaseFormat() *ImageFormat {
	if m != nil {
		return m.BaseFormat
	}
	return nil
}

// Network describes the network information of a pod.
type Network struct {
	// Name of the network that a pod belongs to, required.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Pod's IPv4 address within the network, optional if IPv6 address is given.
	Ipv4 string `protobuf:"bytes,2,opt,name=ipv4" json:"ipv4,omitempty"`
	// Pod's IPv6 address within the network, optional if IPv4 address is given.
	Ipv6 string `protobuf:"bytes,3,opt,name=ipv6" json:"ipv6,omitempty"`
}

func (m *Network) Reset()         { *m = Network{} }
func (m *Network) String() string { return proto.CompactTextString(m) }
func (*Network) ProtoMessage()    {}

// App describes the information of an app that's running in a pod.
type App struct {
	// Name of the app, required.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Image used by the app, required. However, this may only contain the image id
	// if it is returned by ListPods().
	Image *Image `protobuf:"bytes,2,opt,name=image" json:"image,omitempty"`
	// State of the app. optional, non-empty only if it's returned by InspectPod().
	State AppState `protobuf:"varint,3,opt,name=state,enum=v1alpha.AppState" json:"state,omitempty"`
	// Exit code of the app. optional, only valid if it's returned by InspectPod() and
	// the app has already exited.
	ExitCode int32 `protobuf:"zigzag32,4,opt,name=exit_code" json:"exit_code,omitempty"`
}

func (m *App) Reset()         { *m = App{} }
func (m *App) String() string { return proto.CompactTextString(m) }
func (*App) ProtoMessage()    {}

func (m *App) GetImage() *Image {
	if m != nil {
		return m.Image
	}
	return nil
}

// Pod describes a pod's information.
type Pod struct {
	// ID of the pod, in the form of a UUID, required.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// PID of the pod, optional, only valid if it's returned by InspectPod(). A negative value means the pod has exited.
	Pid int32 `protobuf:"zigzag32,2,opt,name=pid" json:"pid,omitempty"`
	// State of the pod, required.
	State PodState `protobuf:"varint,3,opt,name=state,enum=v1alpha.PodState" json:"state,omitempty"`
	// List of apps in the pod, required.
	Apps []*App `protobuf:"bytes,4,rep,name=apps" json:"apps,omitempty"`
	// Network information of the pod, optional, non-empty if the pod is running in private net.
	// Note that a pod can be in multiple networks.
	Networks []*Network `protobuf:"bytes,5,rep,name=networks" json:"networks,omitempty"`
	// JSON-encoded byte array that represents the pod manifest of the pod, required.
	Manifest []byte `protobuf:"bytes,6,opt,name=manifest,proto3" json:"manifest,omitempty"`
}

func (m *Pod) Reset()         { *m = Pod{} }
func (m *Pod) String() string { return proto.CompactTextString(m) }
func (*Pod) ProtoMessage()    {}

func (m *Pod) GetApps() []*App {
	if m != nil {
		return m.Apps
	}
	return nil
}

func (m *Pod) GetNetworks() []*Network {
	if m != nil {
		return m.Networks
	}
	return nil
}

type KeyValue struct {
	// Key part of the key-value pair.
	Key string `protobuf:"bytes,1,opt" json:"Key,omitempty"`
	// Value part of the key-value pair.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *KeyValue) Reset()         { *m = KeyValue{} }
func (m *KeyValue) String() string { return proto.CompactTextString(m) }
func (*KeyValue) ProtoMessage()    {}

// PodFilter defines the condition that the returned pods need to satisfy in ListPods().
// The conditions are combined by 'AND'.
type PodFilter struct {
	// If not empty, the pods that have any of the ids will be returned.
	Ids []string `protobuf:"bytes,1,rep,name=ids" json:"ids,omitempty"`
	// If not empty, the pods that have any of the states will be returned.
	States []PodState `protobuf:"varint,2,rep,name=states,enum=v1alpha.PodState" json:"states,omitempty"`
	// If not empty, the pods that have any of the apps will be returned.
	AppNames []string `protobuf:"bytes,3,rep,name=app_names" json:"app_names,omitempty"`
	// If not empty, the pods that have any of the images(in the apps) will be returned
	ImageIds []string `protobuf:"bytes,4,rep,name=image_ids" json:"image_ids,omitempty"`
	// If not empty, the pods that are in any of the networks will be returned.
	NetworkNames []string `protobuf:"bytes,5,rep,name=network_names" json:"network_names,omitempty"`
	// If not empty, the pods that have any of the annotations will be returned.
	Annotations []*KeyValue `protobuf:"bytes,6,rep,name=annotations" json:"annotations,omitempty"`
}

func (m *PodFilter) Reset()         { *m = PodFilter{} }
func (m *PodFilter) String() string { return proto.CompactTextString(m) }
func (*PodFilter) ProtoMessage()    {}

func (m *PodFilter) GetAnnotations() []*KeyValue {
	if m != nil {
		return m.Annotations
	}
	return nil
}

// ImageFilter defines the condition that the returned images need to satisfy in ListImages().
// The conditions are combined by 'AND'.
type ImageFilter struct {
	// If not empty, the images that have any of the ids will be returned.
	Ids []string `protobuf:"bytes,1,rep,name=ids" json:"ids,omitempty"`
	// if not empty, the images that have any of the prefixes in the name will be returned.
	Prefixes []string `protobuf:"bytes,2,rep,name=prefixes" json:"prefixes,omitempty"`
	// If not empty, the images that have any of the base names will be returned.
	// For example, both 'coreos.com/etcd' and 'k8s.io/etcd' will be returned if 'etcd' is included,
	// however 'k8s.io/etcd-backup' will not be returned.
	BaseNames []string `protobuf:"bytes,3,rep,name=base_names" json:"base_names,omitempty"`
	// If not empty, the images that have any of the keywords in the name will be returned.
	// For example, both 'kubernetes-etcd', 'etcd:latest' will be returned if 'etcd' is included,
	Keywords []string `protobuf:"bytes,4,rep,name=keywords" json:"keywords,omitempty"`
	// If not empty, the images that have any of the labels will be returned.
	Labels []*KeyValue `protobuf:"bytes,5,rep,name=labels" json:"labels,omitempty"`
	// If set, the images that are imported after this timestamp will be returned.
	ImportedAfter int64 `protobuf:"varint,6,opt,name=imported_after" json:"imported_after,omitempty"`
	// If set, the images that are imported before this timestamp will be returned.
	ImportedBefore int64 `protobuf:"varint,7,opt,name=imported_before" json:"imported_before,omitempty"`
	// If not empty, the images that have any of the annotations will be returned.
	Annotations []*KeyValue `protobuf:"bytes,8,rep,name=annotations" json:"annotations,omitempty"`
}

func (m *ImageFilter) Reset()         { *m = ImageFilter{} }
func (m *ImageFilter) String() string { return proto.CompactTextString(m) }
func (*ImageFilter) ProtoMessage()    {}

func (m *ImageFilter) GetLabels() []*KeyValue {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *ImageFilter) GetAnnotations() []*KeyValue {
	if m != nil {
		return m.Annotations
	}
	return nil
}

// Info describes the information of rkt on the machine.
type Info struct {
	// Version of rkt, required, in the form of Semantic Versioning 2.0.0 (http://semver.org/).
	RktVersion string `protobuf:"bytes,1,opt,name=rkt_version" json:"rkt_version,omitempty"`
	// Version of appc, required, in the form of Semantic Versioning 2.0.0 (http://semver.org/).
	AppcVersion string `protobuf:"bytes,2,opt,name=appc_version" json:"appc_version,omitempty"`
	// Latest version of the api that's supported by the service, required, in the form of Semantic Versioning 2.0.0 (http://semver.org/).
	ApiVersion string `protobuf:"bytes,3,opt,name=api_version" json:"api_version,omitempty"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}

// Event describes the events that will be received via ListenEvents().
type Event struct {
	// Type of the event, required.
	Type EventType `protobuf:"varint,1,opt,name=type,enum=v1alpha.EventType" json:"type,omitempty"`
	// ID of the subject that causes the event, required.
	// If the event is a pod or app event, the id is the pod's uuid.
	// If the event is an image event, the id is the image's id.
	Id string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	// Name of the subject that causes the event, required.
	// If the event is a pod event, the name is the pod's name.
	// If the event is an app event, the name is the app's name.
	// If the event is an image event, the name is the image's name.
	From string `protobuf:"bytes,3,opt,name=from" json:"from,omitempty"`
	// Timestamp of when the event happens, it is the seconds since epoch, required.
	Time int64 `protobuf:"varint,4,opt,name=time" json:"time,omitempty"`
	// Data of the event, in the form of key-value pairs, optional.
	Data []*KeyValue `protobuf:"bytes,5,rep,name=data" json:"data,omitempty"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

func (m *Event) GetData() []*KeyValue {
	if m != nil {
		return m.Data
	}
	return nil
}

// EventFilter defines the condition that the returned events needs to satisfy in ListImages().
// The condition are combined by 'AND'.
type EventFilter struct {
	// If not empty, then only returns the events that have the listed types.
	Types []EventType `protobuf:"varint,1,rep,name=types,enum=v1alpha.EventType" json:"types,omitempty"`
	// If not empty, then only returns the events whose 'id' is included in the listed ids.
	Ids []string `protobuf:"bytes,2,rep,name=ids" json:"ids,omitempty"`
	// If not empty, then only returns the events whose 'from' is included in the listed names.
	Names []string `protobuf:"bytes,3,rep,name=names" json:"names,omitempty"`
	// If set, then only returns the events after this timestamp.
	// If the server starts after since_time, then only the events happened after the start of the server will be returned.
	// If since_time is a future timestamp, then no events will be returned until that time.
	SinceTime int64 `protobuf:"varint,4,opt,name=since_time" json:"since_time,omitempty"`
	// If set, then only returns the events before this timestamp.
	// If it is a future timestamp, then the event stream will be closed at that moment.
	UntilTime int64 `protobuf:"varint,5,opt,name=until_time" json:"until_time,omitempty"`
}

func (m *EventFilter) Reset()         { *m = EventFilter{} }
func (m *EventFilter) String() string { return proto.CompactTextString(m) }
func (*EventFilter) ProtoMessage()    {}

// Request for GetInfo().
type GetInfoRequest struct {
}

func (m *GetInfoRequest) Reset()         { *m = GetInfoRequest{} }
func (m *GetInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetInfoRequest) ProtoMessage()    {}

// Response for GetInfo().
type GetInfoResponse struct {
	Info *Info `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *GetInfoResponse) Reset()         { *m = GetInfoResponse{} }
func (m *GetInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetInfoResponse) ProtoMessage()    {}

func (m *GetInfoResponse) GetInfo() *Info {
	if m != nil {
		return m.Info
	}
	return nil
}

// Request for ListPods().
type ListPodsRequest struct {
	Filter *PodFilter `protobuf:"bytes,1,opt,name=filter" json:"filter,omitempty"`
}

func (m *ListPodsRequest) Reset()         { *m = ListPodsRequest{} }
func (m *ListPodsRequest) String() string { return proto.CompactTextString(m) }
func (*ListPodsRequest) ProtoMessage()    {}

func (m *ListPodsRequest) GetFilter() *PodFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Response for ListPods().
type ListPodsResponse struct {
	Pods []*Pod `protobuf:"bytes,1,rep,name=pods" json:"pods,omitempty"`
}

func (m *ListPodsResponse) Reset()         { *m = ListPodsResponse{} }
func (m *ListPodsResponse) String() string { return proto.CompactTextString(m) }
func (*ListPodsResponse) ProtoMessage()    {}

func (m *ListPodsResponse) GetPods() []*Pod {
	if m != nil {
		return m.Pods
	}
	return nil
}

// Request for InspectPod().
type InspectPodRequest struct {
	// ID of the pod which we are querying status for, required.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *InspectPodRequest) Reset()         { *m = InspectPodRequest{} }
func (m *InspectPodRequest) String() string { return proto.CompactTextString(m) }
func (*InspectPodRequest) ProtoMessage()    {}

// Response for InspectPod().
type InspectPodResponse struct {
	Pod *Pod `protobuf:"bytes,1,opt,name=pod" json:"pod,omitempty"`
}

func (m *InspectPodResponse) Reset()         { *m = InspectPodResponse{} }
func (m *InspectPodResponse) String() string { return proto.CompactTextString(m) }
func (*InspectPodResponse) ProtoMessage()    {}

func (m *InspectPodResponse) GetPod() *Pod {
	if m != nil {
		return m.Pod
	}
	return nil
}

// Request for ListImages().
type ListImagesRequest struct {
	Filter *ImageFilter `protobuf:"bytes,1,opt,name=filter" json:"filter,omitempty"`
}

func (m *ListImagesRequest) Reset()         { *m = ListImagesRequest{} }
func (m *ListImagesRequest) String() string { return proto.CompactTextString(m) }
func (*ListImagesRequest) ProtoMessage()    {}

func (m *ListImagesRequest) GetFilter() *ImageFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Response for ListImages().
type ListImagesResponse struct {
	Images []*Image `protobuf:"bytes,1,rep,name=images" json:"images,omitempty"`
}

func (m *ListImagesResponse) Reset()         { *m = ListImagesResponse{} }
func (m *ListImagesResponse) String() string { return proto.CompactTextString(m) }
func (*ListImagesResponse) ProtoMessage()    {}

func (m *ListImagesResponse) GetImages() []*Image {
	if m != nil {
		return m.Images
	}
	return nil
}

// Request for InspectImage().
type InspectImageRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *InspectImageRequest) Reset()         { *m = InspectImageRequest{} }
func (m *InspectImageRequest) String() string { return proto.CompactTextString(m) }
func (*InspectImageRequest) ProtoMessage()    {}

// Response for InspectImage().
type InspectImageResponse struct {
	Image *Image `protobuf:"bytes,1,opt,name=image" json:"image,omitempty"`
}

func (m *InspectImageResponse) Reset()         { *m = InspectImageResponse{} }
func (m *InspectImageResponse) String() string { return proto.CompactTextString(m) }
func (*InspectImageResponse) ProtoMessage()    {}

func (m *InspectImageResponse) GetImage() *Image {
	if m != nil {
		return m.Image
	}
	return nil
}

// Request for ListenEvents().
type ListenEventsRequest struct {
	Filter *EventFilter `protobuf:"bytes,1,opt,name=filter" json:"filter,omitempty"`
}

func (m *ListenEventsRequest) Reset()         { *m = ListenEventsRequest{} }
func (m *ListenEventsRequest) String() string { return proto.CompactTextString(m) }
func (*ListenEventsRequest) ProtoMessage()    {}

func (m *ListenEventsRequest) GetFilter() *EventFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Response for ListenEvents().
type ListenEventsResponse struct {
	// Aggregate multiple events to reduce round trips, optional as the response can contain no events.
	Events []*Event `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *ListenEventsResponse) Reset()         { *m = ListenEventsResponse{} }
func (m *ListenEventsResponse) String() string { return proto.CompactTextString(m) }
func (*ListenEventsResponse) ProtoMessage()    {}

func (m *ListenEventsResponse) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

// Request for GetLogs().
type GetLogsRequest struct {
	// ID of the pod which we will get logs from, required.
	PodId string `protobuf:"bytes,1,opt,name=pod_id" json:"pod_id,omitempty"`
	// Name of the app within the pod which we will get logs
	// from, optional. If not set, then the logs of all the
	// apps within the pod will be returned.
	AppName string `protobuf:"bytes,2,opt,name=app_name" json:"app_name,omitempty"`
	// Number of most recent lines to return, optional.
	Lines int32 `protobuf:"varint,3,opt,name=lines" json:"lines,omitempty"`
	// If true, then a response stream will not be closed,
	// and new log response will be sent via the stream, default is false.
	Follow bool `protobuf:"varint,4,opt,name=follow" json:"follow,omitempty"`
	// If set, then only the logs after the timestamp will
	// be returned, optional.
	SinceTime int64 `protobuf:"varint,5,opt,name=since_time" json:"since_time,omitempty"`
	// If set, then only the logs before the timestamp will
	// be returned, optional.
	UntilTime int64 `protobuf:"varint,6,opt,name=until_time" json:"until_time,omitempty"`
}

func (m *GetLogsRequest) Reset()         { *m = GetLogsRequest{} }
func (m *GetLogsRequest) String() string { return proto.CompactTextString(m) }
func (*GetLogsRequest) ProtoMessage()    {}

// Response for GetLogs().
type GetLogsResponse struct {
	// List of the log lines that returned, optional as the response can contain no logs.
	Lines []string `protobuf:"bytes,1,rep,name=lines" json:"lines,omitempty"`
}

func (m *GetLogsResponse) Reset()         { *m = GetLogsResponse{} }
func (m *GetLogsResponse) String() string { return proto.CompactTextString(m) }
func (*GetLogsResponse) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("v1alpha.ImageType", ImageType_name, ImageType_value)
	proto.RegisterEnum("v1alpha.AppState", AppState_name, AppState_value)
	proto.RegisterEnum("v1alpha.PodState", PodState_name, PodState_value)
	proto.RegisterEnum("v1alpha.EventType", EventType_name, EventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for PublicAPI service

type PublicAPIClient interface {
	// GetInfo gets the rkt's information on the machine.
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
	// ListPods lists rkt pods on the machine.
	ListPods(ctx context.Context, in *ListPodsRequest, opts ...grpc.CallOption) (*ListPodsResponse, error)
	// InspectPod gets detailed pod information of the specified pod.
	InspectPod(ctx context.Context, in *InspectPodRequest, opts ...grpc.CallOption) (*InspectPodResponse, error)
	// ListImages lists the images on the machine.
	ListImages(ctx context.Context, in *ListImagesRequest, opts ...grpc.CallOption) (*ListImagesResponse, error)
	// InspectImage gets the detailed image information of the specified image.
	InspectImage(ctx context.Context, in *InspectImageRequest, opts ...grpc.CallOption) (*InspectImageResponse, error)
	// ListenEvents listens for the events, it will return a response stream
	// that will contain event objects.
	ListenEvents(ctx context.Context, in *ListenEventsRequest, opts ...grpc.CallOption) (PublicAPI_ListenEventsClient, error)
	// GetLogs gets the logs for a pod, if the app is also specified, then only the logs
	// of the app will be returned.
	//
	// If 'follow' in the 'GetLogsRequest' is set to 'true', then the response stream
	// will not be closed after the first response, the future logs will be sent via
	// the stream.
	GetLogs(ctx context.Context, in *GetLogsRequest, opts ...grpc.CallOption) (PublicAPI_GetLogsClient, error)
}

type publicAPIClient struct {
	cc *grpc.ClientConn
}

func NewPublicAPIClient(cc *grpc.ClientConn) PublicAPIClient {
	return &publicAPIClient{cc}
}

func (c *publicAPIClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := grpc.Invoke(ctx, "/v1alpha.PublicAPI/GetInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicAPIClient) ListPods(ctx context.Context, in *ListPodsRequest, opts ...grpc.CallOption) (*ListPodsResponse, error) {
	out := new(ListPodsResponse)
	err := grpc.Invoke(ctx, "/v1alpha.PublicAPI/ListPods", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicAPIClient) InspectPod(ctx context.Context, in *InspectPodRequest, opts ...grpc.CallOption) (*InspectPodResponse, error) {
	out := new(InspectPodResponse)
	err := grpc.Invoke(ctx, "/v1alpha.PublicAPI/InspectPod", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicAPIClient) ListImages(ctx context.Context, in *ListImagesRequest, opts ...grpc.CallOption) (*ListImagesResponse, error) {
	out := new(ListImagesResponse)
	err := grpc.Invoke(ctx, "/v1alpha.PublicAPI/ListImages", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicAPIClient) InspectImage(ctx context.Context, in *InspectImageRequest, opts ...grpc.CallOption) (*InspectImageResponse, error) {
	out := new(InspectImageResponse)
	err := grpc.Invoke(ctx, "/v1alpha.PublicAPI/InspectImage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicAPIClient) ListenEvents(ctx context.Context, in *ListenEventsRequest, opts ...grpc.CallOption) (PublicAPI_ListenEventsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_PublicAPI_serviceDesc.Streams[0], c.cc, "/v1alpha.PublicAPI/ListenEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicAPIListenEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PublicAPI_ListenEventsClient interface {
	Recv() (*ListenEventsResponse, error)
	grpc.ClientStream
}

type publicAPIListenEventsClient struct {
	grpc.ClientStream
}

func (x *publicAPIListenEventsClient) Recv() (*ListenEventsResponse, error) {
	m := new(ListenEventsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *publicAPIClient) GetLogs(ctx context.Context, in *GetLogsRequest, opts ...grpc.CallOption) (PublicAPI_GetLogsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_PublicAPI_serviceDesc.Streams[1], c.cc, "/v1alpha.PublicAPI/GetLogs", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicAPIGetLogsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PublicAPI_GetLogsClient interface {
	Recv() (*GetLogsResponse, error)
	grpc.ClientStream
}

type publicAPIGetLogsClient struct {
	grpc.ClientStream
}

func (x *publicAPIGetLogsClient) Recv() (*GetLogsResponse, error) {
	m := new(GetLogsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for PublicAPI service

type PublicAPIServer interface {
	// GetInfo gets the rkt's information on the machine.
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	// ListPods lists rkt pods on the machine.
	ListPods(context.Context, *ListPodsRequest) (*ListPodsResponse, error)
	// InspectPod gets detailed pod information of the specified pod.
	InspectPod(context.Context, *InspectPodRequest) (*InspectPodResponse, error)
	// ListImages lists the images on the machine.
	ListImages(context.Context, *ListImagesRequest) (*ListImagesResponse, error)
	// InspectImage gets the detailed image information of the specified image.
	InspectImage(context.Context, *InspectImageRequest) (*InspectImageResponse, error)
	// ListenEvents listens for the events, it will return a response stream
	// that will contain event objects.
	ListenEvents(*ListenEventsRequest, PublicAPI_ListenEventsServer) error
	// GetLogs gets the logs for a pod, if the app is also specified, then only the logs
	// of the app will be returned.
	//
	// If 'follow' in the 'GetLogsRequest' is set to 'true', then the response stream
	// will not be closed after the first response, the future logs will be sent via
	// the stream.
	GetLogs(*GetLogsRequest, PublicAPI_GetLogsServer) error
}

func RegisterPublicAPIServer(s *grpc.Server, srv PublicAPIServer) {
	s.RegisterService(&_PublicAPI_serviceDesc, srv)
}

func _PublicAPI_GetInfo_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PublicAPIServer).GetInfo(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _PublicAPI_ListPods_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ListPodsRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PublicAPIServer).ListPods(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _PublicAPI_InspectPod_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(InspectPodRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PublicAPIServer).InspectPod(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _PublicAPI_ListImages_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ListImagesRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PublicAPIServer).ListImages(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _PublicAPI_InspectImage_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(InspectImageRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(PublicAPIServer).InspectImage(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _PublicAPI_ListenEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PublicAPIServer).ListenEvents(m, &publicAPIListenEventsServer{stream})
}

type PublicAPI_ListenEventsServer interface {
	Send(*ListenEventsResponse) error
	grpc.ServerStream
}

type publicAPIListenEventsServer struct {
	grpc.ServerStream
}

func (x *publicAPIListenEventsServer) Send(m *ListenEventsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PublicAPI_GetLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetLogsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PublicAPIServer).GetLogs(m, &publicAPIGetLogsServer{stream})
}

type PublicAPI_GetLogsServer interface {
	Send(*GetLogsResponse) error
	grpc.ServerStream
}

type publicAPIGetLogsServer struct {
	grpc.ServerStream
}

func (x *publicAPIGetLogsServer) Send(m *GetLogsResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _PublicAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1alpha.PublicAPI",
	HandlerType: (*PublicAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _PublicAPI_GetInfo_Handler,
		},
		{
			MethodName: "ListPods",
			Handler:    _PublicAPI_ListPods_Handler,
		},
		{
			MethodName: "InspectPod",
			Handler:    _PublicAPI_InspectPod_Handler,
		},
		{
			MethodName: "ListImages",
			Handler:    _PublicAPI_ListImages_Handler,
		},
		{
			MethodName: "InspectImage",
			Handler:    _PublicAPI_InspectImage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenEvents",
			Handler:       _PublicAPI_ListenEvents_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetLogs",
			Handler:       _PublicAPI_GetLogs_Handler,
			ServerStreams: true,
		},
	},
}
