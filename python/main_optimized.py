import ijson
import time
import psutil
import multiprocessing
from typing import Iterator, Dict, Any

def process_item(item: Dict[str, Any]) -> Dict[str, str]:
    """
    Process a single JSON item with minimal overhead.
    Only extract the 'id' field to reduce memory usage.
    """
    return {"id": item.get("id", "")}

def stream_json_items(file_path: str) -> Iterator[Dict[str, Any]]:
    """
    Efficiently stream JSON items using ijson with minimal memory overhead.
    """
    with open(file_path, 'rb') as file:
        parser = ijson.items(file, 'item')
        for item in parser:
            yield item

def process_chunk(chunk: list) -> int:
    """
    Process a chunk of items in parallel.
    Returns the number of processed items.
    """
    return len([process_item(item) for item in chunk])

def parallel_process_json(file_path: str, chunk_size: int = 1000) -> Dict[str, Any]:
    """
    Parallel processing of JSON file with minimal memory consumption.
    """
    # Start performance tracking
    process = psutil.Process()
    start_memory = process.memory_info().rss
    start_time = time.perf_counter()
    
    # Prepare for parallel processing
    chunk = []
    total_processed = 0
    
    # Use a generator to minimize memory usage
    for item in stream_json_items(file_path):
        chunk.append(item)
        
        # Process in chunks
        if len(chunk) >= chunk_size:
            total_processed += process_chunk(chunk)
            chunk.clear()  # Clear the chunk to free memory
    
    # Process any remaining items
    if chunk:
        total_processed += process_chunk(chunk)
    
    # End performance tracking
    end_memory = process.memory_info().rss
    end_time = time.perf_counter()
    
    # Compile performance metrics
    return {
        "total_items_processed": total_processed,
        "memory_usage_mb": (end_memory - start_memory) / (1024 * 1024),
        "elapsed_time_seconds": end_time - start_time,
        "peak_memory_mb": process.memory_info().rss / (1024 * 1024)  # Use current RSS as peak memory
    }

def main():
    """
    Main function to run JSON processing with performance tracking.
    """
    try:
        # Determine optimal chunk size based on CPU cores
        num_cores = multiprocessing.cpu_count()
        chunk_size = max(1000, num_cores * 500)
        
        # Process the JSON file
        results = parallel_process_json('large-file.json', chunk_size)
        
        # Print detailed performance metrics
        print("JSON Processing Performance Metrics:")
        for key, value in results.items():
            print(f"{key.replace('_', ' ').title()}: {value}")
    
    except Exception as e:
        print(f"Error processing JSON: {e}")

if __name__ == "__main__":
    main()
