package server

import (
	"context"
	"github.com/Kong/kuma/pkg/core/xds"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	go_cp_server "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/pkg/errors"
	"sync"
)

type DataplaneMetadataTracker struct {
	mutex            sync.RWMutex
	metadataForProxy map[xds.ProxyId]*xds.DataplaneMetadata
	streams          map[int64]xds.ProxyId
}

func NewDataplaneMetadataTracker() *DataplaneMetadataTracker {
	return &DataplaneMetadataTracker{
		mutex:            sync.RWMutex{},
		metadataForProxy: map[xds.ProxyId]*xds.DataplaneMetadata{},
		streams:          map[int64]xds.ProxyId{},
	}
}

func (d *DataplaneMetadataTracker) Metadata(proxyId xds.ProxyId) *xds.DataplaneMetadata {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	metadata, found := d.metadataForProxy[proxyId]
	if found {
		return metadata
	} else {
		return &xds.DataplaneMetadata{}
	}
}

var _ go_cp_server.Callbacks = &DataplaneMetadataTracker{}

func (d *DataplaneMetadataTracker) OnStreamOpen(context.Context, int64, string) error {
	return nil
}

func (d *DataplaneMetadataTracker) OnStreamClosed(stream int64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	proxyId, found := d.streams[stream]
	if found {
		delete(d.streams, stream)
		delete(d.metadataForProxy, proxyId)
	}
}

func (d *DataplaneMetadataTracker) OnStreamRequest(stream int64, req *v2.DiscoveryRequest) error {
	proxyId, err := xds.ParseProxyId(req.Node)
	if err != nil {
		return errors.Wrapf(err, "could not parse proxy id of %s", req.Node.Id)
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.streams[stream] = *proxyId
	d.metadataForProxy[*proxyId] = xds.DataplaneMetadataFromNode(req.Node)
	return nil
}

func (d *DataplaneMetadataTracker) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {
}

func (d *DataplaneMetadataTracker) OnFetchRequest(context.Context, *v2.DiscoveryRequest) error {
	return nil
}

func (d *DataplaneMetadataTracker) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {
}
