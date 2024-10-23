import http from 'k6/http';
import { sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
  vus: 20, // 가상 사용자 수를 점진적으로 늘려가며 테스트
  duration: '5m', // 테스트 기간
  thresholds: {
    http_req_failed: ['rate<0.01'], // 에러율 임계값 설정
    http_req_duration: ['p(95)<500'], // 95% 응답 시간 임계값 설정
  },
};

export default function () {
  // const url = 'https://firestore-go-264172533638.asia-northeast3.run.app/read'; // Firestore 엔드포인트
  const url = 'http://localhost:8080/read'; // Firestore 엔드포인트

  const res = http.get(url);

  // 응답 확인
  if (res.status !== 200) {
    console.error(`Request failed with status: ${res.status}`);
  }

  sleep(1); // 요청 간의 지연 시간
}

export function handleSummary(data) {
	return {
		"summary.html": htmlReport(data),
	}
}
