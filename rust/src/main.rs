use serde_json::Value;
use std::fs::File;
use std::io::BufReader;
use std::time::Instant;
use std::mem;
use sys_info;
use rayon::prelude::*;

fn get_memory_usage() -> Result<usize, sys_info::Error> {
    let mem_info = sys_info::mem_info()?;
    Ok((mem_info.total - mem_info.free) as usize)
}

fn process_large_json(file_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    let start_memory = get_memory_usage()?;
    let start_time = Instant::now();

    let file = File::open(file_path)?;
    let reader = BufReader::new(file);
    let json: Value = serde_json::from_reader(reader)?;
    let mut out = Vec::new();

    if let Some(items) = json.as_array() {
        out = items
            .par_iter()  // Parallel iterator
            .filter_map(|item| {
                item.get("id").map(|id| serde_json::json!({ "id": id }))
            })
            .collect();
    }

    let end_time = start_time.elapsed();
    let end_memory = get_memory_usage()?;
    
    let memory_usage = (end_memory - start_memory) as f64 / 1024.0 / 1024.0; // Convert to MB
    let output_size = out.len();
    let output_memory = mem::size_of_val(&out) as f64 / 1024.0 / 1024.0; // Memory of output vector

    println!(
        "Benchmark Results:\n\
        - Total Items Processed: {}\n\
        - Memory Usage: {:.2} MB\n\
        - Output Vector Memory: {:.2} MB\n\
        - Elapsed Time: {:.4} seconds",
        output_size,
        memory_usage,
        output_memory,
        end_time.as_secs_f64()
    );

    Ok(())
}

fn main() {
    if let Err(e) = process_large_json("large-file.json") {
        eprintln!("Error processing JSON: {}", e);
    }
}
