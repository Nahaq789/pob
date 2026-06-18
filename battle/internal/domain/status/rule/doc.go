// Package rule は状態異常関連のルールと handler を提供する。
//
// 「状態異常であることが原因で発生する事象」（やけどなら攻撃半減、
// 毒なら毎ターン1/8ダメージ等）はすべて status の責務であり、
// pure function とタイミング込みの handler を本パッケージに配置する。
//
// 各 handler は phase.EndOfTurnHandler 等の共通 interface を実装し、
// 起動時に registry にフラットに登録される。
//
// このパッケージは pokemon を import しない（循環参照回避）。
// Pokemon へのアクセスは ctx 経由で行う。
package rule
