package forge

import "testing"

func rates() []float64 { return []float64{1.0, 0.95, 0.4, 0.6, 0.1, 0.1, 0.1} }

func TestRegisterReadBack(t *testing.T) {
	enhancer := New(rates(), 8, 3, 3, 100)
	enhancer.Register("sword", 2, 1000)
	if enhancer.Level("sword") != 2 || enhancer.Gold("sword") != 1000 {
		t.Fatal("登记后等级与金币应该能读回来")
	}
}

func TestUnknownItemFails(t *testing.T) {
	enhancer := New(rates(), 8, 3, 3, 100)
	if _, _, err := enhancer.Enhance("r1", "ghost", 0.5, false); err != ErrNoItem {
		t.Fatalf("未知装备应该报 ErrNoItem：%v", err)
	}
}

func TestMaxLevelFails(t *testing.T) {
	enhancer := New(rates(), 8, 3, 3, 100)
	enhancer.Register("sword", 8, 1000)
	if _, _, err := enhancer.Enhance("r1", "sword", 0.5, false); err != ErrMaxLevel {
		t.Fatalf("满级应该报 ErrMaxLevel：%v", err)
	}
}

func TestGoldIsCharged(t *testing.T) {
	enhancer := New(rates(), 8, 3, 3, 100)
	enhancer.Register("sword", 2, 1000)
	enhancer.Enhance("r1", "sword", 0.5, false)
	if enhancer.Gold("sword") != 900 {
		t.Fatalf("每次尝试都应该扣 100 金币：%v", enhancer.Gold("sword"))
	}
}

func TestFailsStartsAtZero(t *testing.T) {
	if New(rates(), 8, 3, 3, 100).Fails("sword") != 0 {
		t.Fatal("连续失败次数初始应该是 0")
	}
}
