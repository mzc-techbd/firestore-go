# firestore-go
Google Cloud Firestore에서 데이터를 읽는 간단한 Go 테스트앱으로, Prometheus를 사용하여 성능 모니터링

# Feature
- 지정된 Firestore collection 및 document에서 단순 읽기 작업 수행
- 읽기 작업의 총 개수(`firestore_reads_total`) 및 지연 시간(`firestore_read_latency_seconds`)을 추적하는 Prometheus 지표를 생성
- `/read` 엔드포인트에서 Firestore 샘플데이터를 반환
- `/metrics` 엔드포인트에서 Prometheus 사용

## When use Cloud Run
```bash
gcloud auth application-default print-access-token > token.txt
AUTH_TOKEN=$(cat token.txt)
```

```bash
for i in $(seq 1 100); do
  curl -X POST \
    -H "Content-Type: application/json" \
    -d '' \
    https://firestore-go-<SOME_CLOUD_RUN_URL>.asia-northeast3.run.app/read
done
```

## Local test
```bash
## check gcloud info before start
# gcloud init

# get require packages form go.mod
go mod tidy

# start from local
go run .
```

## Requirements before start using container
1. Create GCP Service account.
2. Assign role created service account.
3. Download service account key file.
4. move to `./docker/service-account.json`

## Do not shared service account key file.

## Docker
```bash
# init
docker compose up -d --build

# clean up
docker compose down
```

## When use k6
### Save output to file format
```bash
# Save JSON format
k6 run --out json=results.json loadtest.js 
# Save CSV format
k6 run --out csv=results.csv loadtest.js 
```

### Export metrics to Prometheus
```bash
K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write \
K6_PROMETHEUS_RW_TREND_AS_NATIVE_HISTOGRAM=true \
k6 run -o experimental-prometheus-rw loadtest-detail.js
```