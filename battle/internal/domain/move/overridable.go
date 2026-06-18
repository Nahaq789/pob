package move

// OverridableStatusMoves は既存の状態異常を強制的に上書きできる技IDの集合。
// 例: ねむる(173) は既に状態異常がついていても自身を ねむり状態 にできる。
//
// 上書き可否の判定は呼び出し側（move handler）の責務。
// status.Status の API は SetMainStatus（nil時のみ）と ForceSetMainStatus（強制）の
// 2メソッドに分離されており、move 側が技 ID を見て使い分ける。
var OverridableStatusMoves = map[int]struct{}{
	173: {}, // ねむる
}
