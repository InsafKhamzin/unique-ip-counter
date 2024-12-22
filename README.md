## Unique ips counter

You have a simple text file with IPv4 addresses. One line is one address, line by line:

```
145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
```

The file is unlimited in size and can occupy tens and hundreds of gigabytes.

You should calculate the number of unique addresses in this file using as little memory and time as possible.

### How to run 
To run

```
  make run path=C:/ip_addresses
```

To test

```
  make test
```

### Solution
- IPv4 is 32 bits long, so it's possible to have 2^32 unique addresses.
- Each unique address could be represented as a bitflag
- Bitflag array will be predictable in size, 2^32 bits (~536MB)
- Program reads the ip line, converts to int8 and shifts to the corresponding byte index
- To speed up the process, parallel processing where involved by spawning goroutines and assigning file chunks boundries to them
- There is a global bitset that each goroutine writes to, it should be locked properly to ensure data integrity
- By applying frequent locks to the global bitset we risk degrading the perfomance of the program
- To solve that, the global bitset was sharder to 256 smaller 2^24 bitsets, index of shards array is the first number ip
- Program locks corresponding bitset array, this saves us from too frequent locking
- Eventually, the final Count of ip addresses is aggregated from each of 256 bitset


### Benchmark

- The program was able to count 1bln unique addresses from 110GB file in about 60sec on average.
- Hardware: 20 core intel i7-12700K, SSD Samsung 970 evo plus