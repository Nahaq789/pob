package phase

import (
	"testing"

	"pob/battle/internal/domain/damage"
)

type registryTestEntry struct{ tag int }

func (h *registryTestEntry) Handle(ctx EntryContext) Result { return Result{} }

type registryTestExit struct{ tag int }

func (h *registryTestExit) Handle(ctx ExitContext) Result { return Result{} }

// registryTestHandler は DamageModHandler / MoveHandler / OnHitHandler / PostMoveHandler を兼ねるダミー。
type registryTestHandler struct{ tag int }

func (h *registryTestHandler) Mod(ctx MoveResolveContext) damage.DamageMod { return damage.DamageMod{} }
func (h *registryTestHandler) AfterEffect(ctx MoveResolveContext) Result   { return Result{} }
func (h *registryTestHandler) HitCount() int                               { return 1 }
func (h *registryTestHandler) OnHit(ctx MoveResolveContext) Result         { return Result{} }
func (h *registryTestHandler) TurnEnd(ctx PostMoveContext) Result          { return Result{} }

func registryTestTag(h any) (int, bool) {
	switch v := h.(type) {
	case *registryTestEntry:
		return v.tag, true
	case *registryTestExit:
		return v.tag, true
	case *registryTestHandler:
		return v.tag, true
	}
	return 0, false
}

func registryTestLookup[H any](m map[int]H, id int) (int, bool) {
	h, ok := m[id]
	if !ok {
		return 0, false
	}
	return registryTestTag(h)
}

func TestRegistry_RegisterById(t *testing.T) {
	tests := []struct {
		name     string
		register func(r *Registry, id, tag int)
		lookup   func(r *Registry, id int) (int, bool)
	}{
		{
			"RegisterEntryAbility",
			func(r *Registry, id, tag int) { r.RegisterEntryAbility(id, &registryTestEntry{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.entryAbilityHandlers, id) },
		},
		{
			"RegisterEntryItem",
			func(r *Registry, id, tag int) { r.RegisterEntryItem(id, &registryTestEntry{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.entryItemHandlers, id) },
		},
		{
			"RegisterExitAbility",
			func(r *Registry, id, tag int) { r.RegisterExitAbility(id, &registryTestExit{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.exitAbilityHandlers, id) },
		},
		{
			"RegisterExitItem",
			func(r *Registry, id, tag int) { r.RegisterExitItem(id, &registryTestExit{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.exitItemHandlers, id) },
		},
		{
			"RegisterResolveAbility",
			func(r *Registry, id, tag int) { r.RegisterResolveAbility(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.resolveAbilityHandlers, id) },
		},
		{
			"RegisterResolveItem",
			func(r *Registry, id, tag int) { r.RegisterResolveItem(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.resolveItemHandlers, id) },
		},
		{
			"RegisterResolveMove",
			func(r *Registry, id, tag int) { r.RegisterResolveMove(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.resolveMoveHandlers, id) },
		},
		{
			"RegisterOnHitAbility",
			func(r *Registry, id, tag int) { r.RegisterOnHitAbility(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.onHitAbilityHandlers, id) },
		},
		{
			"RegisterOnHitItem",
			func(r *Registry, id, tag int) { r.RegisterOnHitItem(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.onHitItemHandlers, id) },
		},
		{
			"RegisterPostMoveAbility",
			func(r *Registry, id, tag int) { r.RegisterPostMoveAbility(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.postMoveAbilityHandlers, id) },
		},
		{
			"RegisterPostMoveItem",
			func(r *Registry, id, tag int) { r.RegisterPostMoveItem(id, &registryTestHandler{tag}) },
			func(r *Registry, id int) (int, bool) { return registryTestLookup(r.postMoveItemHandlers, id) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" 登録: 同じIDで取得できる", func(t *testing.T) {
			r := NewRegistry()
			tt.register(r, 10, 1)
			tag, ok := tt.lookup(r, 10)
			if !ok {
				t.Fatal("expected handler to be registered")
			}
			if tag != 1 {
				t.Errorf("expected tag 1, got %d", tag)
			}
		})

		t.Run(tt.name+" 未登録ID: 取得できない", func(t *testing.T) {
			r := NewRegistry()
			tt.register(r, 10, 1)
			if _, ok := tt.lookup(r, 11); ok {
				t.Error("expected no handler for unregistered id")
			}
		})

		t.Run(tt.name+" 同じIDで再登録: 上書きされる", func(t *testing.T) {
			r := NewRegistry()
			tt.register(r, 10, 1)
			tt.register(r, 10, 2)
			tag, ok := tt.lookup(r, 10)
			if !ok {
				t.Fatal("expected handler to be registered")
			}
			if tag != 2 {
				t.Errorf("expected tag 2, got %d", tag)
			}
		})

		t.Run(tt.name+" 別IDで登録: 両方保持される", func(t *testing.T) {
			r := NewRegistry()
			tt.register(r, 10, 1)
			tt.register(r, 20, 2)
			if tag, ok := tt.lookup(r, 10); !ok || tag != 1 {
				t.Errorf("id=10: expected tag 1, got %d (ok=%v)", tag, ok)
			}
			if tag, ok := tt.lookup(r, 20); !ok || tag != 2 {
				t.Errorf("id=20: expected tag 2, got %d (ok=%v)", tag, ok)
			}
		})
	}
}

func TestRegistry_AddResolveBase(t *testing.T) {
	t.Run("1件登録: 1件保持される", func(t *testing.T) {
		r := NewRegistry()
		r.AddResolveBase(&registryTestHandler{1})
		if len(r.resolveBaseHandlers) != 1 {
			t.Errorf("expected 1 handler, got %d", len(r.resolveBaseHandlers))
		}
	})

	t.Run("複数登録: 登録順が保持される", func(t *testing.T) {
		r := NewRegistry()
		for _, tag := range []int{1, 2, 3} {
			r.AddResolveBase(&registryTestHandler{tag})
		}
		if len(r.resolveBaseHandlers) != 3 {
			t.Fatalf("expected 3 handlers, got %d", len(r.resolveBaseHandlers))
		}
		for i, want := range []int{1, 2, 3} {
			got, _ := registryTestTag(r.resolveBaseHandlers[i])
			if got != want {
				t.Errorf("index %d: expected tag %d, got %d", i, want, got)
			}
		}
	})

	t.Run("同じハンドラーを2回登録: 上書きされず2件になる", func(t *testing.T) {
		r := NewRegistry()
		h := &registryTestHandler{1}
		r.AddResolveBase(h)
		r.AddResolveBase(h)
		if len(r.resolveBaseHandlers) != 2 {
			t.Errorf("expected 2 handlers, got %d", len(r.resolveBaseHandlers))
		}
	})
}
