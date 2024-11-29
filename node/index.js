import { promises as fs } from "fs";
import { loadavg } from "os";
import { performance } from "perf_hooks";

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

const processFile = async () => {
  try {
    const data = await fs.readFile("large-file.json", "utf8");
    const json = JSON.parse(data);
    const out = json.map((item) => ({ id: item.id }));
    console.log(out.length);
  } catch (err) {
    console.error(err);
  }
};

const init = async () => {
  const start = performance.now();
  logSystemMetrics();
  await processFile().then(() => {
    const end = performance.now();
    console.log(`Time taken: ${end - start} milliseconds`);
    logSystemMetrics();
  });
};

(async () => {
  // for (let index = 0; index < 10; index++) {
  //   console.log(index);
  await init();
  // }
})();
