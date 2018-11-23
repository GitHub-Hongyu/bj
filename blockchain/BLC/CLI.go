package BLC

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

// CLI结构
type CLI struct {
	BC *BlockChain 
}

// 用法展示
func PrintUsage()  {
	fmt.Println("Usage:")
	fmt.Printf("\tcreateblockchain -address address -- 创建区块链\n")
	fmt.Printf("\taddblock -data DATA -- 交易数据\n")
	fmt.Printf("\tprintchain -- 输出区块链信息\n")
	fmt.Printf("\tsend -from FROM -to TO -amount AMOUNT -- 发起转账交易\n")
	fmt.Printf("\tgetbalance -address FROM -- 查询余额\n")
	fmt.Printf("\tcreatewallet -- 创建钱包\n")
	fmt.Printf("\taddresses -- 获取钱包地址列表\n")
	fmt.Printf("\ttest -- 测试相关方法\n")
	fmt.Printf("\tstartnode -- 启动服务")
}

// 检测参数数量
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		// 如果参数数量不对，直接退出程序
		os.Exit(1)
	}
}

// 添加区块
//func (cli *CLI) addBlock(txs []*Transaction) {
//	if !dbExists() {
//		fmt.Println("数据库不存在...")
//		os.Exit(1)
//	}
//	blockchain := BlockchainObject() // 获取区块链对象
//	defer blockchain.DB.Close()
//	blockchain.AddBlock(txs)
//}

// 命令行运行函数
func (cli *CLI)Run()  {
	// 1. 检测参数数量
	IsValidArgs()
	// 获取环境变量
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	// 2. 新建命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBLCWithGenesisCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	// 创建钱包
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	// 发送交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	// 查询余额
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	// 获取钱包地址列表
	getAddresslistsCmd := flag.NewFlagSet("addresses", flag.ExitOnError)
	// 测试
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	// 添加启动服务命令
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	// 3. 获取命令行参数
	flagAddBlockArg := addBlockCmd.String("data", "send 100 BTC to everyone","交易数据")
	flagCreateBlockchainArg := createBLCWithGenesisCmd.String("address","","the address of create blockchain")
	// 转账命令行参数
	flagFromArg := sendCmd.String("from", "", "转账源地址...")
	flagToArg := sendCmd.String("to", "", "转账目标地址...")
	flagAmount := sendCmd.String("amount", "", "转账金额...")

	// 查询余额命令行参数
	flagBalanceArg := getBalanceCmd.String("address", "", "查询余额...")
	switch os.Args[1] {
	case "startnode":
		if err := startNodeCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of startnode failed! %v\n", err)
		}
	case "test":
		if err := testCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of test failed! %v\n", err)
		}
	case "addresses":
		if err := getAddresslistsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of create wallet failed! %v\n", err)
		}
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of create wallet failed! %v\n", err)
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]);nil != err {
			log.Panicf("parse cmd of getbalance failed! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of send failed! %v\n", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of add block failed! %v\n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of printchain failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of create block chain failed! %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	// 启动服务
	if startNodeCmd.Parsed() {
		cli.startNode(nodeID)
	}
	// 添加测试相关命令
	if testCmd.Parsed() {
		cli.TestResetUTXO()
	}

	// 添加获取钱包地址列表命令
	if getAddresslistsCmd.Parsed() {
		cli.getAddresses(nodeID)
	}
	// 添加创建钱包命令
	if createWalletCmd.Parsed() {
		cli.CreateWallets(nodeID)
	}

	// 添加查询余额命令
	if getBalanceCmd.Parsed() {
		if *flagBalanceArg == "" {
			fmt.Println("未指定查询地址...")
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagBalanceArg, nodeID)
	}
	// 添加转账命令
	if sendCmd.Parsed() {
		if *flagFromArg == "" {
			fmt.Println("源地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagToArg == "" {
			fmt.Println("目标地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagAmount == "" {
			fmt.Println("金额不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagAmount))

		cli.send(JSONToSlice(*flagFromArg), JSONToSlice(*flagToArg), JSONToSlice(*flagAmount), nodeID)
	}
	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		//cli.addBlock([]*Transaction{})
	}

	// 输出区块链信息命令
	if printChainCmd.Parsed() {
		cli.printchain(nodeID)
	}

	// 创建区块链
	if createBLCWithGenesisCmd.Parsed() {
		cli.createBlockchainWithGenesis(*flagCreateBlockchainArg, nodeID)
	}
}