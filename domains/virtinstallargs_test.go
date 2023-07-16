package domains_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"git.voidnet.tech/kev/easysandbox/domains"
)

func TestGetVirtInstallArgs(t *testing.T) {
	type args struct {
		domainName string
		args       []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "No optional args",
			args: args{
				domainName: "testdomain",
				args:       nil,
			},
			want: []string{
				"--name", "testdomain",
				"--disk", domains.DomainsDir + "testdomain" + "/" + "root.qcow2" + ",target.bus=sata",
				"--disk", domains.DomainsDir + "testdomain" + "/" + "home.qcow2" + ",target.bus=sata",
				"--import", "--hvm", "--network", "bridge=virbr0", "--virt-type", "kvm", "--install", "no_install=yes", "--noreboot",
				"--memory", "4096", "--os-variant", "linux2022", "--vcpus", fmt.Sprintf("%d", runtime.NumCPU()),
			},
		},
		{
			name: "With optional args",
			args: args{
				domainName: "testdomain",
				args:       []string{"--vcpus", "2", "--memory", "2048"},
			},
			want: []string{
				"--name", "testdomain",
				"--disk", domains.DomainsDir + "testdomain" + "/" + "root.qcow2" + ",target.bus=sata",
				"--disk", domains.DomainsDir + "testdomain" + "/" + "home.qcow2" + ",target.bus=sata",
				"--import", "--hvm", "--network", "bridge=virbr0", "--virt-type", "kvm", "--install", "no_install=yes", "--noreboot",
				"--vcpus", "2", "--memory", "2048", "--os-variant", "linux2022",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domains.GetVirtInstallArgsString(tt.args.domainName, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVirtInstallArgs() = %v, want \n %v", got, tt.want)
			}
		})
	}
}
