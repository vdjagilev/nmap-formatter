package formatter

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// filterExpr filters NMAPRun.Hosts by given expression
func filterExpr(r NMAPRun, code string) (NMAPRun, error) {
	program, err := expr.Compile(
		fmt.Sprintf("filter(Host, { %s })", code),
		expr.Env(r),
	)

	if err != nil {
		return r, err
	}

	output, err := expr.Run(program, r)
	if err != nil {
		return r, err
	}

	hosts, err := convertToHosts(output)

	if err != nil {
		return r, err
	}

	r.Host = hosts
	return r, nil
}

// convertToHosts converts output from expression engine to []Host
func convertToHosts(output interface{}) ([]Host, error) {
	outputInterfaces, ok := output.([]interface{})
	if !ok {
		return nil, fmt.Errorf("output is not []interface{}")
	}

	hosts := make([]Host, len(outputInterfaces))
	for i, v := range outputInterfaces {
		host, ok := v.(Host)
		if !ok {
			return nil, fmt.Errorf("element is not Host")
		}
		hosts[i] = host
	}

	return hosts, nil
}
