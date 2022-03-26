CREATE TABLE IF NOT EXISTS users (id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, username VARCHAR(40) NOT NULL, email VARCHAR(255) NOT NULL, password VARBINARY(60) NOT NULL);

CREATE TABLE IF NOT EXISTS threads (id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, name VARCHAR(40) NULL, type TINYINT UNSIGNED NOT NULL, created_at TIMESTAMP NOT NULL);

CREATE TABLE IF NOT EXISTS messages (id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, thread_id INT UNSIGNED NOT NULL, user_id INT UNSIGNED NOT NULL, date TIMESTAMP NOT NULL, FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE, FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);

CREATE TABLE IF NOT EXISTS users_threads (user_id INT UNSIGNED NOT NULL, thread_id INT UNSIGNED NOT NULL, added_by INT UNSIGNED NOT NULL, date TIMESTAMP NOT NULL, seen BOOLEAN NOT NULL, PRIMARY KEY (user_id, thread_id), FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,  FOREIGN KEY (added_by) REFERENCES users(id) ON DELETE CASCADE);

@fetch chats
SELECT a.`user_id` AS `user_id`,
       b.`id`      AS `thread_id`,
       b.`name`    AS `thread_name`,
       b.`type`    AS `thread_type`,
       c.`id`      AS `message_id`,
       c.`text`    AS `message_text`,
       a.`seen`,
       c.`date` AS `message_date`,
       c.`sender_id`,
       d.`username`  AS `sender_username`,
       f.`participants`
FROM   `users_threads` a
       JOIN `threads` b
         ON a.`thread_id` = b.`id` /* join users_threads and threads so we can get chat room name, id, type */
       JOIN (SELECT messages.`id`,
                    `thread_id`,
                    `text`,
                    `date`,
                    `user_id` AS `sender_id`
             FROM   messages
                    JOIN (SELECT Max(id) AS `id`
                          FROM   messages
                          WHERE  thread_id IN (SELECT thread_id
                                               FROM   users_threads
                                               WHERE  user_id = ?)  /* gets all threads of the user */
                          AND id > ?   /* for pagination */                
                            GROUP  BY thread_id
                            ORDER BY id DESC
                            LIMIT 20) b
                      ON messages.id = b.id) c 
         ON b.`id` = c.`thread_id` /* gets the last message sender, text, date*/
       JOIN users d
         ON d.id = c.sender_id /*maps the last message sender id to his user name */
       JOIN (SELECT thread_id,
                    Group_concat(b.username, "") AS `participants`
             FROM   users_threads a
                    JOIN users b
                      ON a.user_id = b.id
             WHERE  thread_id IN (SELECT thread_id
                                  FROM   users_threads
                                  WHERE  user_id = ?)
                    AND user_id != ?
             GROUP  BY thread_id) f
         ON f.thread_id = a.thread_id
WHERE  a.`user_id` = ?
ORDER  BY c.`id` DESC; 

@check if two users already have a thread between them
SELECT ut1.thread_id
FROM   users_threads ut1
       JOIN users_threads ut2
         ON ut1.thread_id = ut2.thread_id
       JOIN threads t
         ON ut1.thread_id = t.id
WHERE  ut1.user_id = ?
       AND ut2.user_id = ?
       AND t.type = 0; 