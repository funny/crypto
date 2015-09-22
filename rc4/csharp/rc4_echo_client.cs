using System;
using System.IO;
using System.Net.Sockets;

class MainClass
{
	private static DH64 dh64 = new DH64();
	private static Random random = new Random();

	public static void Main(string[] args) {
		TcpClient conn = new TcpClient("127.0.0.1", 10010);

		Stream stream = conn.GetStream();
		BinaryWriter writer = new BinaryWriter(stream);
		BinaryReader reader = new BinaryReader(stream);

		byte[] key = KeyExchange(writer, reader);
		writer = new BinaryWriter(new RC4Stream(stream, key));

		byte[] buffer = new byte[1024];
		for (;;) {
			int length = WriteRandomBytes(writer, buffer);
			uint n = reader.ReadUInt32();
			byte[] recv = reader.ReadBytes((int)n);
			if (!ByteArrayEquals(buffer, recv, length)) {
				Console.WriteLine("send != recv");
				return;
			}
		}
	}

	private static byte[] KeyExchange(BinaryWriter w, BinaryReader r) {
		ulong privateKey;
		ulong publicKey;
		dh64.KeyPair(out privateKey, out publicKey);

		Console.WriteLine("client public key: {0}", publicKey);
		w.Write(publicKey);

		ulong srvPublicKey = r.ReadUInt64();
		Console.WriteLine("server public key: {0}", srvPublicKey);

		ulong secret = dh64.Secret(privateKey, srvPublicKey);
		Console.WriteLine("secret: {0}", secret);

		using (MemoryStream ms = new MemoryStream()) {
			new BinaryWriter(ms).Write(secret);
			return ms.GetBuffer();
		}
	}

	private static int WriteRandomBytes(BinaryWriter w, byte[] buffer) {
		int length = random.Next(buffer.Length);
		for (int i = 0; i < length; i ++) {
			buffer[i] = (byte)random.Next(256);
		}
		Console.WriteLine(length);
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