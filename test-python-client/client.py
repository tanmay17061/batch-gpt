import os
from openai import OpenAI

# Create a custom OpenAI client that points to the batch-gpt server.
client = OpenAI(
    api_key="dummy_api_key",
    base_url="http://localhost:8080/v1"
)

def main():
    try:
        chat_completion = client.chat.completions.create(
            messages=[
                {
                    "role": "user",
                    "content": "Say this is a test",
                }
            ],
            model="gpt-3.5-turbo",
        )
        print("Response from server:")
        print(f"ID: {chat_completion.id}")
        print(f"Object: {chat_completion.object}")
        print(f"Created: {chat_completion.created}")
        print(f"Model: {chat_completion.model}")
        for choice in chat_completion.choices:
            print(f"Choice {choice.index}:")
            print(f"  Role: {choice.message.role}")
            print(f"  Content: {choice.message.content}")
            print(f"  Finish Reason: {choice.finish_reason}")
        print("Usage:")
        print(f"  Prompt Tokens: {chat_completion.usage.prompt_tokens}")
        print(f"  Completion Tokens: {chat_completion.usage.completion_tokens}")
        print(f"  Total Tokens: {chat_completion.usage.total_tokens}")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
