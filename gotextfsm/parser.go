package gotextfsm

// TextFSMValue represents a field and its corresponding regex.

func Parser(templateFilePath string, dataFilePath string) []Record {
	templateFields, startPatterns := ParseTemplateFile(templateFilePath)
	return ParseCLIOutput(dataFilePath, templateFields, startPatterns)
}
