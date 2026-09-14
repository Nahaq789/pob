package other

import "pob/battle/internal/domain/status"

// MoveBlocker は特定の技の使用を封じる OtherStatus が実装するインターフェース。
type MoveBlocker interface {
	BlocksMoveId(moveId int) bool
}

// MoveForcer は特定の技の使用を強制する OtherStatus が実装するインターフェース。
type MoveForcer interface {
	ForcesMoveId() int
}

// Flincher はひるみによる行動ブロックを処理するインターフェース。
type Flincher interface {
	Handle(name string) string
}

// PreMoveChecker はアンコール・かなしばり等の技選択制限を処理するインターフェース。
type PreMoveChecker interface {
	Handle(moveId int, others []status.OtherStatus) (string, bool)
}

// SelfHitter はこんらん自傷の判定を処理するインターフェース。
type SelfHitter interface {
	CheckSelfHit(name string) (string, bool)
}
