-- name: GetUser :one
SELECT id FROM USERS
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO USERS (
  "name", email, country, province, age, gender, proficiency_level, studied_quran_before, job_title, lisper
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id;

-- name: GetUsersRanks :many
SELECT 
    u.id, 
    u.name, 
    u.email, 
    COALESCE(SUM(ut.audio_duration), 0)::float AS total_duration,
    RANK() OVER (ORDER BY COALESCE(SUM(ut.audio_duration), 0) DESC)::int AS rank
FROM USERS u
LEFT JOIN USER_TASKS ut ON u.id = ut.user_id
GROUP BY u.id, u.name, u.email, u.created_at
ORDER BY total_duration DESC, u.created_at ASC
LIMIT 10;

-- name: GetUserRank :one
SELECT 
    u.id, 
    u.name, 
    u.email, 
    COALESCE(SUM(ut.audio_duration), 0)::float AS total_duration,
    (
        SELECT COUNT(*) + 1
        FROM USERS u2
        WHERE (
            SELECT COALESCE(SUM(ut2.audio_duration), 0)
            FROM USER_TASKS ut2
            WHERE ut2.user_id = u2.id
        ) > (
            SELECT COALESCE(SUM(ut3.audio_duration), 0)
            FROM USER_TASKS ut3
            WHERE ut3.user_id = $1
        )
    )::int AS rank
FROM USERS u
LEFT JOIN USER_TASKS ut ON u.id = ut.user_id
WHERE u.id = $1
GROUP BY u.id, u.name, u.email;

-- name: GetTask :one
SELECT t.id, ay.surah, ay.ayah, ay.text, ay.glyphs, ay.page
FROM TASKS t
LEFT JOIN USER_TASKS ut ON t.id = ut.task_id
JOIN AYAHS ay ON t.ayah = ay.id
WHERE t.is_deleted = FALSE
GROUP BY t.id, ay.id
ORDER BY COUNT(ut.id) ASC, RANDOM()
LIMIT 1;

-- name: GetTaskByID :one
SELECT t.id, ay.surah, ay.ayah, ay.text, ay.glyphs, ay.page
FROM TASKS t
JOIN AYAHS ay ON t.ayah = ay.id
WHERE t.id = $1;

-- name: GetUsersTasks :many
SELECT ut.audio_url, ut.audio_duration, ay.surah, ay.ayah, ay.text
FROM USER_TASKS ut
JOIN TASKS t ON ut.task_id = t.id
JOIN AYAHS ay ON t.ayah = ay.id
GROUP BY ut.id, t.ayah, ay.id
ORDER BY ut.created_at DESC;

-- name: CreateUserTask :exec
INSERT INTO USER_TASKS (
  user_id, task_id, audio_url, audio_duration
) VALUES (
  $1, $2, $3, $4
);

-- name: GetUserTotalDuration :one
SELECT COALESCE(SUM(audio_duration), 0) AS total_duration FROM USER_TASKS 
WHERE user_id = $1;

-- name: GetTotalDurations :one
SELECT COALESCE(SUM(audio_duration), 0) AS total_duration FROM USER_TASKS;

-- name: GetAyahs :many
SELECT * FROM AYAHS;

-- name: GetAdminUser :one
SELECT id, username, password_hash FROM ADMINS
WHERE username = $1 LIMIT 1;

-- name: InsertAdminUser :one
INSERT INTO ADMINS (username, password_hash) VALUES ($1, $2) RETURNING id;