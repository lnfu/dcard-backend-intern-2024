CREATE INDEX idx_advertisement_cond_advertisement_id ON advertisement_cond (advertisement_id);
CREATE INDEX idx_advertisement_cond_cond_id ON advertisement_cond (cond_id);

CREATE INDEX idx_cond_gender_cond_id ON cond_gender (cond_id);
CREATE INDEX idx_cond_gender_gender_id ON cond_gender (gender_id);

CREATE INDEX idx_cond_country_cond_id ON cond_country (cond_id);
CREATE INDEX idx_cond_country_country_id ON cond_country (country_id);

CREATE INDEX idx_cond_platform_cond_id ON cond_platform (cond_id);
CREATE INDEX idx_cond_platform_platform_id ON cond_platform (platform_id);