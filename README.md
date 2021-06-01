# GTP DSCP Copy

Tool to copy DSCP field of inner GTP packet to outer IP header written in Golang. It creates an iptable u32 command and appends it to PREROUTING-MANGLE - Check https://github.com/shynuu/slice-aware-ntn for full description
