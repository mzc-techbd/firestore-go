import http from 'k6/http';
import { sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
  stages: [
    { duration: '10m', target: 500 }, // 예상 AVG QPS에 도달할 때까지
    { duration: '1m', target: 0 }, // 1분 동안 점진적으로 종료
    { duration: '10m', target: 2000 }, // 예상 MAX QPS 유지하면서
    { duration: '1m', target: 0 }, // 1분 동안 점진적으로 종료
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'], // 에러율 임계값 설정
    http_req_duration: ['p(95)<500'], // 95% 응답 시간 임계값 설정
  },
};

export default function () {
  // const url = 'https://firestore-go-264172533638.asia-northeast3.run.app/read'; // Cloud run 엔드포인트
  const url = 'http://localhost:8080/read';

  const res = http.get(url);

  // 응답 확인
  if (res.status !== 200) {
    console.error(`Request failed with status: ${res.status}`);
  }

  // sleep(1); // 요청 간의 지연 시간 - stages 옵션을 사용하면 sleep 제거
}

export function handleSummary(data) {
	return {
		"summary.html": htmlReport(data),
	}
}
