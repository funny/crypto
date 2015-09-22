package main

import (
	"crypto/cipher"
	"crypto/rc4"
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

	lsn, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Print("listen failed: ", err)
		return
	}
	log.Print("server wait")

	conn, err := lsn.Accept()
	if err != nil {
		log.Print("accept failed: ", err)
		return
	}
	log.Print("server accpet")

	writer, reader, err := conn_init(conn)
	if err != nil {
		log.Print("conn init failed: ", err)
		return
	}

	var (
		lastPrintTime   = time.Now()
		sendPacketCount uint64
		recvPacketCount uint64
	)

	for {
		recv := reader.ReadPacket(binary.SplitByUvarint)
		if reader.Error() != nil {
			log.Print("receive failed: ", reader.Error())
			return
		}
		recvPacketCount += 1

		writer.WritePacket(recv, binary.SplitByUvarint)
		if writer.Error() != nil {
			log.Print("send failed: ", writer.Error())
			return
		}
		sendPacketCount += 1

		if time.Since(lastPrintTime) > time.Second*2 {
			lastPrintTime = time.Now()
			log.Print("server: ", recvPacketCount, sendPacketCount)
		}
	}

	log.Print("server exit")
}

// Do DH64 key exchange and return a RC4 reader.
func conn_init(conn net.Conn) (*binary.Writer, *binary.Reader, error) {
	var (
		writer = binary.NewWriter(conn)
		reader = binary.NewReader(conn)
	)

	rand.Seed(time.Now().UnixNano())

	privateKey, publicKey := dh64.KeyPair()
	log.Print("server send public key: ", publicKey)
	writer.WriteUint64BE(publicKey)
	if writer.Error() != nil {
		return nil, nil, writer.Error()
	}

	log.Print("server wait client public key")
	clientPublicKey := reader.ReadUint64BE()
	if reader.Error() != nil {
		return nil, nil, reader.Error()
	}

	log.Print("client public key: ", clientPublicKey)
	secert := dh64.Secret(privateKey, clientPublicKey)
	log.Print("secert: ", secert)

	key := make([]byte, 8)
	binary.PutUint64BE(key, secert)
	rc4stream, err := rc4.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	reader = binary.NewReader(cipher.StreamReader{
		R: conn,
		S: rc4stream,
	})

	return writer, reader, nil
}
