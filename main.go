package main

import (
	"bufio"
	"fmt"
	"ollama-cli-teacher/cache"
	"ollama-cli-teacher/ollama"
	"ollama-cli-teacher/ui"
	"os"
	"strings"
)

func main() {
	client := ollama.NewOllamaClient("http://localhost:11434")
	cacheManager := cache.NewCacheManager("categories_cache.json")
	reader := bufio.NewReader(os.Stdin)

	for {
		ui.ClearScreen()
		ui.PrintHeader("Welcome to the Learning CLI!")
		fmt.Print("Enter the topic you want to learn about: ")
		topic, _ := reader.ReadString('\n')
		topic = strings.TrimSpace(topic)

		var categories string
		var err error

		if cacheManager.HasTopic(topic) {
			categories, err = cacheManager.GetCategories(topic)
			if err != nil {
				fmt.Printf("Error loading cached categories: %v\n", err)
				return
			}
		} else {
			prompt := fmt.Sprintf("List the main categories related to learning '%s'.", topic)
			categories, err = client.GenerateResponse(prompt)
			if err != nil {
				fmt.Printf("Error generating response: %v\n", err)
				return
			}
			cacheManager.SaveCategories(topic, categories)
		}

		for {
			ui.ClearScreen()
			ui.PrintHeader("Here are the categories:")
			fmt.Println(categories)

			fmt.Print("Enter the number of the category you want to dive into (or type 'back' to choose a new topic): ")
			categoryChoice, _ := reader.ReadString('\n')
			categoryChoice = strings.TrimSpace(categoryChoice)

			if strings.ToLower(categoryChoice) == "back" {
				break
			}

			ui.ClearScreen()
			prompt := fmt.Sprintf("Provide more in-depth information on category %s related to '%s'.", categoryChoice, topic)
			response, err := client.GenerateResponse(prompt)
			if err != nil {
				fmt.Printf("Error generating response: %v\n", err)
				return
			}
			ui.PrintHeader("Here's more detailed information:")
			fmt.Println(response)

			fmt.Println("Now you can ask any questions related to this topic.")
			for {
				fmt.Print("Enter your question (or type 'back' to select another category, 'exit' to quit): ")
				question, _ := reader.ReadString('\n')
				question = strings.TrimSpace(question)

				if strings.ToLower(question) == "back" {
					break
				} else if strings.ToLower(question) == "exit" {
					os.Exit(0)
				}

				prompt = fmt.Sprintf("Answer the following question about '%s': %s", topic, question)
				response, err = client.GenerateResponse(prompt)
				if err != nil {
					fmt.Printf("Error generating response: %v\n", err)
					return
				}
				ui.PrintHeader("Answer:")
				fmt.Println(response)
			}
		}
	}
}
