-- Create the database
CREATE DATABASE np_db;

-- Use the database
USE np_db;

-- Users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(15) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    membership_tier ENUM('Basic', 'Premium', 'VIP') DEFAULT 'Basic',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD COLUMN verification_token VARCHAR(255),
ADD COLUMN is_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN token_expiry DATETIME;

INSERT INTO users (email, phone, password_hash, membership_tier, created_at, updated_at, verification_token, is_verified, token_expiry)
VALUES
('vipUser@example.com', '1234567890','$2a$10$YvTv2Z9f0DZM2zc7moYsKeEzZsj1USM2/lUKx6NCaLdjplDd88.bq', 'VIP', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'token123', TRUE, NULL);

-- Vehicles table
CREATE TABLE vehicles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    license_plate VARCHAR(15) NOT NULL UNIQUE,
    model VARCHAR(255) NOT NULL,
    charge_level DECIMAL(5, 2) NOT NULL,
    cleanliness ENUM('Clean', 'Moderate', 'Dirty') DEFAULT 'Clean',
    cost DECIMAL(5,2) not NULL, 
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO vehicles (license_plate, model, charge_level, cleanliness, cost, location)
VALUES
('ABC1234', 'Tesla Model 3', 85.50, 'Clean', 20.00, 'Downtown'),
('XYZ5678', 'Nissan Leaf', 60.25, 'Moderate', 10.00, 'Uptown');


-- Reservations table
CREATE TABLE reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    total_price DECIMAL(10, 2),
    status ENUM('Active', 'Completed', 'Cancelled') DEFAULT 'Active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE
);

-- Payments table
CREATE TABLE payments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    user_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method ENUM('Credit Card', 'Debit Card', 'PayPal') NOT NULL,
    status ENUM('Pending', 'Completed', 'Refunded') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

