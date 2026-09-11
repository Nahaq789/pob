package phase

type Result struct {
	Messages      []string       // 特性や技の追加効果等のメッセージ
	NextPhase     Phase          // 次のフェーズ
	DamageContext *DamageContext // 混乱自傷等でpre_damageが次フェーズ用にセット
	Err           error          // 異常系
}
