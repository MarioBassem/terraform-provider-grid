package converter

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DeploymentDeployer struct {
	Id               string
	Node             int16
	Disks            []Disk
	ZDBList          []ZDB
	VMList           []VM
	QSFSList         []QSFS
	IPRange          string
	NetworkName      string
	FQDN             []GatewayFQDNProxy
	GatewayNames     []GatewayNameProxy
	NodeDeploymentID map[uint32]uint64
}

type VM struct {
	Name          string
	Flist         string
	FlistChecksum string
	Publicip      bool //not PublicIP
	Publicip6     bool //not PublicIP6
	Planetary     bool
	Corex         bool
	Computedip    string //not ComputedIP
	Computedip6   string //moded
	YggIP         string
	IP            string
	Description   string
	Cpu           int
	Memory        int
	RootfsSize    int
	Entrypoint    string
	Mounts        []Mount
	Zlogs         []Zlog
	EnvVars       map[string]string

	NetworkName string
}

type Mount struct {
	DiskName   string
	MountPoint string
}

type Zlog struct {
	Output string
}

type Disk struct {
	Name        string
	Size        int
	Description string
}
type ZDB struct {
	Name        string
	Password    string
	Public      bool
	Size        int
	Description string
	Mode        string
	Ips         []string //Ips not IPs
	Port        uint32
	Namespace   string
}

type QSFS struct {
	Name                 string
	Description          string
	Cache                int
	MinimalShards        uint32
	ExpectedShards       uint32
	RedundantGroups      uint32
	RedundantNodes       uint32
	MaxZDBDataDirSize    uint32
	EncryptionAlgorithm  string
	EncryptionKey        string
	CompressionAlgorithm string
	Metadata             Metadata
	Groups               Groups

	MetricsEndpoint string
}
type Metadata struct {
	Type                string
	Prefix              string
	EncryptionAlgorithm string
	EncryptionKey       string
	Backends            Backends
}
type Group struct {
	Backends Backends
}
type Backend ZdbBackend
type Groups []Group
type Backends []Backend

type ZdbBackend struct {
	Address   string `json:"address" toml:"address"`
	Namespace string `json:"namespace" toml:"namespace"`
	Password  string `json:"password" toml:"password"`
}

type GatewayNameProxy struct {
	// Name the fully qualified domain name to use (cannot be present with Name)
	Name string

	// Passthrough whether to pass tls traffic or not
	TLSPassthrough bool

	// Backends are list of backend ips
	Backends []Backend

	// FQDN deployed on the node
	FQDN string
}

type GatewayFQDNProxy struct {
	// Name the fully qualified domain name to use (cannot be present with Name)
	Name string

	// Passthrough whether to pass tls traffic or not
	TLSPassthrough bool

	// Backends are list of backend ips
	Backends []Backend

	// FQDN deployed on the node
	FQDN string
}

