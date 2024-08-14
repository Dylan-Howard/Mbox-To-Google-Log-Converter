package main

const IMPORT_DIRECTORY = "../import"
const EXPORT_DIRECTORY = "../export"

func main() {
	generator := MboxGenerator {
		ImportDirectory: IMPORT_DIRECTORY,
		ExportDirectory: EXPORT_DIRECTORY,
	};

	boxes := generator.GetMboxes();

	// fmt.Println(len(boxes))

	generator.GenerateLogs(boxes);
}
