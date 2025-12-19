-- 001_init.sql

-- Patients
CREATE TABLE IF NOT EXISTS patients (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL,
  name        TEXT NOT NULL,
  active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL
);

-- Availability Slots
CREATE TABLE IF NOT EXISTS availability_slots (
  id          TEXT PRIMARY KEY,
  doctor_id   TEXT NOT NULL,
  start_at    TIMESTAMPTZ NOT NULL,
  end_at      TIMESTAMPTZ NOT NULL,
  status      TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,
  CONSTRAINT chk_slot_time CHECK (start_at < end_at),
  CONSTRAINT chk_slot_status CHECK (status IN ('AVAILABLE','BOOKED','BLOCKED'))
);

CREATE INDEX IF NOT EXISTS idx_slots_doctor_time
ON availability_slots (doctor_id, start_at);

-- Appointments
CREATE TABLE IF NOT EXISTS appointments (
  id          TEXT PRIMARY KEY,
  doctor_id   TEXT NOT NULL,
  patient_id  TEXT NOT NULL,
  slot_id     TEXT NOT NULL,
  price_cents BIGINT NOT NULL,
  status      TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,
  CONSTRAINT chk_appt_price CHECK (price_cents >= 0),
  CONSTRAINT chk_appt_status CHECK (status IN ('SCHEDULED','PAID','CANCELED','COMPLETED')),
  CONSTRAINT uq_appointments_slot UNIQUE (slot_id)
);

CREATE INDEX IF NOT EXISTS idx_appt_doctor_created
ON appointments (doctor_id, created_at DESC);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
  id             TEXT PRIMARY KEY,
  appointment_id TEXT NOT NULL,
  provider       TEXT NOT NULL,
  amount_cents   BIGINT NOT NULL,
  status         TEXT NOT NULL,
  external_ref   TEXT,
  created_at     TIMESTAMPTZ NOT NULL,
  updated_at     TIMESTAMPTZ NOT NULL,
  CONSTRAINT chk_pay_amount CHECK (amount_cents >= 0),
  CONSTRAINT chk_pay_status CHECK (status IN ('PENDING','APPROVED','REJECTED','REFUNDED','CANCELED'))
);

CREATE INDEX IF NOT EXISTS idx_payment_appt
ON payments (appointment_id);
