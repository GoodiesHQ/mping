# mping
Simple command line utility to ping multiple targets. By default, it will ping indefinitely until stopped by `Ctrl+C`.

## Flags

    GLOBAL OPTIONS:
        --version, -V               Show the version number and exit
        --ipv4, -4                  Use IPv4 for name resolution
        --ipv6, -6                  Use IPv6 for name resolution
        --count uint, -c uint       Number of echo requests to send to each target (0 = unlimited) (default: 0)
        --interval float, -i float  Interval (in seconds) between sending each packet (minimum 0.01) (default: 1)
        --timeout float, -t float   Timeout (in seconds) to wait for each reply (minimum 0.01) (default: 1)
        --verbose, -v               Enable more verbose logging for debug output
        --help, -h                  show help


## Targets

Positional arguments are treated as a **target** which consts of two parts, a required host and an optional label in the form: `Host[=Label]`

1) The Host, which can be either:
    - IPv4 or IPv6 address
    - FQDN/hostname
    - A URL containing `://` (in which case the label will always be the extracted host)
2) An optional Label, which is assigned by using `=` after the Host (except if it is a URL because the `=` will be treated as part of the URL). If no Label is provided, the Host is used as the label.

#### Examples:

| Target Provided | Host | Output Label |
| --------------- | ---- | ----- |
| `example.com` | `example.com` | `example.com` |
| `example.com=Example` | `example.com` | `Example` |
| `1.1.1.1` | `1.1.1.1` | `1.1.1.1` |
| `1.1.1.1=CFDNS4` | `1.1.1.1` | `CFDNS4` |
| `2606:4700:4700::1111` | `2606:4700:4700::1111` | `2606:4700:4700::1111` |
| `2606:4700:4700::1111=CFDNS6` | `2606:4700:4700::1111` | `CFDNS6` |
| `localhost` | `localhost` | `localhost` |
| `localhost=Me` | `localhost` | `Me` |
| `https://www.example.com/index.php?module=greeting` | `www.example.com` | `www.example.com` |

## Example Output

Pinging both a DNS-resolved host (IPv4 by default) and a direct IPv4 address.

    # mping -c 4 example.com 1.1.1.1

           example.com 1.1.1.1  
           ----------- ---------
        1) 23ms        16ms     
        2) 21ms        16ms     
        3) 21ms        16ms     
        4) 21ms        21ms     

    Total  4           4        
    Fails  0           0        
    Loss   0.0%        0.0%     
    RTT    21ms        17ms


Pinging an IPv4 and an IPv6 address with provided labels to shorten the output.

    # mping -c 4 1.1.1.1=CFDNS4 2606:4700:4700::1111=CFDNS6 notarealhost.local

           CFDNS4    CFDNS6    notarealhost.local
           --------- --------- ------------------
        1) 16ms      16ms      -           
        2) 17ms      14ms      -           
        3) 17ms      13ms      -           
        4) 20ms      14ms      -           

    Total  4         4         4           
    Fails  0         0         4           
    Loss   0.0%      0.0%      100.0%      
    RTT    17ms      14ms      - 

Verbose output showing resolutions and debug information using IPv6 for DNS resolution. Demonstrating flags can go after positional arguments, too.

    # mping -c 4 1.1.1.1=CFDNS4 2606:4700:4700::1111=CFDNS6 one.one.one.one example.com=Example -v -6

    3:25AM DBG Verbose logging enabled
    3:25AM DBG Target CFDNS4 is a valid IP (1.1.1.1), no resolution needed
    3:25AM DBG Target CFDNS6 is a valid IP (2606:4700:4700::1111), no resolution needed
    3:25AM DBG Resolved Target: 1.1.1.1 (CFDNS4)
    3:25AM DBG Resolved Target: 2606:4700:4700::1111 (CFDNS6)
    3:25AM DBG Resolved Target: 2606:4700:4700::1001 (one.one.one.one)
    3:25AM DBG Resolved Target: 2600:1406:5e00:6::17ce:bc12 (Example)
    3:25AM DBG Latency display enabled
    3:25AM DBG Pinging 4 targets count=4

           CFDNS4    CFDNS6    one.one.one.one Example  
           --------- --------- --------------- ---------
        1) 17ms      15ms      15ms            15ms     
        2) 17ms      22ms      14ms            19ms     
        3) 17ms      13ms      15ms            13ms     
        4) 37ms      36ms      36ms            36ms     

    Total  4         4         4               4        
    Fails  0         0         0               0        
    Loss   0.0%      0.0%      0.0%            0.0%     
    RTT    22ms      22ms      20ms            21ms

Similar to the last command, but setting the interval to 0.015 seconds (15ms). If the timeout is greater than the interval, the timeout will be reduced to be the value of the provided interval. Any response which takes more than 15ms will result in a -    status.

    # mping -c 4 1.1.1.1=CFDNS4 2606:4700:4700::1111=CFDNS6 one.one.one.one example.com=Example -6 -i 0.015

           CFDNS4    CFDNS6    one.one.one.one Example  
           --------- --------- --------------- ---------
        1) -         -         14ms            -        
        2) 14ms      13ms      14ms            14ms     
        3) -         -         14ms            -        
        4) -         14ms      -               -        

    Total  4         4         4               4        
    Fails  3         2         1               3        
    Loss   75.0%     50.0%     25.0%           75.0%    
    RTT    14ms      13ms      14ms            14ms

Running without a count, it will run until you it `Ctrl+C`. The labels will re-print after every 10 pings.

    # mping 1.1.1.1=CFDNS4 2606:4700:4700::1111=CFDNS6 one.one.one.one example.com=Example

           CFDNS4    CFDNS6    one.one.one.one Example  
           --------- --------- --------------- ---------
        1) 15ms      15ms      15ms            24ms     
        2) 20ms      20ms      20ms            20ms     
        3) 19ms      19ms      19ms            13ms     
        4) 21ms      15ms      21ms            15ms     
        5) 16ms      15ms      14ms            15ms     
        6) 17ms      15ms      17ms            14ms     
        7) 15ms      14ms      15ms            14ms     
        8) 17ms      15ms      13ms            15ms     
        9) 17ms      17ms      17ms            17ms     
       10) 17ms      13ms      17ms            17ms     
    
           CFDNS4    CFDNS6    one.one.one.one Example  
           --------- --------- --------------- ---------
       11) 18ms      16ms      15ms            15ms     
       12) 18ms      15ms      16ms            15ms     
       13) 18ms      16ms      16ms            15ms     
       14) 17ms      15ms      13ms            15ms     
       15) 17ms      14ms      16ms            14ms     
       16) 15ms      16ms      18ms            15ms     
       17) 18ms      15ms      16ms            16ms     
       18) 17ms      14ms      14ms            14ms     
       19) 16ms      13ms      13ms            13ms     
       20) 18ms      13ms      16ms            16ms     
    
           CFDNS4    CFDNS6    one.one.one.one Example  
           --------- --------- --------------- ---------
       21) 15ms      15ms      16ms            16ms     
       22) 18ms      19ms      15ms            14ms     
       23) 18ms      16ms      16ms            16ms     
       24) 16ms      15ms      15ms            15ms     
       25) 17ms      17ms      17ms            17ms     
    
    Total  25        25        25              25       
    Fails  0         0         0               0        
    Loss   0.0%      0.0%      0.0%            0.0%     
    RTT    17ms      16ms      16ms            16ms