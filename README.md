# Commerce — Microservices (Clean Architecture) — Go + gRPC + PostgreSQL 17

Implementasi **Clean Architecture** untuk 2 service:
- **order-service** (REST, Gin) — menerima order dan memanggil inventory via gRPC.
- **inventory-service** (gRPC) — mengelola stok produk (PostgreSQL 17), menyediakan `CheckStock`, `ReserveStock`, `ReleaseStock`.

Fitur:
- Struktur **Clean Architecture**: `domain` (entity + port), `usecase` (bisnis), `infrastructure` (DB, transport), `cmd/server` (wiring).
- Konfigurasi via **`.env`** setiap service.
- **gRPC reflection** aktif → bisa debug dengan `grpcurl`.
- **Healthcheck** Postgres + **retry** koneksi di service saat startup.
- **Docker Compose** (monorepo build), Go **1.24**, Postgres **17**.

## Cara Menjalankan
```bash
# 1) Generate protobuf
make proto-all

# 2) Build semua service
make build-all

# 3) Jalankan semua container
docker compose up --build
```

### Seed & Uji Coba
Seed produk:
```bash
curl -X POST http://localhost:9090/seed -H "Content-Type: application/json" -d '{
  "products": [
    {"sku":"SKU-001","name":"Produk A","stock":10},
    {"sku":"SKU-002","name":"Produk B","stock":5}
  ]
}'
```

Buat order:
```bash
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d '{
  "items": [{"sku":"SKU-001","qty":2},{"sku":"SKU-002","qty":1}]
}'
```

Release stok (via grpcurl):
```bash
grpcurl -plaintext -d '{"items":[{"sku":"SKU-001","qty":2}]}' localhost:50051 inventory.v1.InventoryService/ReleaseStock
```

Check stok (via grpcurl):
```bash
grpcurl -plaintext -d '{"items":[{"sku":"SKU-001","qty":10}]}' localhost:50051 inventory.v1.InventoryService/ReleaseStock
```

## Struktur Direktori
```
proto/                          # kontrak protobuf (+ Makefile generator)
inventory-service/
  cmd/server/                   # wiring + bootstrap
  internal/
    config/                     # loader .env
    domain/inventory/           # entity + repository port
    usecase/                    # business rules
  infrastructure/
    persistence/postgres/       # implementasi repository GORM
    transport/grpcsvr/          # gRPC server adapter
    transport/httpadmin/        # HTTP admin (seed)
order-service/
  cmd/server/
  internal/
    config/
    domain/order/               # entity + repo port
    domain/inventory/           # port untuk InventoryClient
    usecase/                    # usecase PlaceOrder
  infrastructure/
    persistence/postgres/
    transport/httpapi/          # REST handler
    transport/grpcclient/       # gRPC client ke inventory
design/system-design.png        # diagram PNG
docker-compose.yml
Makefile
go.work
```

## Diagram Sistem (PNG)
Lihat file: `design/system-design.png`
