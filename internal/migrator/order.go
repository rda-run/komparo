package migrator

func OrderChanges(changes []SchemaChange) []SchemaChange {
	var toCreate, toDrop, toAlter, toRecreate []SchemaChange

	for _, c := range changes {
		switch c.ChangeType {
		case ChangeCreate:
			toCreate = append(toCreate, c)
		case ChangeDrop:
			toDrop = append(toDrop, c)
		case ChangeAlter:
			toAlter = append(toAlter, c)
		case ChangeRecreate:
			toRecreate = append(toRecreate, c)
		}
	}

	orderCreate := []string{"Extension", "Enum", "Sequence", "Table", "Column",
		"Constraint", "Index", "View", "MaterializedView", "Function", "Trigger"}

	orderedCreate := sortByTypeOrder(toCreate, orderCreate)

	orderDrop := reverse(orderCreate)
	orderedDrop := sortByTypeOrder(toDrop, orderDrop)

	orderedAlter := sortByTypeOrder(toAlter, orderCreate)
	orderedRecreate := sortByTypeOrder(toRecreate, orderCreate)

	result := append([]SchemaChange{}, orderedDrop...)
	result = append(result, orderedAlter...)
	result = append(result, orderedRecreate...)
	result = append(result, orderedCreate...)

	return result
}

func sortByTypeOrder(changes []SchemaChange, order []string) []SchemaChange {
	orderMap := make(map[string]int)
	for i, t := range order {
		orderMap[t] = i
	}

	result := make([]SchemaChange, len(changes))
	copy(result, changes)

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if orderMap[result[i].ObjectType] > orderMap[result[j].ObjectType] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func reverse(s []string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}
