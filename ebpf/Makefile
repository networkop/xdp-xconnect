
BPF        = xconnect.o
BPFCC      = clang
BPFCFLAGS  = -target bpf -Wall -O2 -emit-llvm -g
BPFCFLAGS += -Iinclude
BPFLC      = llc
BPFLCFLAGS = -march=bpf -mcpu=probe -filetype=obj

all: $(BPF)

%.o: %.c
	$(BPFCC) $(BPFCFLAGS) -c $< -o - | \
	$(BPFLC) $(BPFLCFLAGS) -o $@

clean:
	rm -f *.o