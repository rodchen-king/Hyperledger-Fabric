# 账本举例
## 传统账本
假设三个程序员A,B,C住在一起，大家平时生活在一起，生活上难免现金周转不开，出现向别人借钱的状况。大家是将现金交易写在一个黑板上的。经过一个月以后，出现下面的账本。
```
A借B 100元，2018/09/01
A借C 200元，2018/09/05
A还B 50元， 2018/09/10
C借B 150元，2018/09/15
C还B 100元，2018/09/20
```

但是有一天，A将第三条记录修改为
```
A还B 100元， 2018/09/10
```
这样账本就不够安全。

## 借鉴数字摘要（Digital Digest）技术
没添加一条记录前，首先计算前面所有历史记录的摘要值，也就是一个hash值。那现在我们将上面的内容改写成带有数字摘要的。
数字摘要的计算方式：首先计算前面所有记录的sha256(前面所有记录)值，然后和将要添加的计算一起写在账本的最下方。
例如：
```
第一条记录摘要值：括号内为空，因为前面没有记录，sha256() = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
第二天记录摘要值：sha256(e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 A借B 100元，2018/09/01)
第三条记录摘要值：sha256(e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 A借B 100元，2018/09/01
a6a2a7fdeffee901ffbc3a9423036f844b5a8004ff5458aa4e1ff63f1096a248 A借C 200元，2018/09/05)
......
```
```
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 A借B 100元，2018/09/01
a6a2a7fdeffee901ffbc3a9423036f844b5a8004ff5458aa4e1ff63f1096a248 A借C 200元，2018/09/05
4f49d2619f653edda7b22e1cfbab5d51961844cfc53adffea4c5ca23f4c6bbc9 A还B 50元， 2018/09/10
878843ca6dd19cae8ed20e515a4e2df5ddd06fcb989129bf8aa68ee8b9a2b0cf C借B 150元，2018/09/15
ebc6d917e3d723386780e3c2489418f0ff52116ca21bd1aa309f42ce0fa3f8f5 C还B 100元，2018/09/20
```

如果此时A将第三条记录内容修改为
```
A还B 100元， 2018/09/10
```
```
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 A借B 100元，2018/09/01
a6a2a7fdeffee901ffbc3a9423036f844b5a8004ff5458aa4e1ff63f1096a248 A借C 200元，2018/09/05
4f49d2619f653edda7b22e1cfbab5d51961844cfc53adffea4c5ca23f4c6bbc9 A还B 100元， 2018/09/10
878843ca6dd19cae8ed20e515a4e2df5ddd06fcb989129bf8aa68ee8b9a2b0cf C借B 150元，2018/09/15
ebc6d917e3d723386780e3c2489418f0ff52116ca21bd1aa309f42ce0fa3f8f5 C还B 100元，2018/09/20
```

那么我们可以通过计算第四条记录的摘要值里发现前面记录被修改过，因为第四条记录的值为：
**8bd9f0621c05578a96ea4c1745cfbb99e2aba54ecb3e003d77e255a1f0cc389a**

## 增量摘要值
上面的方法存问题，是随着账本的增加，每一次计算都需要将前面所有的记录计算一次，计算成本过高，基于sha256算法的特性。内容上一个字母的变化计算结果都会有很大的区别。所以这里采用下面方式进行计算：
```
第一条记录为: "", A借B 100元，2018/09/01, sha("", A借B 100元，2018/09/01) => "", A借B 100元, f748da464d2693dd9a80febf27b251a8dc3592387ee35741620c9c96cdd71fb7
然后将第一条记录得到的摘要值传给第二条记录。
第二条记录为：
f748da464d2693dd9a80febf27b251a8dc3592387ee35741620c9c96cdd71fb7, A借C 200元, 2018/09/05, sha(f748da464d2693dd9a80febf27b251a8dc3592387ee35741620c9c96cdd71fb7, A借C 200元) 
=> f748da464d2693dd9a80febf27b251a8dc3592387ee35741620c9c96cdd71fb7, A借C 200元, ffaa5c0ccb857fba5019f29e2b12cc6cf395d667f582c735c32be104c171d684
```

这样做的原因是因为，如果第一条记录内容有任何一个字母更改，得到的摘要值都会有很大的变化，从第一条记录中就会发现当前的摘要值不对，有人更改了内容。如果连摘要值都修改了，那么第二条记录也可以发现与前一条记录的摘要值不相等，就会发现第一条记录的内容改了，而且摘要值也被改了。所以我们这种那前一个记录计算的摘要值和当前记录内容一起生成行的摘要值也是可以达到(2)的效果，而且计算起来很方便。

这里有一个问题就是如果修改记录的人把所有的摘要值都重新计算一便，那么所有人都没法发现了。这个问题先抛出来，后面再谈。

# 区块链底层结构
区块链网络底层结构类似上面第三种。但是有点不同，现在我们将上面的方式换一下。
就是每一天的所有记录为一个单位计算摘要值，而不是每一条记录都记录摘要值。这样我们得到的结论就是若干个记录会在一起进行摘要值的计算，并且其中会带有前一个单位的摘要值和当前单位的摘要值。然后就可以开始类比区块链网络结构了。


## blockchain结构

### 区块结构体定义
    
    type Block struct {
        Index int64
        TimeStamp int64
        PrevBlockHash string
        Hash string
        Data string
    }
    
### 区块链结构体定义
    
    type BlockChain struct {
        Block []*Block
    }

