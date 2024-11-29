import json
import psutil  # type: ignore
import time
import ijson

def process_large_json(file_path):
    # Start benchmarking
    process = psutil.Process()
    start_memory = process.memory_info().rss
    start_cpu = process.cpu_percent(interval=None)
    start_time = time.perf_counter()
    out = []

    with open(file_path, 'rb') as file:
        parser = ijson.items(file, 'item')
        for item in parser:
            obj = {
                "id": item["id"]
            }
            out.append(obj)
    # out = json.dumps(out)
    # End benchmarking
    end_memory = process.memory_info().rss
    end_cpu = process.cpu_percent(interval=None)
    end_time = time.perf_counter()

    memory_usage = (end_memory - start_memory) / (1024 * 1024)  # Convert to MB
    cpu_usage = end_cpu - start_cpu
    elapsed_time = end_time - start_time

    print({
        "out": len(out),
        "memory_usage_mb": memory_usage,
        "cpu_usage_percent": cpu_usage,
        "elapsed_time_seconds": elapsed_time
    })

# for i in range(10):
#     print(i)
process_large_json('large-file.json')
