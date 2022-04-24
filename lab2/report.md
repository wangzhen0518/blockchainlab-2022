# 实验名称

**区块链共识协议**
汪震 PB19000078

# 实验目的及要求

1. 实现区块链上的 POW 证明算法
2. 理解区块链上的难度调整的作用

# 实验原理

比特币通过计算出，在限定范围的哈希值，以达到工作量证明。
采用的哈希算法是 sha256, 限定的范围为小于 1<<(256-targetBits), 实际中，targetBits 由出块速度动态调整，本次实验中，targetBits 初始值固定为 5, 每出一个块递增 1。
计算哈希值的对象为，区块头+计数器，计数器初始为 0, 逐次递增 1。
区块头包括，上⼀个区块哈希值(32 位)，当前区块数据对应哈希（32 位，即区块数据的 merkle 根），时间戳，区块难度，计数器(因动态变化而单独说明)。
计数器不断递增，直到计算出的哈希值在限定范围(<1<<(256-targetBits))。

# 实验平台

- 操作系统：Kubuntu 22.04
- go: go1.18.1 linux/amd64

# 实验步骤

1. 完成`proofofwork.go/Run()`函数, 具体代码如下

   ```go
   func (pow *ProofOfWork) Run() (int, []byte) {
       nonce := 0
       hashInt := new(big.Int)
       target := big.NewInt(1)
       target.Lsh(target, 256-pow.block.Bits)

       temp := []byte{}
       temp = append(temp, pow.block.PrevBlockHash...)
       temp = append(temp, pow.block.HashData()...)
       temp = append(temp, IntToHex(pow.block.Timestamp)...)
       temp = append(temp, IntToHex(int64(pow.block.Bits))...)

       temp1 := append(temp, IntToHex(int64(nonce))...)
       temp2 := sha256.Sum256(temp1)
       hashInt.SetBytes(temp2[:])
       for target.Cmp(hashInt) <= 0 {
       	nonce++
       	temp1 = append(temp, IntToHex(int64(nonce))...)
       	temp2 = sha256.Sum256(temp1)
       	hashInt.SetBytes(temp2[:])
       }
       nonce--
       return nonce, temp2[:]
   }
   ```

   依次将区块头内容加入 temp, 将 temp+计数器赋值给 temp1, 哈希结果赋值给 temp2, 不断循环。

2. 完成`proofofwork.go/Validate()`函数，具体代码如下

   ```go
   func (pow *ProofOfWork) Validate() bool {
   	target := big.NewInt(1)
   	target.Lsh(target, 256-pow.block.Bits)
   	hashInt := new(big.Int)
   	hashInt.SetBytes(pow.block.Hash)
   	if hashInt.Cmp(target) < 0 {
   		return true
   	} else {
   		return false
   	}
   }
   ```

   将计算的结果与其限定范围进行比较，判断是否满足条件。

3. 修改了`main.go`的输出，增加了`Bits`(难度值)和`Nonce`(计数器)的输出。

   ```go
   {
   	Name:    "printchain",
   	Aliases: []string{"p"},
   	Usage:   "printchain",
   	Action: func(c *cli.Context) error {
   		bci := bc.Iterator()
   		for {
   			block := bci.Next()
   			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
   			fmt.Printf("Data: %s\n", block.Data)
   			fmt.Printf("Hash: %x\n", block.Hash)
   			pow := NewProofOfWork(block)
   			fmt.Printf("PoW: %s\n", strconv.FormatBoo(pow.Validate()))
   			fmt.Printf("Bits: %d\n", block.Bits)
   			fmt.Printf("Nonce: %d\n", block.Nonce)
   			fmt.Println()
   			if len(block.PrevBlockHash) == 0 {
   				break
   			}
   		}
   		return nil
   	},
   },
   ```

4. 增加了 11 个块，内容为 fourth block~fourteenth block。难度从 5 到 15. 输出见`file.txt`文件，可以看出计算量(Nonce 大小)大体是递增的，不过由于总难度不高，且哈希并不单调，所以计算量并不单调递增，存在波动。


# 实验结果
通过具体实现POW算法，更细致地了解了区块链的共识协议和工作原理，体会到区块链上难度调整的作用。

# 附
修改了的文件都放在`code`目录下，`file.txt`为增加块之后的输出。