### 计算区块的 Hash 值

    func calculateHash(block Block) string {
        data := string(block.Index) + string(block.TimeStamp) + block.PrevBlockHash + block.Data
        blockInBytes := sha256.Sum256([]byte(data))
        return hex.EncodeToString(blockInBytes[:])
    }    

## 运行

    // 进入目录
    $ cd $GOPATH/src/BlockChain
    
    // 执行命令
    $ go run ./cmd/main.go
    
    Index: 0
    PrevBlockHash:
    CurrHash: 0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f
    Data: Genesis Block
    Timestamp: 1539874076

    Index: 1
    PrevBlockHash: 0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f
    CurrHash: b11305449703848e79f02f0ba7f7db6bdd085a4a5ea50382ea4cca77644c376b
    Data: A借B 100元，2018/09/01
    Timestamp: 1539874076

    Index: 2
    PrevBlockHash: b11305449703848e79f02f0ba7f7db6bdd085a4a5ea50382ea4cca77644c376b
    CurrHash: 751c3793ee3492f5e050c6b662f4d832bc125dde0aae813147e5459abc23f29a
    Data: A借C 200元，2018/09/05
    Timestamp: 1539874076

    Index: 3
    PrevBlockHash: 751c3793ee3492f5e050c6b662f4d832bc125dde0aae813147e5459abc23f29a
    CurrHash: fe1b5200d9c079f8533f5c5cb0b80ca6438adb9cb5fe3756d6a03910ca50fd1f
    Data: A还B 50元， 2018/09/10
    Timestamp: 1539874076

    Index: 4
    PrevBlockHash: fe1b5200d9c079f8533f5c5cb0b80ca6438adb9cb5fe3756d6a03910ca50fd1f
    CurrHash: 717abf4372e4530adec052e048f75eae4b95f1072d2e66fa4e28ba25112e544e
    Data: C借B 150元，2018/09/15
    Timestamp: 1539874076

    Index: 5
    PrevBlockHash: 717abf4372e4530adec052e048f75eae4b95f1072d2e66fa4e28ba25112e544e
    CurrHash: ae658dd10062bd4034c661ff3423186c4af729de2ab4968f245fcc898b3fe6bb
    Data: C还B 100元，2018/09/20
    Timestamp: 1539874076

## 通过 RPC 接口访问数据

### 启动服务

    // 进入目录
    $ cd $GOPATH/src/BlockChain
    
    // 开启http服务监听
    $ go run rpc/Server.go
    

### 查看数据

打开浏览器输入 URL 地址：`http://localhost:8000/blockchain/get`

    {
        "Block": [
            {
                "Index": 0,
                "TimeStamp": 1539490490,
                "PrevBlockHash": "",
                "Hash": "0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f",
                "Data": "Genesis Block"
            }
        ]
    }

### 写入区块链数据

打开浏览器输入 URL 地址：`http://localhost:8000/blockchain/write?data=A借B 100元，2018/09/01`

    {
        "Block": [{
            "Index": 0,
            "TimeStamp": 1539874242,
            "PrevBlockHash": "",
            "Hash": "0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f",
            "Data": "Genesis Block"
        }, {
            "Index": 1,
            "TimeStamp": 1539874384,
            "PrevBlockHash": "0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f",
            "Hash": "b11305449703848e79f02f0ba7f7db6bdd085a4a5ea50382ea4cca77644c376b",
            "Data": "A借B 100元，2018/09/01"
        }]
    }

### 写入所有交易
    {
        "Block": [{
            "Index": 0,
            "TimeStamp": 1539874242,
            "PrevBlockHash": "",
            "Hash": "0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f",
            "Data": "Genesis Block"
        }, {
            "Index": 1,
            "TimeStamp": 1539874384,
            "PrevBlockHash": "0d8845eb2da42f75aef4ee920f644975d73347e0331d17b37209c4f32ef4867f",
            "Hash": "b11305449703848e79f02f0ba7f7db6bdd085a4a5ea50382ea4cca77644c376b",
            "Data": "A借B 100元，2018/09/01"
        }, {
            "Index": 2,
            "TimeStamp": 1539874483,
            "PrevBlockHash": "b11305449703848e79f02f0ba7f7db6bdd085a4a5ea50382ea4cca77644c376b",
            "Hash": "751c3793ee3492f5e050c6b662f4d832bc125dde0aae813147e5459abc23f29a",
            "Data": "A借C 200元，2018/09/05"
        }, {
            "Index": 3,
            "TimeStamp": 1539874498,
            "PrevBlockHash": "751c3793ee3492f5e050c6b662f4d832bc125dde0aae813147e5459abc23f29a",
            "Hash": "fe1b5200d9c079f8533f5c5cb0b80ca6438adb9cb5fe3756d6a03910ca50fd1f",
            "Data": "A还B 50元， 2018/09/10"
        }, {
            "Index": 4,
            "TimeStamp": 1539874504,
            "PrevBlockHash": "fe1b5200d9c079f8533f5c5cb0b80ca6438adb9cb5fe3756d6a03910ca50fd1f",
            "Hash": "717abf4372e4530adec052e048f75eae4b95f1072d2e66fa4e28ba25112e544e",
            "Data": "C借B 150元，2018/09/15"
        }, {
            "Index": 5,
            "TimeStamp": 1539874511,
            "PrevBlockHash": "717abf4372e4530adec052e048f75eae4b95f1072d2e66fa4e28ba25112e544e",
            "Hash": "ae658dd10062bd4034c661ff3423186c4af729de2ab4968f245fcc898b3fe6bb",
            "Data": "C还B 100元，2018/09/20"
        }]
    }
