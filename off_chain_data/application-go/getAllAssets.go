package main

import (
	"encoding/json"
	"fmt"
	atb "offchaindata/contract"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
)
// 这里的连接过程可以总结为：
// 1. 通过 newConnectOptions(clientConnection) 获取身份和连接选项。
// 2. 使用 client.Connect(id, options...) 创建 gateway 客户端，连接到 Fabric 网络。
// 3. 通过 gateway.GetNetwork(channelName) 获取指定通道的网络对象。
// 4. 通过 network.GetContract(chaincodeName) 获取链码合约对象。
// 5. 使用合约对象（这里封装为 atb.NewAssetTransferBasic(contract)）调用链码方法（如 GetAllAssets）。
// 6. 最后关闭 gateway 连接，完成资源释放。
//
// 该流程实现了从 gRPC 连接到 Fabric 网络、获取合约并调用链码方法的完整链路。
//
// 主要步骤：连接 -> 获取网络 -> 获取合约 -> 调用方法 -> 关闭连接

func getAllAssets(clientConnection grpc.ClientConnInterface) error {
	id, options := newConnectOptions(clientConnection)
	gateway, err := client.Connect(id, options...)
	if err != nil {
		return err
	}
	defer gateway.Close()

	contract := gateway.GetNetwork(channelName).GetContract(chaincodeName)
	smartContract := atb.NewAssetTransferBasic(contract)
	assets, err := smartContract.GetAllAssets()
	if err != nil {
		return err
	}

	formatted, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(formatted))

	return nil
}
