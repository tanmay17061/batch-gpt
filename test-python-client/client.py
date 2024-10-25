import argparse
import sys
from openai import OpenAI
from datetime import datetime
import time

def format_timestamp(timestamp):
    # Convert Unix timestamp (seconds since epoch) to datetime object using system timezone
    if isinstance(timestamp, int):
        # If timestamp is in milliseconds, convert to seconds
        if timestamp > 1e10:  # Likely in milliseconds
            timestamp = timestamp / 1000
        return datetime.fromtimestamp(timestamp).strftime('%Y-%m-%d %H:%M:%S %Z')
    return "N/A"

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
        print(f"Created: {format_timestamp(chat_completion.created)}")
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
        print(f"Created At: {format_timestamp(response.batch.created_at)}")
        print(f"Expires At: {format_timestamp(response.batch.expires_at)}")
        print(f"Request Counts: {response.batch.request_counts}")
    except Exception as e:
        print(f"An error occurred: {e}")

def status_all_batches(client, status_filter=None):
    try:
        response = client.batches.list()
        if not response.data:
            print("No batches received from the batches API")
            return

        # Filter batches based on status
        filtered_batches = response.data
        if status_filter:
            if status_filter == 'completed':
                filtered_batches = [batch for batch in response.data if batch.status == 'completed']
            elif status_filter == 'not_completed':
                filtered_batches = [batch for batch in response.data if batch.status != 'completed']

        if not filtered_batches:
            print(f"No batches found with filter: {status_filter}")
            return

        print(f"All batches{' (filtered: ' + status_filter + ')' if status_filter else ''}:")
        for batch in filtered_batches:
            print(f"Batch ID: {batch.id}")
            print(f"Status: {batch.status}")
            print(f"Created At: {format_timestamp(batch.created_at)}")
            print(f"Expires At: {format_timestamp(batch.expires_at)}")
            print(f"Request Counts: {batch.request_counts}")
            print("---")

        # Print summary
        print("\nSummary:")
        print(f"Total batches shown: {len(filtered_batches)}")
        if status_filter:
            print(f"Filter applied: {status_filter}")
            total_batches = len(response.data)
            print(f"Total batches before filter: {total_batches}")

    except Exception as e:
        print(f"An error occurred: {e}")

def main():
    parser = argparse.ArgumentParser(description="Interact with the batch-gpt server.")
    parser.add_argument("--api", choices=['chat_completions', 'status_single_batch', 'status_all_batches'],
                        required=True, help="The API endpoint to call.")
    parser.add_argument("--content", help="The content for chat completion (required for chat_completions).")
    parser.add_argument("--batch_id", help="The batch ID (required for status_single_batch).")
    parser.add_argument("--status_filter",
                       choices=['completed', 'not_completed'],
                       help="Filter batches by status (only for status_all_batches).")
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
        status_all_batches(client, args.status_filter)

if __name__ == "__main__":
    main()
