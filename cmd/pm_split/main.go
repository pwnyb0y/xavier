package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultInputFile = "/Users/jm/Downloads/OpenAI.postman_collection.json"
)

func splitCollectionByFolder(collectionFile string, outputLocation string, targetFolder string) {
	file, err := os.Open(collectionFile)
	if err != nil {
		fmt.Printf("Failed to open collection file: %v\n", err)
		return
	}
	defer file.Close()

	var collection map[string]interface{}
	err = json.NewDecoder(file).Decode(&collection)
	if err != nil {
		fmt.Printf("Failed to decode collection JSON: %v\n", err)
		return
	}

	if outputLocation == "" {
		outputLocation, err = os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current working directory: %v\n", err)
			return
		}
	}

	searchAllFolders := false
	if targetFolder == "" {
		targetFolder = getTargetFolderName(collection)
		searchAllFolders = true
	}

	folderPath := filepath.Join(outputLocation, targetFolder)
	err = os.MkdirAll(folderPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create output folder: %v\n", err)
		return
	}

	if items, ok := collection["item"].([]interface{}); ok {
		findTargetFolder(items, folderPath, targetFolder, searchAllFolders)
	}
}

func findTargetFolder(items []interface{}, outputLocation string, targetFolder string, searchAllFolders bool) {
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			log.Printf("Target Folder: %v\n", targetFolder)
			if searchAllFolders || itemMap["name"] == targetFolder {
				fmt.Printf("Found target folder: %v\n", itemMap["name"])
				extractRequestsFromFolder(itemMap, outputLocation)
			}

			if nestedItems, ok := itemMap["item"].([]interface{}); ok {
				findTargetFolder(nestedItems, outputLocation, targetFolder, searchAllFolders)
			}
		}
	}
}

func extractRequestsFromFolder(folder map[string]interface{}, outputLocation string) {
	log.Printf("Extracting folder: %v\n", folder["name"])
	if items, ok := folder["item"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				log.Printf("Extracting item: %v\n", itemMap["name"])
				if _, ok := itemMap["request"]; ok {
					saveRequestToFile(itemMap, outputLocation)
				} else if _, ok := itemMap["item"].([]interface{}); ok {
					extractRequestsFromFolder(itemMap, outputLocation)
				}
			}
		}
	}
}

func saveRequestToFile(item map[string]interface{}, folderPath string) {
	log.Printf("Saving request: %v\n", item["name"])
	request := item["request"].(map[string]interface{})
	response := item["response"]
	references := item["item"]

	combinedInfo := map[string]interface{}{
		"request":    request,
		"response":   response,
		"references": references,
	}

	fileName := getRequestFileName(request)
	filePath := filepath.Join(folderPath, fileName)

	data, err := json.MarshalIndent(combinedInfo, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}
}

func getRequestFileName(request map[string]interface{}) string {
	url := request["url"].(map[string]interface{})
	path := url["path"].([]interface{})
	name := path[len(path)-1]
	return fmt.Sprintf("%v.json", name)
}

func getTargetFolderName(collection map[string]interface{}) string {
	name := collection["info"].(map[string]interface{})["name"]
	return fmt.Sprintf("%v", name)
}

func main() {
	collectionFile := flag.String("collection", defaultInputFile, "Path to the Postman collection JSON file")
	outputLocation := flag.String("output", "", "Location to output the split collection. Defaults to the current directory")
	targetFolder := flag.String("folder", "", "Name of the folder to extract. If not provided, the collection name is used")
	flag.Parse()

	splitCollectionByFolder(*collectionFile, *outputLocation, *targetFolder)
}
