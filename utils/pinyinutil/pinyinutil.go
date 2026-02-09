// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pinyinutil

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// ExtractInitials 从完整拼音中提取首字母
// 例如："limingyue" -> "lmy", "zhangsan" -> "zs", "wangwu" -> "ww"
// 参数:
// - pinyin: 完整拼音字符串，多个拼音用空格或连字符分隔，如 "li ming yue" 或 "li-ming-yue"
// 返回: 首字母字符串，如 "lmy"
func ExtractInitials(pinyin string) string {
	if pinyin == "" {
		return ""
	}

	// 统一处理分隔符：将空格、连字符、下划线等统一替换为空格
	pinyin = regexp.MustCompile(`[\s\-_]+`).ReplaceAllString(pinyin, " ")
	pinyin = strings.TrimSpace(pinyin)
	if pinyin == "" {
		return ""
	}

	// 按空格分割
	parts := strings.Fields(pinyin)
	if len(parts) == 0 {
		return ""
	}

	// 提取每个部分的首字母
	initials := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 取第一个字符（小写）
		firstChar := strings.ToLower(string(part[0]))
		// 只保留字母
		if matched, _ := regexp.MatchString(`^[a-z]$`, firstChar); matched {
			initials = append(initials, firstChar)
		}
	}

	return strings.Join(initials, "")
}

// IsLikelyInitials 判断字符串是否可能是首字母搜索
// 如果字符串全是小写字母且长度较短（1-10个字符），则可能是首字母搜索
// 参数:
// - s: 待判断的字符串
// 返回: true 表示可能是首字母搜索，false 表示不是
func IsLikelyInitials(s string) bool {
	if s == "" {
		return false
	}
	// 去除空格
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// 长度限制：1-10个字符
	if len(s) < 1 || len(s) > 10 {
		return false
	}
	// 必须全是小写字母
	matched, _ := regexp.MatchString(`^[a-z]+$`, s)
	return matched
}

// ExtractInitialsFromNickname 从中文昵称提取拼音首字母
// 例如："林大伯" -> "ldb", "李明月" -> "lmy", "张三" -> "zs"
// 参数:
// - nickname: 中文昵称，可能包含中文、英文、数字等
// 返回: 拼音首字母字符串，如 "ldb"
// 注意:
// - 如果昵称包含中文，会将中文转换为拼音并提取首字母
// - 如果昵称只包含英文或数字，会保留英文首字母，忽略数字
// - 如果昵称为空，返回空字符串
func ExtractInitialsFromNickname(nickname string) string {
	if nickname == "" {
		return ""
	}

	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return ""
	}

	// 使用 go-pinyin 将中文转换为拼音（不带声调，用空格分隔）
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal // 不带声调，如 "zhong guo ren"
	pinyinParts := pinyin.LazyPinyin(nickname, args)

	if len(pinyinParts) == 0 {
		// 如果没有中文，可能是纯英文或数字，尝试提取英文首字母
		return extractEnglishInitials(nickname)
	}

	// 提取每个拼音的首字母
	initials := make([]string, 0, len(pinyinParts))
	for _, part := range pinyinParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 取第一个字符（小写）
		firstChar := strings.ToLower(string(part[0]))
		// 只保留字母
		if matched, _ := regexp.MatchString(`^[a-z]$`, firstChar); matched {
			initials = append(initials, firstChar)
		}
	}

	return strings.Join(initials, "")
}

// GeneratePinyinAndInitials 从名称生成拼音和首字母（支持中文、英文、数字）
// 参数:
//   - name: 名称，可能包含中文、英文、数字等
// 返回:
//   - pinyin: 完整拼音字符串（小写，去除空格，如 "zhangsan" 或 "techteam"）
//   - pinyinInitials: 首字母字符串（小写，如 "zs" 或 "tt"）
// 注意:
//   - 如果名称包含中文，会将中文转换为拼音并提取首字母
//   - 如果名称只包含英文或数字，会保留英文拼音，提取首字母
//   - 如果名称为空，返回空字符串
func GeneratePinyinAndInitials(name string) (string, string) {
	if name == "" {
		return "", ""
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}

	// 使用 go-pinyin 将中文转换为拼音（不带声调，用空格分隔）
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal // 不带声调，如 "zhong guo ren"
	pinyinParts := pinyin.LazyPinyin(name, args)

	var pinyinStr string
	var initials string

	if len(pinyinParts) > 0 {
		// 有中文，使用拼音
		// 生成完整拼音：去除空格，连接所有拼音部分
		pinyinStr = strings.ToLower(strings.Join(pinyinParts, ""))
		// 提取首字母
		initials = ExtractInitials(strings.Join(pinyinParts, " "))
	} else {
		// 没有中文，可能是纯英文或数字
		// 转换为小写，去除空格和特殊字符，只保留字母和数字
		lowerName := strings.ToLower(name)
		pinyinStr = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(lowerName, "")
		// 提取英文首字母
		initials = extractEnglishInitials(name)
	}

	return pinyinStr, initials
}

// extractEnglishInitials 从英文或混合字符串中提取首字母
// 例如："John Doe" -> "jd", "ABC123" -> "a", "test" -> "t"
func extractEnglishInitials(s string) string {
	if s == "" {
		return ""
	}

	// 按空格分割，提取每个单词的首字母
	parts := strings.Fields(s)
	if len(parts) == 0 {
		// 如果没有空格，取整个字符串的首字母
		if len(s) > 0 {
			firstChar := strings.ToLower(string(s[0]))
			if matched, _ := regexp.MatchString(`^[a-z]$`, firstChar); matched {
				return firstChar
			}
		}
		return ""
	}

	initials := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 取第一个字符（小写）
		firstChar := strings.ToLower(string(part[0]))
		// 只保留字母
		if matched, _ := regexp.MatchString(`^[a-z]$`, firstChar); matched {
			initials = append(initials, firstChar)
		}
	}

	return strings.Join(initials, "")
}
