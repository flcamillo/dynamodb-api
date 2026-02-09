// load-test.js
// k6 Load Testing Script for DynamoDB API
// Usage: k6 run load-test.js

import http from 'k6/http';
import { check, group, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:7000';

export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Ramp-up
    { duration: '3m', target: 50 },   // Stay at 50
    { duration: '2m', target: 100 },  // Ramp-up to 100
    { duration: '3m', target: 100 },  // Stay at 100
    { duration: '2m', target: 0 },    // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.1'],
  },
};

export function setup() {
  // Create test data
  const response = http.post(
    `${BASE_URL}/eventos`,
    JSON.stringify({
      date: new Date().toISOString(),
      statusCode: 200,
      statusMessage: 'Setup test event',
      metadata: { test: 'setup' },
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );

  const eventId = response.json('id');
  console.log(`Created test event: ${eventId}`);
  return { eventId };
}

export default function (data) {
  group('Health Check', () => {
    const res = http.get(`${BASE_URL}/health`);
    check(res, {
      'status is 200': (r) => r.status === 200,
      'response is OK': (r) => r.body === 'OK',
    });
  });

  group('Create Event', () => {
    const payload = JSON.stringify({
      date: new Date().toISOString(),
      statusCode: 200,
      statusMessage: 'Load test event',
      metadata: {
        test_run: __ENV.TEST_RUN || 'default',
        timestamp: new Date().toISOString(),
      },
    });

    const res = http.post(`${BASE_URL}/eventos`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
      'status is 201': (r) => r.status === 201,
      'has event id': (r) => r.json('id') !== undefined,
      'has expiration': (r) => r.json('expiration') !== undefined,
    });

    if (res.status === 201) {
      const eventId = res.json('id');

      group('Get Event', () => {
        const getRes = http.get(`${BASE_URL}/eventos/${eventId}`);
        check(getRes, {
          'status is 200': (r) => r.status === 200,
          'event id matches': (r) => r.json('id') === eventId,
        });
      });

      group('Update Event', () => {
        const updatePayload = JSON.stringify({
          date: new Date().toISOString(),
          statusCode: 202,
          statusMessage: 'Updated by load test',
          metadata: { updated: true },
        });

        const updateRes = http.put(`${BASE_URL}/eventos/${eventId}`, updatePayload, {
          headers: { 'Content-Type': 'application/json' },
        });

        check(updateRes, {
          'status is 200': (r) => r.status === 200,
          'status code updated': (r) => r.json('statusCode') === 202,
        });
      });
    }
  });

  group('Find Events', () => {
    const now = new Date();
    const startDate = new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString();
    const endDate = new Date(now.getTime() + 60 * 1000).toISOString();

    const res = http.get(
      `${BASE_URL}/eventos?startDate=${startDate}&endDate=${endDate}&statusCode=200`
    );

    check(res, {
      'status is 200': (r) => r.status === 200,
      'has items': (r) => r.json('items') !== undefined,
      'has total': (r) => r.json('total') !== undefined,
    });
  });

  sleep(1);
}

export function teardown(data) {
  // Cleanup test data
  if (data.eventId) {
    const res = http.del(`${BASE_URL}/eventos/${data.eventId}`);
    check(res, {
      'deleted successfully': (r) => r.status === 200,
    });
    console.log(`Cleaned up test event: ${data.eventId}`);
  }
}
