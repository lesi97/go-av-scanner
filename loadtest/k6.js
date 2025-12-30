import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

export const options = {
    scenarios: {
        scan: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '10s', target: 50 },
                { duration: '10s', target: 200 },
                { duration: '10s', target: 500 },
                { duration: '10s', target: 1000 },
                { duration: '10s', target: 0 },
            ],
            gracefulRampDown: '30s',
        },
    },
};

const url = 'http://host.docker.internal:8080/api/scan';
const bin = open('./loadtest/loadtest.bin', 'b');
const statusCount = new Counter('http_status');

export default function () {
    const res = http.post(
        url,
        { file: http.file(bin, 'loadtest.bin', 'application/octet-stream') },
        { timeout: '120s' }
    );

    statusCount.add(1, { status: String(res.status) });

    check(res, {
        'status is 200 or 422': (r) => r.status === 200 || r.status === 422,
    });

    sleep(0.1);
}
