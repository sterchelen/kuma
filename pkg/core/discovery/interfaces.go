package discovery

import (
	discovery_proto "github.com/Kong/kuma/api/discovery/v1alpha1"
	mesh_core "github.com/Kong/kuma/pkg/core/resources/apis/mesh"
	"github.com/Kong/kuma/pkg/core/resources/model"
)

// DiscoverySource is a source of discovery information, i.e. Services and Workloads.
type DiscoverySource interface {
	AddConsumer(DiscoveryConsumer)
}

// ServiceInfo combines original proto object with auxiliary information.
type ServiceInfo struct {
	Service *discovery_proto.Service
}

// WorkloadInfo combines original proto object with auxiliary information.
type WorkloadInfo struct {
	Workload *discovery_proto.Workload
	Desc     *WorkloadDescription
}

// WorkloadDescription represents auxiliary information about a Workload.
type WorkloadDescription struct {
	Version   string
	Endpoints []WorkloadEndpoint
}

type WorkloadEndpoint struct {
	Address string
	Port    uint32
}

// DiscoveryConsumer is a consumer of discovery information, i.e. Services and Workloads.
type DiscoveryConsumer interface {
	DataplaneDiscoveryConsumer
}

type DataplaneDiscoveryConsumer interface {
	OnDataplaneUpdate(*mesh_core.DataplaneResource) error
	OnDataplaneDelete(model.ResourceKey) error
}
