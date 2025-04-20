# Redis Packet Decoder


Redis follows RESP protocol https://redis.io/docs/latest/develop/reference/protocol-spec/

There are two major versions of RESP being followed (Version 2 and 3)

1. In RESP2, the protocol starts with ping,+PONG request and response packets
2. In RESP3/RESP2(newer format), the protocol starts with and initial handshake, (hello and the protocol version)


Special Things about Redis,
1. Redis follows a string format for the packets(every item in the packet is in form of a string), and the string uses various symbols to denote the type of data coming in
2. After the symbol, we have the size of the data (in case of Strings, Maps and Arrays), and then the value of the data
3. Each character ends with CRLF which helps denote the end


This module, helps in decoding Redis Packets into a Readable Format(for info about the models https://github.com/khareyash05/redis-packet-decoder/tree/main/models)

Examples - 

```go
import (
    "fmt"
    "log"
    rpd "github.com/khareyash05/redis-packet-decoder"
)

func main() {
    packetToBeProcessed := "*2\r\n$5\r\nhello\r\n$1\r\n3\r\n" // the initial packet in case of RESP3
    decodedPacket, err := rpd.ParseRedis(packetToBeProcessed)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Here is the Decoded Packet", decodedPacket)
}
```