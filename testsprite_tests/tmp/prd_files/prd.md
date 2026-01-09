# Product Requirements Document (PRD) - Club Pulse System

**Target Platform**: TestSprite (AI Automated Testing)
**Application Type**: Full Stack Web Application (Go Backend + Next.js Frontend)
**Architecture**: Monolythic Modular (Clean Architecture)

## 1. Project Overview & Scope

Club Pulse System is a comprehensive SaaS platform for sports club management, designed to handle multi-tenancy (`clubs` table). It orchestrates users, memberships, facility bookings, payments, championships, and team logistics.

### Tech Stack
-   **Frontend**: Next.js 14 (App Router), Tailwind CSS, Shadcn UI (Radix Primitives), React Hook Form, Zod.
-   **Backend**: Go 1.24, Gin Framework, GORM (PostgreSQL), Redis (Session/Caching).
-   **Infrastructure**: Docker Compose, PostgreSQL (pgvector enabled), Redis, SendGrid (Email), Twilio (SMS).
-   **Authentication**: JWT-based (stored in cookies), Google OAuth, efficient Role-Based Access Control (RBAC).

## 2. User Roles & Permissions

Reflected in `users.role` and validated via `authMiddleware` and `tenantMiddleware`.

1.  **SUPER_ADMIN** (`RoleSuperAdmin`): System-wide access. Can manage clubs (tenants) and global configurations.
2.  **ADMIN** (`RoleAdmin`): Club-level administrator. Manages facilities, memberships, tournaments, products, and financial settings for their specific club.
3.  **COACH** (`RoleCoach`): Manages specific `training_groups` and `teams`. Handles `attendance_records` and `travel_events`.
4.  **MEDICAL_STAFF** (`RoleMedicalStaff`): Special access to health data (GDPR Article 9 compliance). Validates `user_documents` (e.g., medical certificates).
5.  **MEMBER** (`RoleMember`): Standard user. Can book `facilities`, pay `subscriptions`, purchase `products`, and view `championships` stats.
6.  **GUEST**: Unauthenticated or limited access. Can view public landing pages, store (read-only), and tournament public pages.

## 3. Data Model (Core Entities)
*Based on `001_initial_schema.sql`*

*   **Clubs**: Tenants. Contains settings, theming, and domain config.
*   **Users**: Global members linked to a `club_id`. Includes GDPR fields (`terms_accepted_at`, `privacy_policy_version`).
*   **Facilities**: Bookable resources (courts, fields) with pricing and capacity.
*   **Bookings**: Reservations with status (`PENDING`, `CONFIRMED`, `PAID`) and Redis locking for concurrency.
*   **Championships**: Tournaments with formats (League/Knockout). Linked to `disciplines`.
*   **Teams**: Rosters participating in championships or training. linked to `travel_events`.
*   **Payments**: Transaction records linked to `orders` or `subscriptions`. via MercadoPago.
*   **User Documents**: Verification assets (IDs, Medical Certs) with expiration tracking.

## 4. Core Modules & Features

### 4.1 Authentication & Authorization (`/auth`)
*   **Endpoints**: `POST /auth/login`, `POST /auth/register`, `POST /auth/refresh`, `POST /auth/google`.
*   **Logic**:
    *   **Login**: Validates credentials -> Generates JWT -> Stores session in Redis -> Sets HttpOnly Cookie.
    *   **RBAC**: Middleware checks `Claims.Role` against required permission.
    *   **GDPR**: Registration requires ticking consent boxes (Terms, Privacy). Database records consent version in `consent_records`.

### 4.2 Facility Booking (`/booking`)
*   **Endpoints**: `GET /facilities`, `POST /bookings`, `GET /bookings/availability`.
*   **Logic**:
    *   **Conflict Detection**: Uses `idx_bookings_availability` and Redis locks (`booking_lock.go`) to prevent double-booking.
    *   **Pricing**: dynamic calculation based on `hourly_rate` or custom slot pricing.
    *   **Validation**: Checks if user has valid `Membership` (if required) or pending debts.

