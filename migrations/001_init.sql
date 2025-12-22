-- 001_init.sql

-- USERS (mínimo)
CREATE TABLE IF NOT EXISTS users (
                                     id            TEXT PRIMARY KEY,
                                     email         TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     active        BOOLEAN NOT NULL DEFAULT TRUE,
                                     created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_users_updated CHECK (updated_at >= created_at)
    );

-- DOCTORS
CREATE TABLE IF NOT EXISTS doctors (
                                       id              TEXT PRIMARY KEY,
                                       user_id         TEXT NOT NULL UNIQUE,
                                       name            TEXT NOT NULL,
                                       registry_type   TEXT NOT NULL,
                                       registry_number TEXT NOT NULL,
                                       specialty       TEXT NOT NULL,
                                       session_minutes INT NOT NULL,
                                       price_cents     BIGINT NOT NULL,
                                       active          BOOLEAN NOT NULL DEFAULT TRUE,
                                       created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_doctor_registry_type CHECK (registry_type IN ('CRM','CRP','OTHER')),
    CONSTRAINT chk_doctor_session CHECK (session_minutes > 0),
    CONSTRAINT chk_doctor_price CHECK (price_cents >= 0),
    CONSTRAINT chk_doctor_updated CHECK (updated_at >= created_at),

    CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

-- PATIENTS
CREATE TABLE IF NOT EXISTS patients (
                                        id          TEXT PRIMARY KEY,
                                        user_id     TEXT NOT NULL UNIQUE,
                                        name        TEXT NOT NULL,
                                        active      BOOLEAN NOT NULL DEFAULT TRUE,
                                        created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_patient_updated CHECK (updated_at >= created_at),
    CONSTRAINT fk_patient_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

-- AVAILABILITY SLOTS
CREATE TABLE IF NOT EXISTS availability_slots (
                                                  id          TEXT PRIMARY KEY,
                                                  doctor_id   TEXT NOT NULL,
                                                  start_at    TIMESTAMPTZ NOT NULL,
                                                  end_at      TIMESTAMPTZ NOT NULL,
                                                  status      TEXT NOT NULL,
                                                  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_slot_time CHECK (start_at < end_at),
    CONSTRAINT chk_slot_status CHECK (status IN ('AVAILABLE','BOOKED','BLOCKED')),
    CONSTRAINT chk_slot_updated CHECK (updated_at >= created_at),

    CONSTRAINT fk_slot_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(id)
    );

CREATE INDEX IF NOT EXISTS idx_slots_doctor_time
    ON availability_slots (doctor_id, start_at);

-- APPOINTMENTS
CREATE TABLE IF NOT EXISTS appointments (
                                            id          TEXT PRIMARY KEY,
                                            doctor_id   TEXT NOT NULL,
                                            patient_id  TEXT NOT NULL,
                                            slot_id     TEXT NOT NULL,
                                            price_cents BIGINT NOT NULL,
                                            status      TEXT NOT NULL,
                                            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_appt_price CHECK (price_cents >= 0),
    CONSTRAINT chk_appt_status CHECK (status IN ('SCHEDULED','PAID','CANCELED','COMPLETED')),
    CONSTRAINT chk_appt_updated CHECK (updated_at >= created_at),

    CONSTRAINT uq_appointments_slot UNIQUE (slot_id),

    CONSTRAINT fk_appt_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    CONSTRAINT fk_appt_patient FOREIGN KEY (patient_id) REFERENCES patients(id),
    CONSTRAINT fk_appt_slot FOREIGN KEY (slot_id) REFERENCES availability_slots(id)
    );

CREATE INDEX IF NOT EXISTS idx_appt_doctor_created
    ON appointments (doctor_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_appt_patient_created
    ON appointments (patient_id, created_at DESC);

-- PAYMENTS
CREATE TABLE IF NOT EXISTS payments (
                                        id             TEXT PRIMARY KEY,
                                        appointment_id TEXT NOT NULL UNIQUE,
                                        provider       TEXT NOT NULL,
                                        amount_cents   BIGINT NOT NULL,
                                        status         TEXT NOT NULL,
                                        external_ref   TEXT,
                                        created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_pay_amount CHECK (amount_cents >= 0),
    CONSTRAINT chk_pay_provider CHECK (provider IN ('STRIPE','MERCADOPAGO','MANUAL')),
    CONSTRAINT chk_pay_status CHECK (status IN ('PENDING','APPROVED','REJECTED','REFUNDED','CANCELED')),
    CONSTRAINT chk_pay_updated CHECK (updated_at >= created_at),

    CONSTRAINT fk_pay_appt FOREIGN KEY (appointment_id) REFERENCES appointments(id)
    );

CREATE INDEX IF NOT EXISTS idx_payment_appt
    ON payments (appointment_id);
