import os
from openai import OpenAI

# Create a custom OpenAI client that points to your local server
client = OpenAI(
    api_key="dummy_api_key",  # The API key doesn't matter for your local server
    base_url="http://localhost:8080/v1"  # Point to your local server
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
        print(chat_completion)
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
