-- Create table: authors
CREATE TABLE IF NOT EXISTS authors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT,
    bio TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create table: books
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    isbn TEXT UNIQUE,
    author_id INTEGER NOT NULL,
    category_id INTEGER,
    description TEXT,
    pages INTEGER DEFAULT 0,
    available INTEGER DEFAULT 1,
    price REAL DEFAULT 0.000000,
    published_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create table: borrowers
CREATE TABLE IF NOT EXISTS borrowers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    address TEXT,
    is_active INTEGER DEFAULT 1,
    max_books INTEGER DEFAULT 5,
    joined_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create table: categories
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create table: loans
CREATE TABLE IF NOT EXISTS loans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id INTEGER NOT NULL,
    borrower_id INTEGER NOT NULL,
    borrowed_at TIMESTAMP,
    due_at TIMESTAMP NOT NULL,
    returned_at TIMESTAMP,
    is_returned INTEGER DEFAULT 0,
    fine_amount REAL DEFAULT 0.000000,
    notes TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

