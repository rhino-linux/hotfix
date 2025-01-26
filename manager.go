package main

import (
	"fmt"
	"os"
	"encoding/json"
	"time"
	"flag"
	"regexp"
)

type Item struct { // Hotfix details
	Script string `json:"script"`
	Description string `json:"description"`
	Published string `json:"published"`
}

type Items map[string]Item // Type for JSON map

// Read JSON file
func readHotfixes(filename string) (Items, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read file
	var items Items
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&items)
	if err != nil && err.Error() != "EOF" { // Allows EOF errors
		return nil, err
	}
	return items, nil
}

// Write to file
func writeHotfixes(filename string, items Items) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(items)
	if err != nil {
		return err
	}
	return nil
}

// Add items
func addHotfix(filename, name, script, description string) error {
    // Validate the title field
	if err := validateField(name, "Title"); err != nil {
		return err
	}
	items, err := readHotfixes(filename)
	if err != nil && !os.IsNotExist (err) {
		return err
	}
	// Check if script file exists
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("Script file does not exist")
	}
	// If file != exist or is empty
	if items == nil {
		items = make(Items)
	}
	// Get date
	date := time.Now().Format("2006-01-02") // Format: YYYY-MM-DD
	items[name] = Item{Script: script, Description: description, Published: date} // Add the item to a map.
	// Write to file
	return writeHotfixes(filename, items)
}

// Validation
func validateField(field, fieldName string) error {
    valid := `^[A-Za-z0-9\-~._+()]+$`
    re := regexp.MustCompile(valid)
    if !re.MatchString(field) {
        return fmt.Errorf("%s contains invalid characters. Allowed characters are: [A-z], [0-9], [-~._+()]", fieldName)
    }
    return nil
}

// Remove hotfixes
func removeHotfix(filename, name string) error {
	// Validate the title field
	if err := validateField(name, "Title"); err != nil {
		return err
	}

	items, err := readHotfixes(filename)
	if err != nil {
		return err
	}
	// Check if item exists
	if _, exists := items[name]; !exists {
		return fmt.Errorf("Item does not exist")
	}
	// Remove item
	delete(items, name)
	// Write to file
	return writeHotfixes(filename, items)
}

// List hotfixes
func listHotfixes(filename string) {
	items, err := readHotfixes(filename)
	if err != nil {
		fmt.Println("Error reading hotfixes:", err)
		return
	}
	// Print hotfixes
	for name, item := range items {
		fmt.Println("Title: ", name)
		fmt.Println("Script: ", item.Script)
		fmt.Println("Description: ", item.Description)
		fmt.Println("Published: ", item.Published)
		fmt.Println()
	}
}

func main() {
	filename := "hotfixes.json"

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run manager.go <command> [options]")
        fmt.Println("  add      Add a new hotfix")
		fmt.Println("  remove   Remove an existing hotfix")
		fmt.Println("  list     List all hotfixes")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
        title := addCmd.String("t", "", "Title of the hotfix (required)")
        description := addCmd.String("d", "", "Description of the hotfix (required)")
        script := addCmd.String("s", "", "Path to the hotfix script (required)")
        addCmd.Parse(os.Args[2:])
        if *title == "" || *description == "" || *script == "" {
			fmt.Println("Error: all flags (-t, -d, -s) are required")
			addCmd.Usage()
			return
		}
        err := addHotfix(filename, *title, *script, *description)
        if err != nil {
            fmt.Println("Error adding hotfix:", err)
            return
        } else {
            fmt.Println("Hotfix added successfully!")
        }
	case "remove":
		if len(os.Args) < 3 {
            fmt.Println("Error: hotfix title is required")
            return
        }
        name := os.Args[2]
        err := removeHotfix(filename, name)
        if err != nil {
            fmt.Println("Error removing hotfix:", err)
            return
        } else {
            fmt.Println("Hotfix removed successfully!")
        }
	case "list":
		listHotfixes(filename)
	default:
		fmt.Println("Invalid operation.")
	}
}
