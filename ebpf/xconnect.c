#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

#define XCONNECT_MAP_SIZE 1024

struct bpf_map_def SEC("maps") xconnect_map = {
	.type = BPF_MAP_TYPE_DEVMAP,
	.key_size = sizeof(int),
	.value_size = sizeof(int),
	.max_entries = XCONNECT_MAP_SIZE,
};


SEC("xdp")
int  xdp_xconnect(struct xdp_md *ctx)
{
	int ingIdx = ctx->ingress_ifindex;

    return bpf_redirect_map(&xconnect_map, ingIdx, 0);
}


char _license[] SEC("license") = "GPL";


