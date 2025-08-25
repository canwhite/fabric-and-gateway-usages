# Store.go 代码分析

## 文件功能概述

`store.go`
是链下数据存储的核心实现，负责将区块链上的状态变化（写入操作）持久化到本地文件系统。
它提供了一个简单的链下存储机制，同时包含故障模拟功能用于测试系统的容错能力。

## 核心结构体

### offChainStore

```go
type offChainStore struct {
    path                                    string
    simulatedFailureCount, transactionCount uint
}
```

**字段说明：**

- `path`: 存储文件路径，所有写入操作将追加到此文件
- `simulatedFailureCount`: 模拟故障频率（每 N 次事务后触发一次失败）
- `transactionCount`: 当前已处理的事务计数器

## 主要方法分析

### 1. 构造函数

```go
func newOffChainStore(path string, simulatedFailureCount uint) *offChainStore
```

创建链下存储实例，初始化文件路径和故障模拟参数。

### 2. 核心写入方法

```go
func (ocs *offChainStore) write(data ledgerUpdate) error
```

**功能流程：**

1. 检查是否需要模拟故障
2. 序列化写入数据为 JSON 格式
3. 将数据追加到存储文件

**参数说明：**

- `data`: 包含区块号、交易 ID 和写入操作集合的账本更新

### 3. 故障模拟机制

```go
func (ocs *offChainStore) simulateFailureIfRequired() error
```

**故障逻辑：**

- 当`simulatedFailureCount > 0`且`transactionCount >= simulatedFailureCount`时触发
- 重置计数器并返回`errExpected`错误
- 用于测试系统在写入失败时的恢复能力

### 4. 数据序列化

```go
func (ocs *offChainStore) marshal(writes []write) (string, error)
```

**处理方式：**

- 将每个 write 结构体序列化为 JSON 格式
- 每个 JSON 对象占一行（便于追加和解析）
- 返回拼接后的字符串

### 5. 文件持久化

```go
func (ocs *offChainStore) persist(marshaledWrites string) error
```

**文件操作：**

- 以追加模式打开文件（`os.O_APPEND|os.O_CREATE|os.O_WRONLY`）
- 写入数据后立即关闭文件
- 权限设置为 0644（所有者可读写，组和其他只读）

## 数据格式

### 存储格式

每行一个 JSON 对象，格式如下：

```json
{
  "channelName": "mychannel",
  "namespace": "basic",
  "key": "asset1",
  "isDelete": false,
  "value": "{\"ID\":\"asset1\",\"Color\":\"red\",\"Size\":5,\"Owner\":\"Alice\",\"AppraisedValue\":100}"
}
```

### 数据结构

```go
type ledgerUpdate struct {
    BlockNumber   uint64
    TransactionID string
    Writes        []write
}

type write struct {
    ChannelName string `json:"channelName"`
    Namespace   string `json:"namespace"`
    Key         string `json:"key"`
    IsDelete    bool   `json:"isDelete"`
    Value       string `json:"value"`
}
```

## 使用场景

### 1. 正常操作

- 监听区块链事件，将状态变化实时同步到本地文件
- 用于数据备份、审计和离线分析

### 2. 容错测试

- 通过设置`SIMULATED_FAILURE_COUNT`环境变量来模拟写入失败
- 测试系统能否正确处理存储故障并重试

## 设计特点

### 优点

1. **简单可靠**: 基于文件系统的简单存储，易于理解和维护
2. **幂等性**: 追加写入模式，即使重复处理也不会破坏数据
3. **可观测性**: JSON 格式便于人工阅读和调试
4. **容错测试**: 内置故障模拟机制

### 局限性

1. **性能瓶颈**: 频繁的文件打开/关闭操作可能影响性能
2. **数据增长**: 追加模式导致文件无限增长，需要额外的清理机制
3. **查询能力**: 纯文件存储不支持复杂查询，需要额外的索引机制
4. **并发安全**: 当前实现未考虑并发写入的场景

## 配置方式

通过环境变量控制行为：

- `STORE_FILE`: 存储文件路径（默认：store.log）
- `SIMULATED_FAILURE_COUNT`: 故障模拟频率（默认：0，不模拟）

## 总结

`store.go`实现了一个最小化的链下数据同步系统，专注于将 Fabric 区块链的状态变化可靠地持久化到本地文件系统。它的设计哲学是简单优先，通过清晰的接口和可预测的行为，为上层应用提供稳定的数据存储服务。虽然功能有限，但为更复杂的链下存储系统提供了良好的基础架构。
