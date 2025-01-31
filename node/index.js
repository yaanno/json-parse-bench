import * as fs from "fs";
import { loadavg } from "os";
import { performance } from "perf_hooks";
import JSONStream from "jsonstream";

const logSystemMetrics = () => {
  const memoryUsage = process.memoryUsage();
  const cpuUsage = process.cpuUsage();
  const loadAverage = loadavg();

  console.log("Memory Usage:", {
    RSS: memoryUsage.rss / 1024 / 1024,
    heapTotal: memoryUsage.heapTotal / 1024 / 1024,
    heapUsed: memoryUsage.heapUsed / 1024 / 1024,
    external: memoryUsage.external / 1024 / 1024,
    arrayBuffers: memoryUsage.arrayBuffers / 1024 / 1024,
  });
  console.log("CPU Usage:", cpuUsage);
  console.log("Load Average:", loadAverage);
};

const processFile = () => {
  const stream = fs.createReadStream('large-file.json', { encoding: 'utf8' });
  const parser = JSONStream.parse('*'); // Adjust the path as needed

  stream.pipe(parser);

  let count = 0;
  parser.on('data', (item) => {
    // Process each item
    // console.log(item.id);
    count++;
  });

  parser.on('end', () => {
    console.log('Finished processing file');
    console.log(`Processed ${count} items`);
  });

  parser.on('error', (err) => {
    console.error('Error parsing JSON:', err);
  });
};

const init = async () => {
  const start = performance.now();
  logSystemMetrics();
  processFile();
  const end = performance.now();
  console.log(`Time taken: ${end - start} milliseconds`);
  logSystemMetrics();
  process.exit();
};

(async () => {
  await init();
})();
