--
-- name: GetActiveAdvertisements :many
SELECT DISTINCT adv.id,
    adv.title,
    adv.start_at,
    adv.end_at
FROM advertisement adv
    LEFT JOIN advertisement_cond adc ON adv.id = adc.advertisement_id
    LEFT JOIN cond ON adc.cond_id = cond.id
    LEFT JOIN cond_gender ON cond.id = cond_gender.cond_id
    LEFT JOIN gender ON cond_gender.gender_id = gender.id
    LEFT JOIN cond_country ON cond.id = cond_country.cond_id
    LEFT JOIN country ON cond_country.country_id = country.id
    LEFT JOIN cond_platform ON cond.id = cond_platform.cond_id
    LEFT JOIN platform ON cond_platform.platform_id = platform.id
WHERE (
        sqlc.narg(age) IS NULL
        OR (
            (
                cond.age_start IS NULL
                OR cond.age_start <= sqlc.narg(age)
            )
            AND (
                cond.age_end IS NULL
                OR cond.age_end >= sqlc.narg(age)
            )
        )
    )
    AND (
        sqlc.narg(gender) IS NULL
        OR gender.code = sqlc.narg(gender)
        OR cond_gender.cond_id IS NULL
    )
    AND (
        sqlc.narg(country) IS NULL
        OR country.code = sqlc.narg(country)
        OR cond_country.cond_id IS NULL
    )
    AND (
        sqlc.narg(platform) IS NULL
        OR platform.name = sqlc.narg(platform)
        OR cond_platform.cond_id IS NULL
    )
    OR adc.id IS NULL
ORDER BY end_at ASC
LIMIT ?, ?;
--
-- name: CreateAdvertisement :execlastid
INSERT INTO advertisement (title, start_at, end_at)
VALUES (
        sqlc.arg(title),
        sqlc.arg(start_at),
        sqlc.arg(end_at)
    );
--
-- name: CreateCondition :execlastid
INSERT INTO cond (age_start, age_end)
VALUES (
        sqlc.arg(age_start),
        sqlc.arg(age_end)
    );
-- 
-- name: CreateAdvertisementCondition :exec
INSERT INTO advertisement_cond (advertisement_id, cond_id)
VALUES (
        sqlc.arg(advertisement_id),
        sqlc.arg(condition_id)
    );
--
-- name: CreateConditionGender :exec
INSERT INTO cond_gender (cond_id, gender_id)
VALUES (
        sqlc.arg(condition_id),
        (
            SELECT id
            FROM gender
            WHERE code = sqlc.arg(gender)
        )
    );
--
-- name: CreateConditionCountry :exec
INSERT INTO cond_country (cond_id, country_id)
VALUES (
        sqlc.arg(condition_id),
        (
            SELECT id
            FROM country
            WHERE code = sqlc.arg(country)
        )
    );
--
-- name: CreateConditionPlatform :exec
INSERT INTO cond_platform (cond_id, platform_id)
VALUES (
        sqlc.arg(condition_id),
        (
            SELECT id
            FROM platform
            WHERE name = sqlc.arg(platform)
        )
    );
--
-- name: GetAllGenders :many
SELECT code
FROM gender;
--
-- name: GetAllCountries :many
SELECT code
FROM country;
--
-- name: GetAllPlatforms :many
SELECT name
FROM platform;