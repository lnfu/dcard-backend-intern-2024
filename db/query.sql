--
-- name: GetActiveAdvertisements :many
SELECT advertisement.id,
    title,
    start_at,
    end_at
FROM advertisement
WHERE advertisement.start_at < NOW()
    AND advertisement.end_at > NOW()
    AND NOT EXISTS (
        SELECT 1
        FROM advertisement_cond
        WHERE advertisement_cond.advertisement_id = advertisement.id
    )
UNION ALL
SELECT advertisement.id,
    title,
    start_at,
    end_at
FROM advertisement
WHERE advertisement.start_at < NOW()
    AND advertisement.end_at > NOW()
    AND EXISTS(
        SELECT 1
        FROM advertisement_cond
        WHERE (
                advertisement_cond.advertisement_id = advertisement.id
                AND advertisement_cond.cond_id IN (
                    SELECT DISTINCT cond.id
                    FROM cond
                        LEFT JOIN cond_gender ON cond.id = cond_gender.cond_id
                        LEFT JOIN cond_country ON cond.id = cond_country.cond_id
                        LEFT JOIN cond_platform ON cond.id = cond_platform.cond_id
                        LEFT JOIN gender ON cond_gender.gender_id = gender.id
                        LEFT JOIN country ON cond_country.country_id = country.id
                        LEFT JOIN platform ON cond_platform.platform_id = platform.id
                    WHERE (1 = 1)
                        AND (
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
                        AND(
                            sqlc.narg(gender) IS NULL
                            OR (
                                gender.code IS NULL
                                OR sqlc.narg(gender) = gender.code
                            )
                        )
                        AND(
                            sqlc.narg(country) IS NULL
                            OR (
                                country.code IS NULL
                                OR sqlc.narg(country) = country.code
                            )
                        )
                        AND(
                            sqlc.narg(platform) IS NULL
                            OR (
                                platform.name IS NULL
                                OR sqlc.narg(platform) = platform.name
                            )
                        )
                )
            )
    )
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
-- name: CheckGender :one
SELECT COUNT(*)
FROM gender
WHERE code = sqlc.arg(gender);
--
-- name: CheckCountry :one
SELECT COUNT(*)
FROM country
WHERE code = sqlc.arg(country);
--
-- name: CheckPlatform :one
SELECT COUNT(*)
FROM platform
WHERE name = sqlc.arg(platform);