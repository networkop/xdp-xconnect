# Cross-connect Linux interfaces with XDP redirect


[![Go Report Card](https://goreportcard.com/badge/github.com/networkop/xdp-xconnect)](https://goreportcard.com/report/github.com/networkop/xdp-xconnect)
[![GoDoc](https://godoc.org/istio.io/istio?status.svg)](https://pkg.go.dev/github.com/networkop/xdp-xconnect)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


![](icon.png)

`xdp-xconnect` daemon is a long-running process that uses a YAML file as its configuration API. For example:

```yaml
links:
    eth0: tap0
    veth2: veth3
```

Given the above YAML file, local Linux interfaces will be cross-connected (`eth0<->tap0` and `veth2<->veth3`) with the following command:

```
sudo xdp-xconnect -conf config.yaml
```

This command will block, listening to any changes to the file and will perform "warm" reconfigurations on the fly.

> Note: due to its nature (loading eBPF progs, maps and interacting with netlink), `xdp-xconnect` requires `NET_ADMIN` capabilities (root privileges are used for simplicity).

## Theory

Each pair of interfaces will have an eBPF program attached to its XDP hook and will use `bpf_redirect_map` eBPF [helper function](https://man7.org/linux/man-pages/man7/bpf-helpers.7.html) to redirect packets directly to the receive queue of the other interface. For mode details read [this](https://github.com/xdp-project/xdp-tutorial/tree/master/packet03-redirecting).

## Installation

Binary

```
go get github.com/networkop/xdp-xconnect
```

Docker:

```
docker pull networkop/xdp-xconnect
```

## Usage

Due to its nature, `xdp-xconnect` must be run with root privileges:

```
sudo xdp-xconnect -conf input.yaml
```


```
docker run --privileged networkop/xdp-xconnect -conf input.yaml
```

```go
import "github.com/networkop/xdp-xconnect/pkg/xdp"

func main() {

    input := map[string]string{"eth1":"tap1"}

    app, err := xdp.NewXconnectApp(input)
	// handle error

    updateCh := make(chan map[string]string, 1)

	app.Launch(ctx, updateCh)
}
```

Test setup (move into integration testing)

```
sudo ip link del dev xconnect-1
sudo ip link del dev xconnect-2
sudo ip link del dev xconnect-3
sudo ip netns del ns1
sudo ip netns del ns2
sudo ip netns del ns3
cp testdata/first.yaml testdata/input.yaml

sudo ip link add dev xconnect-1 type veth peer name xc-1
sudo ip link add dev xconnect-2 type veth peer name xc-2
sudo ip link add dev xconnect-3 type veth peer name xc-3

sudo ip link set dev xconnect-1 up
sudo ip link set dev xconnect-2 up
sudo ip link set dev xconnect-3 up

sudo ip netns add ns1
sudo ip netns add ns2
sudo ip netns add ns3
sudo ip link set xc-1 netns ns1
sudo ip link set xc-2 netns ns2
sudo ip link set xc-3 netns ns3
sudo ip netns exec ns1 ip addr add 169.254.1.10/24 dev xc-1
sudo ip netns exec ns2 ip addr add 169.254.1.20/24 dev xc-2
sudo ip netns exec ns3 ip addr add 169.254.1.30/24 dev xc-3
sudo ip netns exec ns1 ip link set dev xc-1 up
sudo ip netns exec ns2 ip link set dev xc-2 up
sudo ip netns exec ns3 ip link set dev xc-3 up

sudo go run ./main.go -conf testdata/input.yaml &


sudo ip netns exec ns1 ping -c1 -w1 169.254.1.20

cp testdata/second.yaml testdata/input.yaml
sudo ip netns exec ns1 ping -c1 -w1 169.254.1.20

cp testdata/third.yaml testdata/input.yaml
sudo ip netns exec ns1 ping -c1 -w1 169.254.1.30

cp testdata/bad.yaml testdata/input.yaml
```

REQS:
Linux > 4.14 (veth XDP support)


Reading:


https://github.com/takehaya/goxdp-template

https://github.com/hrntknr/nfNat

https://github.com/takehaya/Vinbero

https://github.com/tcfw/vpc

https://github.com/florianl/tc-skeleton

https://github.com/cloudflare/rakelimit

https://github.com/b3a-dev/ebpf-geoip-demo

https://github.com/lmb/ship-bpf-with-go

https://qmonnet.github.io/whirl-offload/2020/04/12/llvm-ebpf-asm/

https://docs.cilium.io/en/stable/bpf/

https://github.com/xdp-project/xdp-tutorial