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

Test setup

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