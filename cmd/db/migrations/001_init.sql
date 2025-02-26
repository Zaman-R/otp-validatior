CREATE TABLE otps (
                      id TEXT PRIMARY KEY,
                      purpose TEXT NOT NULL,
                      delivery_method TEXT NOT NULL,
                      mobile TEXT,
                      email TEXT,
                      hashed_otp TEXT NOT NULL,
                      status TEXT NOT NULL DEFAULT 'pending',
                      created_at TIMESTAMP DEFAULT now(),
                      expires_at TIMESTAMP NOT NULL,
                      retry_count INT DEFAULT 0,
                      retry_limit INT NOT NULL
);
