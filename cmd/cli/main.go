package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aishell/internal/executor"
	"aishell/internal/llm"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	execute     bool
	verbose     bool
)

func init() {
	// Get the executable's directory
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("Warning: Could not determine executable path: %v\n", err)
		return
	}
	exPath := filepath.Dir(ex)

	// Try to load .env from multiple possible locations
	locations := []string{
		".env",                              // Current directory
		filepath.Join(exPath, ".env"),       // Executable's directory
		filepath.Join(exPath, "../.env"),    // One level up
		filepath.Join(exPath, "../../.env"), // Two levels up
	}

	var loaded bool
	for _, loc := range locations {
		if err := godotenv.Load(loc); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		fmt.Println("Warning: Could not load .env file from any location")
	}

	// Verify the API key is available
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		fmt.Println("Error: OPENROUTER_API_KEY is not set in .env file")
		fmt.Println("Please create a .env file in the project root with your OpenRouter API key:")
		fmt.Println("OPENROUTER_API_KEY=your-api-key-here")
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "aishell [flags] [query]",
		Short:   "AI-powered shell command generator",
		Version: "1.0.0",
		Long: `aishell is an AI-powered command-line tool that converts natural language 
descriptions into shell commands. It uses advanced language models to understand
your intent and generate the appropriate command.

You can use it in two ways:
1. Direct mode: Simply provide your query as an argument
2. Interactive mode: Use the -i flag to enter an interactive shell

The tool will generate shell commands based on your natural language input and
can optionally execute them directly.`,
		Example: `  # Generate a command to find large files
  aishell "find files larger than 100MB"

  # Start interactive mode
  aishell -i

  # Generate and automatically execute the command
  aishell -e "show system memory usage"

  # Show detailed output
  aishell -v "list all pdf files"`,
		Run:  runCommand,
		Args: cobra.ArbitraryArgs,
	}

	// Add flags with improved descriptions
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false,
		"Start interactive mode with a command prompt")
	rootCmd.PersistentFlags().BoolVarP(&execute, "execute", "e", false,
		"Automatically execute the generated command without confirmation")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Display detailed output including query processing information")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand(cmd *cobra.Command, args []string) {
	// If interactive flag is set, run interactive mode regardless of arguments
	if interactive {
		runInteractiveMode()
		return
	}

	// Handle direct query mode
	if len(args) == 0 {
		fmt.Println("Please provide a natural language query or use -i for interactive mode")
		os.Exit(1)
	}

	// Process the query from arguments
	query := strings.Join(args, " ")
	processQuery(query)
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n🤖 Interactive mode - Type 'exit' or 'quit' to end the session")

	for {
		fmt.Print("\n💭 Enter query: ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)

		if query == "exit" || query == "quit" {
			fmt.Println("Goodbye! 👋")
			return
		}

		if query == "" {
			continue
		}

		processQuery(query)
	}
}

func processQuery(query string) {
	if verbose {
		fmt.Printf("\n🔄 Generating command for: %s\n", query)
	}

	command, err := llm.GenerateCommand(query)
	if err != nil {
		fmt.Printf("❌ Error generating command: %v\n", err)
		return
	}

	fmt.Printf("\n📝 Generated command:\n%s\n", command)

	if execute {
		executeCommand(command)
		return
	}

	promptAndExecute(command)
}

func promptAndExecute(command string) {
	fmt.Print("\n⚡ Run this command? [y]es/[n]o/[e]dit/[q]uit: ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	switch confirm {
	case "y", "yes":
		executeCommand(command)
	case "e", "edit":
		fmt.Print("🖊️  Edit command: ")
		editedCommand, _ := reader.ReadString('\n')
		editedCommand = strings.TrimSpace(editedCommand)
		if editedCommand != "" {
			executeCommand(editedCommand)
		}
	case "q", "quit":
		fmt.Println("Goodbye! 👋")
		os.Exit(0)
	default:
		fmt.Println("Command not executed.")
	}
}

func executeCommand(command string) {
	if verbose {
		fmt.Printf("\n🚀 Executing command: %s\n", command)
	}

	output, err := executor.ExecuteCommand(command)
	if err != nil {
		fmt.Printf("❌ Error executing command: %v\n", err)
		return
	}

	if output == "" {
		fmt.Println("✅ Command executed successfully (no output)")
	} else {
		fmt.Printf("\n📄 Output:\n%s\n", output)
	}
}
