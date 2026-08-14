package phase

type Result struct {
	Messages  []string // 特性や技の追加効果等のメッセージ
	NextPhase Phase    // 次のフェーズ
	Err       error    // 異常系
}
