package expiry_test

import (
	"testing"

	"pob/battle/internal/domain/status/other"
	"pob/battle/internal/domain/status/other/rule/expiry"
)

func TestRampage_Resolve(t *testing.T) {
	succeeded := other.OtherStatusContext{ActionSucceeded: true}
	failed := other.OtherStatusContext{ActionSucceeded: false}

	t.Run("自然終了: 規定ターン完走でcleared+こんらん", func(t *testing.T) {
		r := expiry.NewRampage(2)

		cleared, addConfusion := r.Resolve(succeeded)
		if cleared || addConfusion {
			t.Errorf("turn1: want (false,false), got (%v,%v)", cleared, addConfusion)
		}

		cleared, addConfusion = r.Resolve(succeeded)
		if !cleared || !addConfusion {
			t.Errorf("turn2(last): want (true,true), got (%v,%v)", cleared, addConfusion)
		}
	})

	t.Run("途中失敗: 残りターンが複数あるときの失敗はこんらんなし", func(t *testing.T) {
		r := expiry.NewRampage(3)

		cleared, addConfusion := r.Resolve(failed)
		if !cleared || addConfusion {
			t.Errorf("want (true,false), got (%v,%v)", cleared, addConfusion)
		}
	})

	t.Run("最終ターン失敗: 残り1ターンで失敗したらこんらん付与", func(t *testing.T) {
		r := expiry.NewRampage(2)
		r.Resolve(succeeded) // 1ターン目成功で残り1へ

		cleared, addConfusion := r.Resolve(failed)
		if !cleared || !addConfusion {
			t.Errorf("want (true,true), got (%v,%v)", cleared, addConfusion)
		}
	})

	t.Run("継続中: 成功かつ残りターンありはclearedしない", func(t *testing.T) {
		r := expiry.NewRampage(3)

		cleared, addConfusion := r.Resolve(succeeded)
		if cleared || addConfusion {
			t.Errorf("want (false,false), got (%v,%v)", cleared, addConfusion)
		}
	})
}
