// Package rule は天候・フィールド・ルーム関連のルールと handler を提供する。
//
// 「天候・フィールド・ルームが原因で発生する事象」（雨で水技1.5倍、
// すなあらしで岩・地面・鋼以外にダメージ等）は ground の責務であり、
// pure function とタイミング込みの handler を本パッケージに配置する。
package rule
