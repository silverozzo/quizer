package model

type Input struct {
	Type  string
	Name  string
	Value string
}

func Fields(lst *[]Input) map[string]string {
	res := make(map[string]string)

	for _, vl := range *lst {
		switch vl.Type {
		case "text":
			res[vl.Name] = "test"

		case "select":
			fallthrough

		case "radio":
			if len(vl.Value) > len(res[vl.Name]) {
				res[vl.Name] = vl.Value
			}
		}
	}

	return res
}
