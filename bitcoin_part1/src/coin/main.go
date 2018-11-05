package main

import (
	"time"
	"fmt"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
)

func main() {
	/*bc := core.NewBlockChain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to lhg")

	for _,block := range bc.Blocks {
		fmt.Printf("Prev.hash:%x\n",block.PrevBlockHash)
		fmt.Printf("Data:%s\n",block.Data)
		fmt.Printf("Hash:%x\n",block.Hash)
		fmt.Printf("Timestamp:%v\n",block.Timestamp)
		fmt.Println(time.Now().Unix())
	}*/

	timestamp := translate("2018-08-22 10:30:11")
	prevBlockHash := []byte("0000000000000000002bb47d1fe591848a44273a006a6193b3720e1d9ebd6e7e")
	//data := []byte("1")
	targetBits := 72
	nonce := 3481895986

	data2 := bytes.Join(
		[][]byte{
			prevBlockHash,
			//data,
			IntToHex(timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce))},
		[]byte{})

	fmt.Println(IntToHex(int64(nonce)))
	hash := sha256.Sum256(data2)
	fmt.Println(string(hash[:]))
}

func translate(toBeCharge string) int64 {
	//获取本地location
	//toBeCharge := "2015-01-01 00:00:00"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	fmt.Println(theTime)                                            //打印输出theTime 2015-01-01 15:15:00 +0800 CST
	fmt.Println(sr)                                                 //打印输出时间戳 1420041600

	//时间戳转日期
	dataTimeStr := time.Unix(sr, 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	fmt.Println(dataTimeStr)

	return sr
}


func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil {
		log.Panic(err)
	}


	return buff.Bytes()
}


/*
func (pow *ProofOfWork) prepareData(nonce int) []byte  {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce))},
		[]byte{})
	return data
}
*/