func getDeployment() DeploymentDeployer {
	dp := DeploymentDeployer{}
	dp.Id = "1234"
	dp.Node = int16(1)
	dp.Disks = []Disk{{"d1", 5, "desc1"}, {"d2", 6, "desc2"}}
	dp.ZDBList = []ZDB{{
		"zdb1",
		"pass1",
		true,
		5,
		"desc1",
		"mod1",
		[]string{"ip1, ip2"},
		1234,
		"namespace1",
	},
		{
			"zdb2",
			"pass2",
			true,
			5,
			"desc2",
			"mod2",
			[]string{"ip3, ip4"},
			5678,
			"namespace2",
		},
	}

	dp.VMList = []VM{{
		"vm1",
		"flist1",
		"flist_checksum1",
		false,
		false,
		true,
		false,
		"computedip_1",
		"computedip6_1",
		"yggip1",
		"ip1",
		"desc1",
		2,
		5,
		3,
		"entrypoint1",
		[]Mount{{"d1", "mp1"}, {"d2", "mp2"}},
		[]Zlog{{"zlog1"}, {"zlog2"}},
		map[string]string{"1": "var1", "2": "var2"},
		"net1",
	},
		{
			"vm2",
			"flist2",
			"flist_checksum2",
			true,
			true,
			false,
			true,
			"computedip_2",
			"computedip6_2",
			"yggip2",
			"ip2",
			"desc2",
			5,
			7,
			4,
			"entrypoint2",
			[]Mount{{"d5", "mp5"}, {"d6", "mp6"}},
			[]Zlog{{"zlog3"}, {"zlog4"}},
			map[string]string{"3": "var3", "4": "var4"},
			"net2",
		},
	}
	dp.QSFSList = []QSFS{
		{
			"name1",
			"desc1",
			1,
			2,
			3,
			4,
			5,
			6,
			"encalgo",
			"key1",
			"comalgo",
			Metadata{
				"tp1",
				"pre1",
				"encalgo",
				"key1",
				Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			},
			Groups{
				{
					Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
				},
				{
					Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
				},
			},
			"endpoint1",
		},
		{
			"name2",
			"desc2",
			1,
			2,
			3,
			4,
			5,
			6,
			"encalgo",
			"key1",
			"comalgo",
			Metadata{
				"tp1",
				"pre1",
				"encalgo",
				"key1",
				Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			},
			Groups{
				{
					Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
				},
				{
					Backends{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
				},
			},
			"endpoint2",
		},
	}
	dp.IPRange = "iprange"
	dp.NetworkName = "net1"
	dp.FQDN = []GatewayFQDNProxy{
		{
			"name1",
			true,
			[]Backend{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			"fqdn1",
		},
		{
			"name2",
			true,
			[]Backend{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			"fqdn2",
		},
	}
	dp.GatewayNames = []GatewayNameProxy{
		{
			"name1",
			false,
			[]Backend{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			"fqdn1",
		},
		{
			"name2",
			false,
			[]Backend{{"add3", "ns3", "pss3"}, {"add4", "ns4", "pss4"}},
			"fqdn2",
		},
	}
	dp.NodeDeploymentID = map[uint32]uint64{1: 2}

	return dp
}

type ResourceData struct {
	data map[string]interface{}
}

func (r *ResourceData) Get(key string) interface{} {
	return r.data[key]
}

func (r *ResourceData) Set(key string, value interface{}) error {
	r.data[key] = value
	return nil
}

func getFilledResourceData() ResourceData {
	rd := ResourceData{}
	rd.data = map[string]interface{}{}
	rd.Set("id", "1234")
	rd.Set("node", 1)
	rd.Set("disks", []interface{}{
		map[string]interface{}{
			"name":        "d1",
			"size":        int(5),
			"description": "desc1",
		}, map[string]interface{}{
			"name":        "d2",
			"size":        int(6),
			"description": "desc2",
		}})
	rd.Set("zdb_list", []interface{}{
		map[string]interface{}{
			"name":        "zdb1",
			"password":    "pass1",
			"public":      true,
			"size":        int(5),
			"description": "desc1",
			"mode":        "mod1",
			"ips":         []interface{}{"ip1, ip2"},
			"port":        int(1234),
			"namespace":   "namespace1",
		},
		map[string]interface{}{
			"name":        "zdb2",
			"password":    "pass2",
			"public":      true,
			"size":        int(5),
			"description": "desc2",
			"mode":        "mod2",
			"ips":         []interface{}{"ip3, ip4"},
			"port":        int(5678),
			"namespace":   "namespace2",
		},
	})
	rd.Set("vm_list", []interface{}{
		map[string]interface{}{
			"name":           "vm1",
			"flist":          "flist1",
			"flist_checksum": "flist_checksum1",
			"publicip":       false,
			"publicip6":      false,
			"planetary":      true,
			"corex":          false,
			"computedip":     "computedip_1",
			"computedip6":    "computedip6_1",
			"ygg_ip":         "yggip1",
			"ip":             "ip1",
			"description":    "desc1",
			"cpu":            int(2),
			"memory":         int(5),
			"rootfs_size":    int(3),
			"entrypoint":     "entrypoint1",
			"mounts": []interface{}{
				map[string]interface{}{
					"disk_name":   "d1",
					"mount_point": "mp1",
				},
				map[string]interface{}{
					"disk_name":   "d2",
					"mount_point": "mp2",
				}},
			"zlogs":        []interface{}{map[string]interface{}{"output": "zlog1"}, map[string]interface{}{"output": "zlog2"}},
			"env_vars":     map[string]interface{}{"1": "var1", "2": "var2"},
			"network_name": "net1",
		},
		map[string]interface{}{
			"name":           "vm2",
			"flist":          "flist2",
			"flist_checksum": "flist_checksum2",
			"publicip":       true,
			"publicip6":      true,
			"planetary":      false,
			"corex":          true,
			"computedip":     "computedip_2",
			"computedip6":    "computedip6_2",
			"ygg_ip":         "yggip2",
			"ip":             "ip2",
			"description":    "desc2",
			"cpu":            int(5),
			"memory":         int(7),
			"rootfs_size":    int(4),
			"entrypoint":     "entrypoint2",
			"mounts": []interface{}{
				map[string]interface{}{
					"disk_name":   "d5",
					"mount_point": "mp5",
				},
				map[string]interface{}{
					"disk_name":   "d6",
					"mount_point": "mp6",
				}},
			"zlogs":        []interface{}{map[string]interface{}{"output": "zlog3"}, map[string]interface{}{"output": "zlog4"}},
			"env_vars":     map[string]interface{}{"3": "var3", "4": "var4"},
			"network_name": "net2",
		},
	})
	rd.Set("qsfs_list", []interface{}{
		map[string]interface{}{
			"name":                  "name1",
			"description":           "desc1",
			"cache":                 int(1),
			"minimal_shards":        int(2),
			"expected_shards":       int(3),
			"redundant_groups":      int(4),
			"redundant_nodes":       int(5),
			"max_zdb_data_dir_size": int(6),
			"encryption_algorithm":  "encalgo",
			"encryption_key":        "key1",
			"compression_algorithm": "comalgo",
			"metadata": map[string]interface{}{
				"type":                 "tp1",
				"prefix":               "pre1",
				"encryption_algorithm": "encalgo",
				"encryption_key":       "key1",
				"backends": []interface{}{
					map[string]interface{}{
						"address":   "add3",
						"namespace": "ns3",
						"password":  "pss3",
					},
					map[string]interface{}{
						"address":   "add4",
						"namespace": "ns4",
						"password":  "pss4",
					},
				},
			},
			"groups": []interface{}{
				map[string]interface{}{
					"backends": []interface{}{
						map[string]interface{}{
							"address":   "add3",
							"namespace": "ns3",
							"password":  "pss3",
						},
						map[string]interface{}{
							"address":   "add4",
							"namespace": "ns4",
							"password":  "pss4",
						},
					},
				},
				map[string]interface{}{
					"backends": []interface{}{
						map[string]interface{}{
							"address":   "add3",
							"namespace": "ns3",
							"password":  "pss3",
						},
						map[string]interface{}{
							"address":   "add4",
							"namespace": "ns4",
							"password":  "pss4",
						},
					},
				},
			},
			"metrics_endpoint": "endpoint1",
		},
		map[string]interface{}{
			"name":                  "name2",
			"description":           "desc2",
			"cache":                 int(1),
			"minimal_shards":        int(2),
			"expected_shards":       int(3),
			"redundant_groups":      int(4),
			"redundant_nodes":       int(5),
			"max_zdb_data_dir_size": int(6),
			"encryption_algorithm":  "encalgo",
			"encryption_key":        "key1",
			"compression_algorithm": "comalgo",
			"metadata": map[string]interface{}{
				"type":                 "tp1",
				"prefix":               "pre1",
				"encryption_algorithm": "encalgo",
				"encryption_key":       "key1",
				"backends": []interface{}{
					map[string]interface{}{
						"address":   "add3",
						"namespace": "ns3",
						"password":  "pss3",
					},
					map[string]interface{}{
						"address":   "add4",
						"namespace": "ns4",
						"password":  "pss4",
					},
				},
			},
			"groups": []interface{}{
				map[string]interface{}{
					"backends": []interface{}{
						map[string]interface{}{
							"address":   "add3",
							"namespace": "ns3",
							"password":  "pss3",
						},
						map[string]interface{}{
							"address":   "add4",
							"namespace": "ns4",
							"password":  "pss4",
						},
					},
				},
				map[string]interface{}{
					"backends": []interface{}{
						map[string]interface{}{
							"address":   "add3",
							"namespace": "ns3",
							"password":  "pss3",
						},
						map[string]interface{}{
							"address":   "add4",
							"namespace": "ns4",
							"password":  "pss4",
						},
					},
				},
			},
			"metrics_endpoint": "endpoint2",
		},
	})
	rd.Set("ip_range", "iprange")
	rd.Set("network_name", "net1")
	rd.Set("fqdn", []interface{}{
		map[string]interface{}{
			"name":            "name1",
			"tls_passthrough": true,
			"backends": []interface{}{
				map[string]interface{}{
					"address":   "add3",
					"namespace": "ns3",
					"password":  "pss3",
				},
				map[string]interface{}{
					"address":   "add4",
					"namespace": "ns4",
					"password":  "pss4",
				},
			},
			"fqdn": "fqdn1",
		},
		map[string]interface{}{
			"name":            "name2",
			"tls_passthrough": true,
			"backends": []interface{}{
				map[string]interface{}{
					"address":   "add3",
					"namespace": "ns3",
					"password":  "pss3",
				},
				map[string]interface{}{
					"address":   "add4",
					"namespace": "ns4",
					"password":  "pss4",
				},
			},
			"fqdn": "fqdn2",
		},
	})
	rd.Set("gateway_names", []interface{}{
		map[string]interface{}{
			"name":            "name1",
			"tls_passthrough": false,
			"backends": []interface{}{
				map[string]interface{}{
					"address":   "add3",
					"namespace": "ns3",
					"password":  "pss3",
				},
				map[string]interface{}{
					"address":   "add4",
					"namespace": "ns4",
					"password":  "pss4",
				},
			},
			"fqdn": "fqdn1",
		},
		map[string]interface{}{
			"name":            "name2",
			"tls_passthrough": false,
			"backends": []interface{}{
				map[string]interface{}{
					"address":   "add3",
					"namespace": "ns3",
					"password":  "pss3",
				},
				map[string]interface{}{
					"address":   "add4",
					"namespace": "ns4",
					"password":  "pss4",
				},
			},
			"fqdn": "fqdn2",
		},
	})
	rd.Set("node_deployment_id", map[string]interface{}{
		"1": 2,
	})
	return rd
}

func getEmptyResourceDeployment() ResourceData {
	rd := ResourceData{}
	rd.data = map[string]interface{}{}
	rd.Set("id", "")
	rd.Set("node", 0)
	rd.Set("disks", []interface{}{})
	rd.Set("zdb_list", []interface{}{})
	rd.Set("vm_list", []interface{}{})
	rd.Set("qsfs_list", []interface{}{})
	rd.Set("ip_range", "")
	rd.Set("network_name", "")
	rd.Set("fqdn", []interface{}{})
	rd.Set("gateway_names", "")
	rd.Set("node_deployment_id", map[string]interface{}{})
	return rd
}

func TestConverter(t *testing.T) {

	dp := getDeployment()
	rd := getFilledResourceData()

	newDP := DeploymentDeployer{}
	newRD := getEmptyResourceDeployment()

	err := Encode(dp, &newRD)
	if err != nil {
		log.Printf("error in encoding: %+v", err)
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, rd, newRD)

	err2 := Decode(&newDP, &rd)
	if err2 != nil {
		log.Printf("error in decoding: %+v", err2)
		assert.Equal(t, nil, err2)
	}
	assert.Equal(t, dp, newDP)

}
