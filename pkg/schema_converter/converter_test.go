package converter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Rdata struct {
	data map[string]interface{}
}

func (r *Rdata) Get(key string) interface{} {
	return r.data[key]
}

func (r *Rdata) Set(key string, value interface{}) error {
	r.data[key] = value
	return nil
}

type Deployment struct {
	Id          string
	Node        uint32
	Disks       []Disk
	ZDBs        []ZDB
	VMs         []VM
	UsedIPs     []string
	NetworkName string
}

type ZDB struct {
	Name        string
	Password    string
	Public      bool
	Size        int
	Description string
	Mode        string
	Ips         []string
	Port        uint32
	Namespace   string
}

type Disk struct {
	Name        string
	Size        int
	Description string
}

type VM struct {
	Name          string
	Flist         string
	FlistChecksum string
	Publicip      bool
	Publicip6     bool
	Planetary     bool
	Corex         bool
	Computedip    string
	Computedip6   string
	YggIp         string
	Ip            string
	Description   string
	Cpu           int
	Memory        int
	RootfsSize    int
	Entrypoint    string
	Mounts        []Mount
	Zlogs         []Zlog
	EnvVars       map[string]string
}
type Mount struct {
	DiskName   string
	MountPoint string
}

type Zlog struct {
	Output string
}

func getDeployment() Deployment {
	dp := Deployment{}
	dp.Id = "1234"
	dp.Node = 123456
	dp.Disks = []Disk{{"d1", 5, "desc1"}, {"d2", 6, "desc2"}}
	dp.NetworkName = "net1"
	dp.UsedIPs = []string{"usedip1", "usedip2"}
	dp.VMs = []VM{{
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
		map[string]string{"env1": "var1", "env2": "var2"},
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
			map[string]string{"env3": "var3", "env4": "var4"},
		}}
	dp.ZDBs = []ZDB{{
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
		}}
	return dp
}

func getFilledRD() Rdata {
	rd := Rdata{}
	rd.data = map[string]interface{}{}
	rd.Set("id", "1234")
	rd.Set("node", uint32(123456))
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
	rd.Set("zdbs", []interface{}{
		map[string]interface{}{
			"name":        "zdb1",
			"password":    "pass1",
			"public":      true,
			"size":        int(5),
			"description": "desc1",
			"mode":        "mod1",
			"ips":         []interface{}{"ip1, ip2"},
			"port":        uint32(1234),
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
			"port":        uint32(5678),
			"namespace":   "namespace2",
		},
	})
	rd.Set("vms", []interface{}{
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
			"zlogs":    []interface{}{map[string]interface{}{"output": "zlog1"}, map[string]interface{}{"output": "zlog2"}},
			"env_vars": map[string]interface{}{"env1": "var1", "env2": "var2"},
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
			"zlogs":    []interface{}{map[string]interface{}{"output": "zlog3"}, map[string]interface{}{"output": "zlog4"}},
			"env_vars": map[string]interface{}{"env3": "var3", "env4": "var4"},
		},
	})
	rd.Set("used_ips", []interface{}{"usedip1", "usedip2"})
	rd.Set("network_name", "net1")
	return rd
}

func getEmptyRD() Rdata {
	rd := Rdata{}
	rd.data = map[string]interface{}{}
	rd.Set("id", "")
	rd.Set("node", uint32(0))
	rd.Set("disks", []map[string]interface{}{})
	rd.Set("zdbs", []map[string]interface{}{})
	rd.Set("vms", []map[string]interface{}{})
	rd.Set("used_ips", []string{})
	rd.Set("network_name", "")
	return rd
}

func TestConverter(t *testing.T) {
	dp := getDeployment()
	rd := getFilledRD()

	newDP := Deployment{}
	newRD := getEmptyRD()

	Encode(dp, &newRD)
	assert.Equal(t, true, reflect.DeepEqual(rd, newRD))

	Decode(&newDP, &rd)
	assert.Equal(t, true, reflect.DeepEqual(newDP, dp))

}
