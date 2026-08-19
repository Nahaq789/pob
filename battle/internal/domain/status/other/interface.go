package other

// MoveBlocker は特定の技の使用を封じる OtherStatus が実装するインターフェース。
type MoveBlocker interface {
	BlocksMoveId(moveId int) bool
}

// MoveForcer は特定の技の使用を強制する OtherStatus が実装するインターフェース。
type MoveForcer interface {
	ForcesMoveId() int
}
