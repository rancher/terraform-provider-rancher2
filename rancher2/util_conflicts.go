package rancher2

func getConflicts(fieldNames []string, fieldName string) []string {
	conflicts := make([]string, 0, len(fieldNames)-1)
	for _, name := range fieldNames {
		if name != fieldName {
			conflicts = append(conflicts, name)
		}
	}
	return conflicts
}
