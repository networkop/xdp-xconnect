# Cross-connect Linux interfaces with XDP redirect


[![Go Report Card](https://goreportcard.com/badge/github.com/networkop/xdp-xconnect)](https://goreportcard.com/report/github.com/networkop/xdp-xconnect)
[![GoDoc](https://godoc.org/istio.io/istio?status.svg)](https://pkg.go.dev/github.com/networkop/xdp-xconnect)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Build Status](https://github.com/networkop/xdp-xconnect/actions/workflows/ci.yml/badge.svg)

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

Each pair of interfaces will have an eBPF program attached to its XDP hook and will use `bpf_redirect_map` eBPF [helper function](https://man7.org/linux/man-pages/man7/bpf-helpers.7.html) to redirect packets directly to the receive queue of the peer interface. For mode details read [this](https://github.com/xdp-project/xdp-tutorial/tree/master/packet03-redirecting).


## Prerequisites

Linux Kernel version > 4.14 (introduced veth XDP support)
Go

## Installation

Binary:

```
go install github.com/networkop/xdp-xconnect@latest
```

Docker:

```
docker pull networkop/xdp-xconnect
```

## Usage

Binary:

```
sudo xdp-xconnect -conf input.yaml
```

Docker:

```
docker run --net host -v$(pwd):/xc --privileged networkop/xdp-xconnect -conf /xc/input.yaml
```

Go code:


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

## Demo

Create three network namespaces and three veth links:


```
sudo ip link add dev xconnect-1 type veth peer name xc-1
sudo ip link add dev xconnect-2 type veth peer name xc-2
sudo ip link add dev xconnect-3 type veth peer name xc-3
sudo ip netns add ns1
sudo ip netns add ns2
sudo ip netns add ns3
```

Move one side of each veth link into a correponding namespace and configure an IP address from `169.254.0.0/16` subnet:

```
sudo ip link set xc-1 netns ns1
sudo ip link set xc-2 netns ns2
sudo ip link set xc-3 netns ns3
sudo ip netns exec ns1 ip addr add 169.254.1.10/24 dev xc-1
sudo ip netns exec ns2 ip addr add 169.254.1.20/24 dev xc-2
sudo ip netns exec ns3 ip addr add 169.254.1.30/24 dev xc-3
```

Bring up both sides of veth links

```
sudo ip link set dev xconnect-1 up
sudo ip link set dev xconnect-2 up
sudo ip link set dev xconnect-3 up
sudo ip netns exec ns1 ip link set dev xc-1 up
sudo ip netns exec ns2 ip link set dev xc-2 up
sudo ip netns exec ns3 ip link set dev xc-3 up
```

At this point there should be no connectivity between IPs in individual namespaces, i.e. the following commands will return no output:

```
sudo ip netns exec ns1 ping 169.254.1.20 &
sudo ip netns exec ns1 ping 169.254.1.30 &
```

Start the `xdp-xconnect` app with and connect `xconnect-1` to itself (see [first.yaml](testdata/first.yaml)):

```
cp testdata/first.yaml testdata/input.yaml
sudo go run ./main.go -conf testdata/input.yaml &

2021/03/03 20:08:41 Parsing config file: testdata/input.yaml
2021/03/03 20:08:41 App configuration: {Links:map[xconnect-1:xconnect-1]}
```

Now update the configuration by replacing the file to connect the first two interfaces (see [second.yaml](testdata/second.yaml)):


```
cp testdata/second.yaml testdata/input.yaml

2021/03/03 20:12:16 Parsing config file: testdata/input.yaml
2021/03/03 20:12:16 Error parsing the configuration file: EOF
2021/03/03 20:12:16 Parsing config file: testdata/input.yaml
2021/03/03 20:12:16 App configuration: {Links:map[xconnect-1:xconnect-2]}
64 bytes from 169.254.1.20: icmp_seq=113 ttl=64 time=0.068 ms
64 bytes from 169.254.1.20: icmp_seq=114 ttl=64 time=0.068 ms
64 bytes from 169.254.1.20: icmp_seq=115 ttl=64 time=0.075 ms
```

This proves that the first and second links are now connected. Now swap the connection over to the third link (see [third.yaml](testdata/third.yaml)):

```
cp testdata/third.yaml testdata/input.yaml

2021/03/03 20:13:53 Parsing config file: testdata/input.yaml
2021/03/03 20:13:53 Error parsing the configuration file: EOF
2021/03/03 20:13:53 Parsing config file: testdata/input.yaml
2021/03/03 20:13:53 App configuration: {Links:map[xconnect-1:xconnect-3]
64 bytes from 169.254.1.30: icmp_seq=207 ttl=64 time=0.075 ms
64 bytes from 169.254.1.30: icmp_seq=208 ttl=64 time=0.070 ms
64 bytes from 169.254.1.30: icmp_seq=209 ttl=64 time=0.071 ms
64 bytes from 169.254.1.30: icmp_seq=210 ttl=64 time=0.071 ms
```

Ping replies now come from the third link. Let's see what happens if we provide the malformed input (see [bad.yaml](testdata/bad.yaml))

```
cp testdata/bad.yaml testdata/input.yaml

2021/03/03 20:15:46 Parsing config file: testdata/input.yaml
2021/03/03 20:15:46 App configuration: {Links:map[xconnect-1:xconnect-asd]}
2021/03/03 20:15:46 Error updating eBPF: 1 error occurred:
	* Link not found

64 bytes from 169.254.1.30: icmp_seq=317 ttl=64 time=0.065 ms
64 bytes from 169.254.1.30: icmp_seq=318 ttl=64 time=0.065 ms
64 bytes from 169.254.1.30: icmp_seq=319 ttl=64 time=0.064 m
```

The app detected the error and cross-connect continues working with its last known configuration. Finally, we can terminate the application which will cleanup all of the configured state:


```
fg
[1]  + 1095444 running    sudo go run ./main.go -conf testdata/input.yaml
^C
2021/03/03 20:16:44 Received syscall:interrupt
2021/03/03 20:16:44 ctx.Done
```

Don't forget to cleanup test interfaces and namespaces:

```
sudo ip link del dev xconnect-1
sudo ip link del dev xconnect-2
sudo ip link del dev xconnect-3
sudo ip netns del ns1
sudo ip netns del ns2
sudo ip netns del ns3
```



## Additional Reading and References

https://github.com/xdp-project/xdp-tutorial

https://docs.cilium.io/en/stable/bpf/

https://qmonnet.github.io/whirl-offload/2020/04/12/llvm-ebpf-asm/

https://github.com/takehaya/goxdp-template

https://github.com/hrntknr/nfNat

https://github.com/takehaya/Vinbero

https://github.com/tcfw/vpc

https://github.com/florianl/tc-skeleton

https://github.com/cloudflare/rakelimit

https://github.com/b3a-dev/ebpf-geoip-demo

https://github.com/lmb/ship-bpf-with-go