### 4.3 Membership & Access (`/membership`, `/access`)
*   **Endpoints**: `GET /subscriptions`, `POST /subscriptions/assign`, `GET /access/check`.
*   **Logic**:
    *   **Tiers**: `membership_tiers` define benefits and pricing.
    *   **Scholarships**: `scholarships` table allows percentage discounts on fees.
    *   **Access Control**: Physical access checks via QR/Pin (simulated) verifying `subscription.status == 'ACTIVE'`.

### 4.4 Championship Management (`/championship`)
*   **Endpoints**: `POST /tournaments`, `POST /matches`, `GET /standings`.
*   **Logic**:
    *   **Fixtures**: Automated generation of `matches` based on registered `teams`.
    *   **Standings**: Calculated on-the-fly or cached, sorting by Points -> Goal Diff -> Goals For.
    *   **Volunteers**: `volunteer_assignments` for managing parents/staff during matches.

### 4.5 Team Logistics (`/team`)
*   **Features**:
    *   **Attendance**: Coaches mark `attendance_records` for `training_groups`.
    *   **Travel**: `travel_events` manage away games, transport details, and `event_rsvps` from players.
    *   **Player Status**: "Traffic Light" system (Green/Red) based on `medical_cert_status` and `membership` status.

### 4.6 Store & Payments (`/store`, `/payment`)
*   **Endpoints**: `GET /products`, `POST /checkout`.
*   **Logic**:
    *   **Cart**: Client-side managed, validated on checkout.
    *   **Checkout**: Integrates MercadoPago. Creates `orders` with `PENDING` status. Webhook updates to `PAID`.

## 5. Detailed User Flows (E2E Candidates)

### Flow A: Member Booking with Payment
1.  **Login**: User logs in as Member.
2.  **Discovery**: Navigate to `/bookings`. View Calendar.
3.  **Selection**: Click 18:00 slot (Red) -> Backend checks Redis lock.
4.  **Confirmation**: Modal shows price. Click "Pay & Book".
5.  **Payment**: Redirection to Fake/Sandbox Gateway -> Success.
6.  **Verification**: Redirect to `/dashboard`. Toast "Reserva Confirmada". Booking appears in list.

### Flow B: Season Management (Admin)
1.  **Setup**: Admin creates a new `Championship` ("Winter League 2026").
2.  **Registration**: Adds 6 `Teams`.
3.  **Fixture**: Clicks "Generate Fixture". System creates `Matches` (Round Robin).
4.  **Result**: Admin updates Match 1 (Team A 3 - 1 Team B).
5.  **Verify**: Standings table updates immediately (Team A: 3pts).

### Flow C: GDPR & Document Compliance
1.  **Upload**: User goes to `/profile`. Uploads "Medical Certificate".
2.  **Review**: Login as `MEDICAL_STAFF`. Go to `/admin/documents`.
3.  **Action**: View document. Click "Approve". Set expiry date.
4.  **Effect**: User's `MedicalCertStatus` becomes `VALID`. Player Status turns Green.

## 6. API Testing Context (Backend)
*   **Base URL**: `http://localhost:8080/api/v1`
*   **Auth Header**: `Authorization: Bearer <token>` (though standard browser flow uses Cookies, API tests might use Header).
*   **Headers**: `X-Tenant-ID` (optional, usually inferred from Domain/User).
*   **Test Data**: Use `cmd/seeder` or `migrations/001_initial_schema.sql` to understand constraints.

## 7. Frontend Testing Context (Selectors)
*   **Tech**: Shadcn UI uses Radix Primitives. Look for `role="dialog"`, `role="checkbox"`.
*   **Ids**: Key elements have `data-testid` (e.g., `booking-slot-{id}`, `btn-login`).
*   **Navigation**: Uses `next/navigation`. `router.push()` triggers client-side transitions.
