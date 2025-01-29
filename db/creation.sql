PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` TEXT NOT NULL UNIQUE,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL,
    `image` TEXT NOT NULL,
    `role` TEXT NOT NULL,
    `token` TEXT UNIQUE,
    `token_exp` INTEGER,  
    `created_at` INTEGER NOT NULL 
);

CREATE TABLE IF NOT EXISTS `posts` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `title` TEXT NOT NULL,
    `body` TEXT NOT NULL,
    `image` TEXT,
    `created_at` INTEGER NOT NULL, 
    `modified_at` INTEGER NOT NULL, 
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `comments` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER NOT NULL,
    `body` TEXT NOT NULL,
    `created_at` INTEGER NOT NULL,
    `modified_at` INTEGER NOT NULL, 
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `categories` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS `posts_categories` (
    `post_id` INTEGER NOT NULL,
    `category_id` INTEGER NOT NULL,
    PRIMARY KEY (`category_id`, `post_id`),
    FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `likes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER,
    `comment_id` INTEGER,
    `like` BOOLEAN NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE
);

