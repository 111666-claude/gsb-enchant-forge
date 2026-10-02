// Package forge 是强化系统：成功率表、失败掉级、保护材料、保底与请求幂等。
// 缺陷：成功率表下标差一、掉级条件用了 <=（阈值差一）、保护材料不生效、
// 保底不生效、请求不幂等、每次强化扫描全部装备。
package forge

import "errors"

// 强化失败的原因。
var (
	ErrNoItem   = errors.New("没有这件装备")
	ErrMaxLevel = errors.New("已经满级")
	ErrGold     = errors.New("金币不足")
)

type memo struct {
	level   int
	success bool
}

// Enhancer 是强化台。
type Enhancer struct {
	rates    []float64
	maxLevel int
	guard    int
	pity     int
	cost     int
	level    map[string]int
	gold     map[string]int
	fails    map[string]int
	memo     map[string]memo
	scanned  int
}

// New 建强化台：rates[level] 是当前等级的成功率。
func New(rates []float64, maxLevel, guard, pity, cost int) *Enhancer {
	return &Enhancer{rates: rates, maxLevel: maxLevel, guard: guard, pity: pity, cost: cost,
		level: map[string]int{}, gold: map[string]int{}, fails: map[string]int{},
		memo: map[string]memo{}}
}

// Register 登记一件装备。
func (e *Enhancer) Register(itemID string, level, gold int) {
	e.level[itemID] = level
	e.gold[itemID] = gold
}

// Level 是当前等级。
func (e *Enhancer) Level(itemID string) int { return e.level[itemID] }

// Gold 是剩余金币。
func (e *Enhancer) Gold(itemID string) int { return e.gold[itemID] }

// Fails 是当前连续失败次数。
func (e *Enhancer) Fails(itemID string) int { return e.fails[itemID] }

// Enhance 尝试强化一次，roll 是 0 到 1 的随机数。
// 缺陷：用 level+1 取成功率；掉级条件用 <=；protect 被忽略；pity 被忽略；
// 同一个 reqID 会被重复结算；每次扫描全部装备。
func (e *Enhancer) Enhance(reqID, itemID string, roll float64, protect bool) (int, bool, error) {
	e.scanned += len(e.level)
	level, ok := e.level[itemID]
	if !ok {
		return 0, false, ErrNoItem
	}
	if level >= e.maxLevel {
		return level, false, ErrMaxLevel
	}
	if e.gold[itemID] < e.cost {
		return level, false, ErrGold
	}
	e.gold[itemID] -= e.cost
	index := level + 1
	if index < 0 {
		index = 0
	}
	if index >= len(e.rates) {
		index = len(e.rates) - 1
	}
	success := roll < e.rates[index]
	if success {
		level++
	} else if level > 0 && level <= e.guard {
		level--
	}
	e.fails[itemID]++
	e.level[itemID] = level
	return level, success, nil
}

// Scanned 是累计扫描量（规模观测）。
func (e *Enhancer) Scanned() int { return e.scanned }
