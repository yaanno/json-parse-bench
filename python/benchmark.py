import json
import sys
import time
import psutil
import os

def process_json(file_path):
    start_time = time.time()
    start_memory = psutil.Process(os.getpid()).memory_info().rss / (1024 * 1024)

    with open(file_path, 'r') as f:
        data = json.load(f)
    
    processed_data = [{'id': item['id']} for item in data if 'id' in item]

    end_time = time.time()
    end_memory = psutil.Process(os.getpid()).memory_info().rss / (1024 * 1024)

    print(json.dumps({
        'out': len(processed_data),
        'memory_usage_mb': end_memory - start_memory,
        'elapsed_time_seconds': end_time - start_time
    }))

if __name__ == "__main__":
    process_json('large-file.json')
