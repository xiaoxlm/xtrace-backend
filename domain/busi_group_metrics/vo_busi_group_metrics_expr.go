package busi_group_metrics

import (
	"fmt"
	"regexp"
	"strings"
)

type busiGroupMetricsExpr struct {
	originalExpr string
	ibn          string
	hostIPs      []string
}

func newBusiGroupMetricsExpr(originalExpr string, ibn string, hostIPs []string) (*busiGroupMetricsExpr, error) {
	b := &busiGroupMetricsExpr{
		originalExpr: originalExpr,
		ibn:          ibn,
		hostIPs:      hostIPs,
	}

	if err := b.check(); err != nil {
		return nil, err
	}

	return b, nil
}

func (metrics *busiGroupMetricsExpr) check() error {
	if len(metrics.hostIPs) < 1 {
		return fmt.Errorf("hostIPs is empty in BusiGroupMetrics")
	}

	return nil
}

// 解析表达式
func (metrics *busiGroupMetricsExpr) getParsedExpr() string {
	metrics.originalExpr = strings.ReplaceAll(metrics.originalExpr, "$IBN", metrics.ibn)

	// 使用正则匹配花括号内的标签内容
	var tmpExpr string
	{
		re := regexp.MustCompile(`\{([^}]+)\}`)
		matches := re.FindStringSubmatch(metrics.originalExpr)
		if len(matches) > 1 {
			// 获取花括号内的标签内容
			labels := matches[1]

			// 分割标签
			labelPairs := strings.Split(labels, ",")
			var newLabels []string

			// 过滤掉host_ip标签
			for _, pair := range labelPairs {
				if !strings.Contains(pair, "host_ip=") {
					newLabels = append(newLabels, strings.TrimSpace(pair))
				}
			}

			// 重新组合标签
			labelStr := strings.Join(newLabels, ", ")

			// 替换原始表达式中的标签部分
			tmpExpr = re.ReplaceAllString(metrics.originalExpr, "{"+labelStr+", %s}")
		}
	}

	// 匹配
	var hostIPLables string
	{
		for _, ip := range metrics.hostIPs {
			hostIPLables += fmt.Sprintf("%s|", ip)
		}

		hostIPLables = `host_ip=~"` + strings.TrimSuffix(hostIPLables, "|") + `"`
	}

	return fmt.Sprintf(tmpExpr, hostIPLables)
}
