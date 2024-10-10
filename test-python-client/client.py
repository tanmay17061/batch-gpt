import argparse
import sys
from openai import OpenAI

def chat_completion(client, content):
    try:
        chat_completion = client.chat.completions.create(
            messages=[
                {
                    "role": "user",
                    "content": content,
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

def status_single_batch(client, batch_id):
    try:
        response = client.batches.retrieve(batch_id)
        print(f"Batch ID: {response.batch.id}")
        print(f"Status: {response.batch.status}")
        print(f"Created At: {response.batch.created_at}")
        print(f"Expires At: {response.batch.expires_at}")
        print(f"Request Counts: {response.batch.request_counts}")
    except Exception as e:
        print(f"An error occurred: {e}")

def status_all_batches(client):
    try:
        response = client.batches.list()
        if not response.data:
            print("No batches received from the batches API")
            return
        print("All batches:")
        for batch_response in response.data:
            batch = batch_response
            print(f"Batch ID: {batch.id}")
            print(f"Status: {batch.status}")
            print(f"Created At: {batch.created_at}")
            print(f"Expires At: {batch.expires_at}")
            print(f"Request Counts: {batch.request_counts}")
            print("---")
    except Exception as e:
        print(f"An error occurred: {e}")

def main():
    parser = argparse.ArgumentParser(description="Interact with the batch-gpt server.")
    parser.add_argument("--api", choices=['chat_completions', 'status_single_batch', 'status_all_batches'],
                        required=True, help="The API endpoint to call.")
    parser.add_argument("--content", help="The content for chat completion (required for chat_completions).")
    parser.add_argument("--batch_id", help="The batch ID (required for status_single_batch).")
    args = parser.parse_args()

    client = OpenAI(
        api_key="dummy_openai_api_key",
        base_url="http://localhost:8080/v1"
    )

    if args.api == 'chat_completions':
        if not args.content:
            print("Error: --content is required for chat_completions.")
            sys.exit(1)
        chat_completion(client, args.content)
    elif args.api == 'status_single_batch':
        if not args.batch_id:
            print("Error: --batch_id is required for status_single_batch.")
            sys.exit(1)
        status_single_batch(client, args.batch_id)
    elif args.api == 'status_all_batches':
        status_all_batches(client)

if __name__ == "__main__":
    main()
