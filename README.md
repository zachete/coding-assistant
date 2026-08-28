# Coding Assistant

A learning pet-project for an AI-powered coding assistant with a terminal interface, built using Go, Bubble Tea, and OpenAI.

<img width="1280" height="528" alt="coding-agent" src="https://github.com/user-attachments/assets/87635ede-969c-4f0d-8660-1cef9f7a6268" />

## Supported Tools

The assistant has the following tools:

- `read_file`: Read the content of a file given its path.
- `write_file`: Write content to a file given its path and content.

## How to Run

1. Create a `.env` file with your `OPENAI_API_KEY`.
2. Run the program:
   ```bash
   go run cmd/main.go
   ```
