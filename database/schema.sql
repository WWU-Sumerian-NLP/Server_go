CREATE TABLE [entities](
    id INTEGER,
    entity_name TEXT,
    entity_type TEXT
);

CREATE TABLE [relations](
    -- id INTEGER,
    relation_type TEXT,
    subject_tag TEXT,
    object_tag TEXT,
    regex_rules TEXT,
    tags TEXT
);

CREATE TABLE [relationships](
    id INTEGER,
    tablet_num TEXT,
    relation_type TEXT,
    subj TEXT,
    obj TEXT,
    providence TEXT,
    time_period TEXT,
    dates_referenced TEXT
);