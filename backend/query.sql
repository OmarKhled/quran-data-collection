-- name: GetUser :one
SELECT id FROM USERS
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO USERS (
  "name", email, country, province, age, gender, proficiency_level, studied_quran_before, job_title, lisper
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id;


-- name: GetTask :one
SELECT t.id, ay.surah, ay.ayah, ay.text, ay.glyphs, ay.page
FROM TASKS t
LEFT JOIN USER_TASKS ut ON t.id = ut.task_id
JOIN AYAHS ay ON t.ayah = ay.id
WHERE t.is_deleted = FALSE
GROUP BY t.id, ay.id
ORDER BY COUNT(ut.id) ASC, RANDOM()
LIMIT 1;

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