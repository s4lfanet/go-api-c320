import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend } from 'k6/metrics';

// Custom metrics
const errors = new Counter('errors');
const timeouts = new Counter('timeouts');
const successRate = new Trend('success_rate');

// Test scenarios
export const options = {
  scenarios: {
    // Scenario 1: Ramp up test (find breaking point)
    ramp_up: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 4 },   // Ramp up to 4 users
        { duration: '1m', target: 4 },    // Stay at 4 users
        { duration: '30s', target: 8 },   // Ramp up to 8 users
        { duration: '1m', target: 8 },    // Stay at 8 users
        { duration: '30s', target: 0 },   // Ramp down
      ],
      gracefulRampDown: '30s',
    },

    // Scenario 2: Spike test (sudden load)
    // spike: {
    //   executor: 'ramping-vus',
    //   startVUs: 0,
    //   stages: [
    //     { duration: '10s', target: 0 },
    //     { duration: '10s', target: 20 },  // Sudden spike
    //     { duration: '1m', target: 20 },
    //     { duration: '10s', target: 0 },
    //   ],
    //   startTime: '5m', // Start after ramp_up
    // },
  },

  thresholds: {
    http_req_duration: ['p(95)<60000'], // 95% requests should complete under 60s (cold cache scenario)
    http_req_failed: ['rate<0.1'],      // Error rate should be less than 10%
    errors: ['count<30'],                // Less than 30 total errors (higher for cold cache)
  },
};

// All 32 board/pon endpoints
const urls = [
  // Board 1
  'http://localhost:8081/api/v1/board/1/pon/1',
  'http://localhost:8081/api/v1/board/1/pon/2',
  'http://localhost:8081/api/v1/board/1/pon/3',
  'http://localhost:8081/api/v1/board/1/pon/4',
  'http://localhost:8081/api/v1/board/1/pon/5',
  'http://localhost:8081/api/v1/board/1/pon/6',
  'http://localhost:8081/api/v1/board/1/pon/7',
  'http://localhost:8081/api/v1/board/1/pon/8',
  'http://localhost:8081/api/v1/board/1/pon/9',
  'http://localhost:8081/api/v1/board/1/pon/10',
  'http://localhost:8081/api/v1/board/1/pon/11',
  'http://localhost:8081/api/v1/board/1/pon/12',
  'http://localhost:8081/api/v1/board/1/pon/13',
  'http://localhost:8081/api/v1/board/1/pon/14',
  'http://localhost:8081/api/v1/board/1/pon/15',
  'http://localhost:8081/api/v1/board/1/pon/16',
  // Board 2
  'http://localhost:8081/api/v1/board/2/pon/1',
  'http://localhost:8081/api/v1/board/2/pon/2',
  'http://localhost:8081/api/v1/board/2/pon/3',
  'http://localhost:8081/api/v1/board/2/pon/4',
  'http://localhost:8081/api/v1/board/2/pon/5',
  'http://localhost:8081/api/v1/board/2/pon/6',
  'http://localhost:8081/api/v1/board/2/pon/7',
  'http://localhost:8081/api/v1/board/2/pon/8',
  'http://localhost:8081/api/v1/board/2/pon/9',
  'http://localhost:8081/api/v1/board/2/pon/10',
  'http://localhost:8081/api/v1/board/2/pon/11',
  'http://localhost:8081/api/v1/board/2/pon/12',
  'http://localhost:8081/api/v1/board/2/pon/13',
  'http://localhost:8081/api/v1/board/2/pon/14',
  'http://localhost:8081/api/v1/board/2/pon/15',
  'http://localhost:8081/api/v1/board/2/pon/16',
];

export default function () {
  // Each VU randomly picks URLs to simulate real usage
  const url = urls[Math.floor(Math.random() * urls.length)];

  const params = {
    timeout: '120s', // Increased client timeout to allow slow SNMP queries (cold cache scenario)
    tags: { name: url.split('/').slice(-2).join('/') }, // Tag by board/pon
  };

  const response = http.get(url, params);

  // Detailed checks
  const checkResult = check(response, {
    'status is 200': (r) => r.status === 200,
    'status is not 408 (timeout)': (r) => r.status !== 408,
    'status is not 429 (rate limit)': (r) => r.status !== 429,
    'status is not 500': (r) => r.status !== 500,
    'response time < 90s': (r) => r.timings.duration < 90000, // Increased for cold cache scenario
    'response has data': (r) => r.body && r.body.length > 0,
  });

  // Track errors
  if (!checkResult) {
    errors.add(1);
  }

  // Track timeouts specifically
  if (response.status === 408 || response.timings.duration > 90000) {
    timeouts.add(1);
    console.log(`⚠️  TIMEOUT on ${url} - Duration: ${response.timings.duration}ms - Status: ${response.status}`);
  }

  // Track rate limits
  if (response.status === 429) {
    console.log(`⚠️  RATE LIMITED on ${url}`);
  }

  // Log slow responses (only for very slow queries)
  if (response.timings.duration > 30000) {
    console.log(`🐌 SLOW response on ${url} - Duration: ${response.timings.duration}ms`);
  }

  // Success rate calculation
  successRate.add(response.status === 200 ? 1 : 0);

  // Small sleep to simulate real user behavior (think time)
  sleep(1);
}

export function handleSummary(data) {
  console.log('\n========================================');
  console.log('📊 LOAD TEST SUMMARY');
  console.log('========================================');
  console.log(`Total Requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Failed Requests: ${data.metrics.http_req_failed.values.passes || 0}`);
  console.log(`Request Rate: ${data.metrics.http_reqs.values.rate.toFixed(2)} req/s`);
  console.log(`\nResponse Times:`);
  console.log(`  Min: ${data.metrics.http_req_duration.values.min.toFixed(2)}ms`);
  console.log(`  Avg: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  console.log(`  Max: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms`);
  console.log(`  p(95): ${(data.metrics.http_req_duration.values['p(95)'] || 0).toFixed(2)}ms`);
  console.log(`  p(99): ${(data.metrics.http_req_duration.values['p(99)'] || 0).toFixed(2)}ms`);
  console.log(`\nErrors: ${data.metrics.errors?.values?.count || 0}`);
  console.log(`Timeouts: ${data.metrics.timeouts?.values?.count || 0}`);
  console.log('========================================\n');

  return {
    'stdout': JSON.stringify(data, null, 2),
  };
}
