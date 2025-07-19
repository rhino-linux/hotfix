package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"encoding/json"
	"time"
	"flag"
	"regexp"
)

type Item struct { // Hotfix details
	Script string `json:"script"`
	Description string `json:"description"`
	Published string `json:"published"`
	Archived string `json:"archived"`
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
func addHotfix(filename, name, script, description, archived string) error {
    // Validate the title field
	if err := validateField(name, "Title"); err != nil {
		return err
	}
	if err := validateArchived(archived); err != nil {
		return err
	}
	items, err := readHotfixes(filename)
	if err != nil && !os.IsNotExist (err) {
		return err
	}
	if items == nil {
		items = make(Items)
	}
	item, exists := items[name]
	var filePath, activePath, archivePath, targetPath string
	activePath = "scripts/active/" + script
	if !exists {
		if archived == "false" {
			filePath = activePath
		} else {
			filePath = "scripts/archived/" + archived + "/" + script
		}
	} else {
		if ((item.Archived == "false" && item.Script == script) || (archived == "false" && item.Script != script)) {
			filePath = activePath
		} else if item.Archived != "false" {
			filePath = "scripts/archived/" + item.Archived + "/" + script
		} else {
			filePath = "scripts/archived/" + archived + "/" + script
		}
	}
	// Check if script file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("'%s' does not exist", filePath)
	}

	if (item.Archived == "false" && archived != "false") {
		targetPath = "scripts/archived/" + archived + "/" + script
	} else if (item.Archived != "false" && archived == "false") {
		targetPath = activePath
	}

	if exists {
		if item.Script == script {
			if ((item.Archived == "false" && archived != "false") || (item.Archived != archived && archived != "false")) {
				targetPath = "scripts/archived/" + archived + "/" + script
			} else if (item.Archived != "false" && archived == "false") {
				targetPath = activePath
			}
			if err := os.MkdirAll(strings.TrimSuffix(targetPath, script), 0755); err != nil {
			    return fmt.Errorf("failed to create directory: %v", err)
			}
			if err := os.Rename(filePath, targetPath); err != nil {
				return fmt.Errorf("failed to move script to archived folder: %v", err)
			}
		} else {
			if item.Archived == "false" {
				archivePath = "scripts/active/" + item.Script
			} else {
				archivePath = "scripts/archived/" + item.Archived + "/" + item.Script
			}
			if err := os.Remove(archivePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete old script file: %v", err)
			}
		}
	}

	// Get date
	date := time.Now().Format("2006-01-02") // Format: YYYY-MM-DD
	items[name] = Item{Script: script, Description: description, Published: date, Archived: archived} // Add the item to a map.
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

func validateArchived(archived string) error {
    valid := `^(false|\d{4}\.\d+)$`
    re := regexp.MustCompile(valid)
    if !re.MatchString(archived) {
        return fmt.Errorf("archived must be either 'false' or in the format YYYY.VERNUM (e.g., 2023.1)")
    }
    return nil
}

// Remove hotfixes
func removeHotfix(filename, name, archived string) error {
	// Validate the title field
	if err := validateField(name, "Title"); err != nil {
		return err
	}
	if err := validateArchived(archived); err != nil {
		return err
	}
	items, err := readHotfixes(filename)
	if err != nil {
		return err
	}
	// Check if item exists
	item, exists := items[name]
	if !exists {
		return fmt.Errorf("Item does not exist")
	}
	var filePath, archivePath string
	if item.Archived == "false" {
		filePath = "scripts/active/" + item.Script
	} else {
		filePath = "scripts/archived/" + item.Archived + "/" + item.Script
	}
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("Script file does not exist")
	}
	if archived == "false" {
		// Remove script and item
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete script file: %v", err)
		}
		delete(items, name)
	} else {
		if item.Archived != archived {
			archivePath = "scripts/archived/" + archived + "/" + item.Script
			err := os.MkdirAll(strings.TrimSuffix(archivePath, item.Script), 0755)
			if err != nil {
			    return fmt.Errorf("failed to create directory: %v", err)
			}
			if err := os.Rename(filePath, archivePath); err != nil {
				return fmt.Errorf("failed to move script to archived folder: %v", err)
			}
		}
		item.Archived = archived
		items[name] = item
	}
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
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Print hotfixes
	for _, name := range keys {
		item := items[name]
		fmt.Println("Title: ", name)
		if item.Archived == "false" {
			fmt.Println("Script: ", "scripts/active/" + item.Script)
		} else {
			fmt.Println("Script: ", "scripts/archived/" + item.Archived + "/" + item.Script)
		}
		fmt.Println("Description: ", item.Description)
		fmt.Println("Published: ", item.Published)
		fmt.Println("Archived: ", item.Archived)
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
        title := addCmd.String("t", "", "Title of the hotfix [required]")
        description := addCmd.String("d", "", "Description of the hotfix [required]")
        script := addCmd.String("s", "", "Path to the hotfix script [required]")
        archived := addCmd.String("a", "false", "Version the hotfix was deprecated [optional]")
        addCmd.Parse(os.Args[2:])
        if *title == "" || *description == "" || *script == "" || *archived == "" {
			fmt.Println("Error: all flags (-t, -d, -s, -a) are required or cannot be empty")
			addCmd.Usage()
			return
		}
        err := addHotfix(filename, *title, *script, *description, *archived)
        if err != nil {
            fmt.Println("Error adding/updating hotfix:", err)
            return
        } else {
            fmt.Println("Hotfix added/updated successfully!")
        }
	case "remove":
		rmCmd := flag.NewFlagSet("remove", flag.ExitOnError)
		title := rmCmd.String("t", "", "Title of the hotfix [required]")
		archived := rmCmd.String("a", "false", "Archive instead of remove; Version the hotfix was deprecated [optional]")
		rmCmd.Parse(os.Args[2:])
        if *title == "" || *archived == "" {
			fmt.Println("Error: all flags (-t, -a) are required or cannot be empty")
			rmCmd.Usage()
			return
		}
        err := removeHotfix(filename, *title, *archived)
        if err != nil {
            fmt.Println("Error removing/archiving hotfix:", err)
            return
        } else {
			if *archived != "false" {
				fmt.Println("Hotfix archived successfully!")
			} else {
				fmt.Println("Hotfix removed successfully!")
			}
		}
	case "list":
		listHotfixes(filename)
	default:
		fmt.Println("Invalid operation.")
	}
}
