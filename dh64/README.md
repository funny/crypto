64位的迪菲－赫尔曼密钥交换算法代码集
===============================

迪菲－赫尔曼密钥交换（英语：Diffie–Hellman key exchange，简称“D-H”） 是一种安全协议。

它可以让双方在完全没有对方任何预先信息的条件下通过不安全信道创建起一个密钥。这个密钥可以在后续的通讯中作为对称密钥来加密通讯内容。

迪菲－赫尔曼密钥交换的同义词包括:

* 迪菲－赫尔曼密钥协商
* 迪菲－赫尔曼密钥创建
* 迪菲－赫尔曼协议
* 指数密钥交换

虽然迪菲－赫尔曼密钥交换本身是一个匿名（无认证）的 密钥交换协议，它却是很多认证协议的基础，并且被用来提供传输层安全协议的短暂模式中的完备的前向安全性。

考虑到手机游戏项目的平台兼容性和效率要求，我们决定采用64位的迪菲－赫尔曼密钥交换算法，所以建立了这个代码仓库用来收集不同语言的算法实现。

用法
====

随机生成一对64位密钥（私钥 + 公钥）：

```
// C
uint64_t my_private_key;
uint64_t my_public_key;

dh64_key_pair(&my_private_key, &my_public_key);

// Go
myPrivateKey, myPublicKey := dh64.KeyPair()

// C#
DH64 dh64 = new DH64();

ulong myPrivateKey;
ulong myPublicKey;

dh64.KeyPair(out myPrivateKey, out myPublicKey);
```

用以上方式获得公钥之后，就可以通过网络把公钥传递给对方，双方互相拿到对方的公钥之后，利用自己手上的私钥和对方的公钥就可以计算出双方一致的密钥，这样就完成了密钥交换过程：

```
// C
uint64_t secert = dh64_secert(my_private_key, another_publick_key);

// Go
secert := dh64.Secert(myPrivateKey, anotherPublicKey);

// C#
ulong secert = dh64.Secert(myPrivateKey, anotherPublicKey);
```

最终这个密钥就可以用于RC4之类的加密算法来对双发的后续通讯内容做加密了。

进阶用法
=======

默认的生成方式使用的是语言内置的伪随机函数，如果要加强安全性可以自己用真随机来生成一个64位的私钥，然后单独生成DH64交换用的公钥：

```
// C
uint64_t my_private_key = my_real_random64();
uint64_t my_public_key = dh64_public_key(my_private_key);

// Go
myPrivateKey := MyRealRandom64();
myPublicKey := dh64.PublicKey(myPrivateKey);

// C#
ulong myPrivateKey = MyRealRandom64();
ulong myPublicKey = dh64.PublicKey(myPrivateKey);
```

NOTE: 请注意自己生成私钥的时候，私钥必须大于等于1，这是DH算法要求的。

相关链接
=======

1. [迪菲－赫尔曼密钥交换](https://zh.wikipedia.org/wiki/%E8%BF%AA%E8%8F%B2%EF%BC%8D%E8%B5%AB%E5%B0%94%E6%9B%BC%E5%AF%86%E9%92%A5%E4%BA%A4%E6%8D%A2)
2. [云风开源的C版本代码](https://gist.github.com/cloudwu/8838724)
3. [DH64密钥交换 + RC4加密的演示](https://github.com/funny/crypto/tree/master/rc4)
