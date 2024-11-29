use serde_json::Value;
use std::fs::File;
use std::io::BufReader;
use std::time::Instant;
use sysinfo::{Pid, ProcessRefreshKind, ProcessesToUpdate, System};

fn process_large_json(file_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    // Start benchmarking
    let mut system = System::new_all();
    system.refresh_all();
    std::thread::sleep(sysinfo::MINIMUM_CPU_UPDATE_INTERVAL);
    // Refresh CPU usage to get actual value.
    system.refresh_processes_specifics(
        ProcessesToUpdate::All,
        true,
        ProcessRefreshKind::new().with_cpu(),
    );
    let process = system
        .process(Pid::from_u32(std::process::id()))
        .ok_or("Process not found")?;
    let start_memory = process.memory();
    let start_cpu = process.cpu_usage();
    let start_time = Instant::now();

    let file = File::open(file_path)?;
    let reader = BufReader::new(file);
    let json: Value = serde_json::from_reader(reader)?;
    let mut out = Vec::new();

    if let Some(items) = json.as_array() {
        for item in items {
            if let Some(id) = item.get("id") {
                out.push(serde_json::json!({ "id": id }));
            }
        }
    }

    // End benchmarking
    let end_memory = process.memory();
    let end_cpu = process.cpu_usage();
    let elapsed_time = start_time.elapsed();

    let memory_usage = (end_memory - start_memory) as f64 / 1024.0; // Convert to MB
    let cpu_usage = end_cpu - start_cpu;

    println!(
        "Output length: {}, Memory usage: {:.2} MB, CPU usage: {:.2}%, Elapsed time: {:.2?} seconds",
        out.len(),
        memory_usage,
        cpu_usage,
        elapsed_time
    );

    Ok(())
}

fn main() {
    if let Err(e) = process_large_json("large-file.json") {
        eprintln!("Error processing JSON: {}", e);
    }
}
