---
description: Fix Facility OpeningHour/ClosingHour field names
---

1. Replace `OpeningHour` with `OpeningTime` in `backend/internal/modules/facilities/application/usecases.go`.
2. Replace `ClosingHour` with `ClosingTime` in `backend/internal/modules/facilities/application/usecases.go`.
3. Replace `OpeningHour` with `OpeningTime` in `backend/internal/modules/facilities/infrastructure/repository/postgres.go`.
4. Replace `ClosingHour` with `ClosingTime` in `backend/internal/modules/facilities/infrastructure/repository/postgres.go`.
5. Replace `OpeningHour` with `OpeningTime` in `backend/internal/modules/booking/application/usecases.go`.
6. Replace `ClosingHour` with `ClosingTime` in `backend/internal/modules/booking/application/usecases.go`.
7. Replace `OpeningHour` with `OpeningTime` in `backend/internal/modules/booking/infrastructure/http/handler_test.go`.
8. Replace `ClosingHour` with `ClosingTime` in `backend/internal/modules/booking/infrastructure/http/handler_test.go`.
9. Replace `OpeningHour` with `OpeningTime` in `backend/internal/modules/booking/application/usecases_test.go`.
10. Replace `ClosingHour` with `ClosingTime` in `backend/internal/modules/booking/application/usecases_test.go`.
11. Replace `OpeningHour` with `OpeningTime` in `backend/tests/e2e/booking_pricing_test.go`.
12. Replace `ClosingHour` with `ClosingTime` in `backend/tests/e2e/booking_pricing_test.go`.

// turbo-all
