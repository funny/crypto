MT19937伪随机算法代码集
=====================

收集整理了不同语言的MT19937随机算法的实现，并做了测试。

可以用于客户端运算服务端验算型的游戏项目，由服务端下发随机种子，之后每次随机的结果都是服务端可验证的。

相关链接
=======

1. [梅森旋转算法](http://zh.wikipedia.org/wiki/%E6%A2%85%E6%A3%AE%E6%97%8B%E8%BD%AC%E7%AE%97%E6%B3%95)
2. [C语言版的实现](http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt64.html)
3. [C#版本的实现](https://github.com/M-S-D/Team-Splitter/blob/master/mt19937.cs)
4. [Go语言版的实现](https://github.com/seehuhn/mt19937/blob/master/mt19937.go)

注1：C#版实现用了静态写法，本仓库收集整理时重构了命名方式并改为非静态用法。
注2：Go语言版没有浮点数生成的部分，本仓库收集整理时添加了这部分代码。
