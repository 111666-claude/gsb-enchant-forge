# enchant-forge

强化系统：成功率表、失败掉级、保护材料、保底与请求幂等，只用 Go 标准库。

```
go run ./cmd/enchant-forge --sample=grade
go run ./cmd/enchant-forge --sample=degrade
go run ./cmd/enchant-forge --sample=protect
go run ./cmd/enchant-forge --sample=pity
go run ./cmd/enchant-forge --sample=dup
go run ./cmd/enchant-forge --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **成功率**：当前等级为 `L` 时取 `rates[L]`，成功条件是 `roll < rates[L]`。
- **费用**：每次尝试（成功或失败）都扣 `cost` 金币；金币不足返回 `ErrGold` 且不扣钱。
- **满级**：等级达到 `maxLevel` 返回 `ErrMaxLevel`，不扣钱、不判定。
- **失败掉级**：失败且等级大于 0 且等级不小于 `guard` 时掉一级；
  `protect` 为真时不掉级，但材料照扣。等级不会低于 0。
- **保底**：连续失败达到 `pity` 次后，下一次必定成功（不看 `roll`）；
  成功一次后连续失败计数清零。
- **幂等**：同一个 `reqID` 只结算一次，重复调用返回上次的等级与成功标记，
  不扣钱、不改连续失败计数。
- **代价**：强化不许扫描全部装备，`Scanned` 不随装备数线性增长。
- **规模**：10 万件装备、每秒 10 万次强化，单次摊还 O(1)。

## 输出契约（不改格式）

```
level=.. success=..
level=..
level=.. gold=..
scanned=..
```

## 目录

```
forge.go              成功率、掉级、保护、保底与幂等
forge_test.go         单元测试
cmd/enchant-forge/    命令行入口
```
