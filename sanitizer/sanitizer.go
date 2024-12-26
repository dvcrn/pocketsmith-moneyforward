package sanitizer

import (
	"strings"

	"github.com/dvcrn/romajiconv"
)

func halfWidthKanaToFullWidth(input string) string {
	replacements := map[string]string{
		"ｱ": "ア", "ｲ": "イ", "ｳ": "ウ", "ｴ": "エ", "ｵ": "オ",
		"ｶ": "カ", "ｷ": "キ", "ｸ": "ク", "ｹ": "ケ", "ｺ": "コ",
		"ｻ": "サ", "ｼ": "シ", "ｽ": "ス", "ｾ": "セ", "ｿ": "ソ",
		"ﾀ": "タ", "ﾁ": "チ", "ﾂ": "ツ", "ﾃ": "テ", "ﾄ": "ト",
		"ﾅ": "ナ", "ﾆ": "ニ", "ﾇ": "ヌ", "ﾈ": "ネ", "ﾉ": "ノ",
		"ﾊ": "ハ", "ﾋ": "ヒ", "ﾌ": "フ", "ﾍ": "ヘ", "ﾎ": "ホ",
		"ﾏ": "マ", "ﾐ": "ミ", "ﾑ": "ム", "ﾒ": "メ", "ﾓ": "モ",
		"ﾔ": "ヤ", "ﾕ": "ユ", "ﾖ": "ヨ",
		"ﾗ": "ラ", "ﾘ": "リ", "ﾙ": "ル", "ﾚ": "レ", "ﾛ": "ロ",
		"ﾜ": "ワ", "ｦ": "ヲ", "ﾝ": "ン",
		"ｧ": "ァ", "ｨ": "ィ", "ｩ": "ゥ", "ｪ": "ェ", "ｫ": "ォ",
		"ｯ": "ッ", "ｬ": "ャ", "ｭ": "ュ", "ｮ": "ョ",
		"ｶﾞ": "ガ", "ｷﾞ": "ギ", "ｸﾞ": "グ", "ｹﾞ": "ゲ", "ｺﾞ": "ゴ",
		"ｻﾞ": "ザ", "ｼﾞ": "ジ", "ｽﾞ": "ズ", "ｾﾞ": "ゼ", "ｿﾞ": "ゾ",
		"ﾀﾞ": "ダ", "ﾁﾞ": "ヂ", "ﾂﾞ": "ヅ", "ﾃﾞ": "デ", "ﾄﾞ": "ド",
		"ﾊﾞ": "バ", "ﾋﾞ": "ビ", "ﾌﾞ": "ブ", "ﾍﾞ": "ベ", "ﾎﾞ": "ボ",
		"ﾊﾟ": "パ", "ﾋﾟ": "ピ", "ﾌﾟ": "プ", "ﾍﾟ": "ペ", "ﾎﾟ": "ポ",
		"｡": "。", "｢": "「", "｣": "」", "､": "、", "･": "・",
		"ｰ": "ー", "-": "ー",
	}

	result := input
	for half, full := range replacements {
		result = strings.ReplaceAll(result, half, full)
	}
	return result
}

func Sanitize(input string) string {
	// convert half-width katakana to full-width katakana
	convertedPayee := halfWidthKanaToFullWidth(input)
	convertedPayee = romajiconv.ConvertFullWidthToHalf(convertedPayee)

	// when /iD is present, add a space so that it's " /iD" instead of "/iD"
	convertedPayee = strings.ReplaceAll(convertedPayee, "/iD", " /iD")
	convertedPayee = strings.ReplaceAll(convertedPayee, "/NFC", " /NFC")

	// strip all consecutive spaces
	convertedPayee = strings.Join(strings.Fields(convertedPayee), " ")

	return convertedPayee
}
