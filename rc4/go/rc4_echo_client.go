package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/rc4"
	"encoding/hex"
	"flag"
	"github.com/funny/binary"
	dh64 "github.com/funny/crypto/dh64/go"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:10010", "server address")
	flag.Parse()

	log.Print("client dial")
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Print("connect failed: ", err)
		return
	}
	log.Print("client connect")

	writer, reader, err := conn_init(conn)
	if err != nil {
		log.Print("conn init failed: ", err)
		return
	}

	var (
		send            = make([]byte, 0, 1024)
		lastPrintTime   = time.Now()
		sendPacketCount uint64
		recvPacketCount uint64
	)

	for {
		send = randomData(send)
		writer.WritePacket(send, binary.SplitByUint32LE)
		if writer.Error() != nil {
			log.Print("send failed: ", writer.Error())
			return
		}
		sendPacketCount += 1

		recv := reader.ReadPacket(binary.SplitByUint32LE)
		if reader.Error() != nil {
			log.Print("receive failed: ", reader.Error())
			return
		}
		recvPacketCount += 1

		if !bytes.Equal(send, recv) {
			log.Print("send != recv")
			log.Print("send: ", hex.EncodeToString(send))
			log.Print("recv: ", hex.EncodeToString(recv))
			return
		}

		if time.Since(lastPrintTime) > time.Second*2 {
			lastPrintTime = time.Now()
			log.Print("client: ", recvPacketCount, sendPacketCount)
		}
	}

	log.Print("client exit")
}

// Do DH64 key exchange and return a RC4 writer.
func conn_init(conn net.Conn) (*binary.Writer, *binary.Reader, error) {
	var (
		writer = binary.NewWriter(conn)
		reader = binary.NewReader(conn)
	)

	rand.Seed(time.Now().UnixNano())

	privateKey, publicKey := dh64.KeyPair()
	log.Print("client send public key: ", publicKey)
	writer.WriteUint64LE(publicKey)
	if writer.Error() != nil {
		return nil, nil, writer.Error()
	}

	log.Print("client wait server public key")
	serverPublicKey := reader.ReadUint64LE()
	if reader.Error() != nil {
		return nil, nil, reader.Error()
	}

	log.Print("server public key: ", serverPublicKey)
	secert := dh64.Secret(privateKey, serverPublicKey)
	log.Print("secert: ", secert)

	key := make([]byte, 8)
	binary.PutUint64LE(key, secert)
	rc4stream, err := rc4.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	writer = binary.NewWriter(cipher.StreamWriter{
		W: conn,
		S: rc4stream,
	})

	return writer, reader, nil
}

func randomData(b []byte) []byte {
	a := b[:rand.Intn(cap(b))]
	for i := 0; i < len(a); i++ {
		a[i] = byte(rand.Int() % 256)
	}
	return a
}
