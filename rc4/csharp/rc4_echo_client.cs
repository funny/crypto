using System;
using System.IO;
using System.Net.Sockets;
using Funny.Crypto;

// mcs rc4_echo_client.cs rc4.cs ../../dh64/csharp/dh64.cs
class MainClass
{
	private static DH64 dh64 = new DH64();
	private static Random random = new Random();

	public static void Main(string[] args) {
		TcpClient conn = new TcpClient("127.0.0.1", 10010);
		Console.WriteLine("client connect");

		Stream stream = conn.GetStream();
		BinaryReader reader = new BinaryReader(stream);
		BinaryWriter writer = ConnInit(reader);

		byte[] buffer = new byte[1024];
		for (;;) {
			int length = WriteRandomBytes(writer, buffer);
			uint n = reader.ReadUInt32();
			byte[] recv = reader.ReadBytes((int)n);
			if (!ByteArrayEquals(buffer, recv, length)) {
				Console.WriteLine("send != recv");
				Console.WriteLine("send: {0}", BitConverter.ToString(buffer, length));
				Console.WriteLine("recv: {0}", BitConverter.ToString(recv));
				return;
			}
		}
	}

	// Do DH64 key exchange and return a RC4 writer
	private static BinaryWriter ConnInit(BinaryReader r) {
		ulong privateKey;
		ulong publicKey;
		dh64.KeyPair(out privateKey, out publicKey);
		Console.WriteLine("client public key: {0}", publicKey);

		new BinaryWriter(r.BaseStream).Write(publicKey);
		ulong srvPublicKey = r.ReadUInt64();
		Console.WriteLine("server public key: {0}", srvPublicKey);

		ulong secret = dh64.Secret(privateKey, srvPublicKey);
		Console.WriteLine("secret: {0}", secret);

		byte[] key;
		using (MemoryStream ms = new MemoryStream()) {
			new BinaryWriter(ms).Write(secret);
			key = ms.ToArray();
		}
		Console.WriteLine("key: {0}", BitConverter.ToString(key).Replace("-", "").ToLower());

		return new BinaryWriter(
			new RC4Stream(r.BaseStream, key)
		);
	}

	private static int WriteRandomBytes(BinaryWriter w, byte[] buffer) {
		int length = random.Next(buffer.Length);
		for (int i = 0; i < length; i ++) {
			buffer[i] = (byte)random.Next(256);
		}
		w.Write((uint)length);
		w.Write(buffer, 0, length);
		return length;
	}

	private static bool ByteArrayEquals(byte[] a1, byte[] a2, int length) {
		if (a2.Length != length)
			return false;
		for (int i=0; i<length; i++)
			if (a1[i]!=a2[i])
				return false;
		return true;
	}
}