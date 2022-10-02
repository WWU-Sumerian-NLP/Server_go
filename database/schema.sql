CREATE TABLE [entities](
    id INTEGER NOT NULL PRIMARY KEY,
    entity_name TEXT,
    entity_tag TEXT,
    unique (entity_name, entity_tag)
);

CREATE TABLE [relations](
    id INTEGER NOT NULL PRIMARY KEY,
    relation_type TEXT,
    subject_tag TEXT,
    object_tag TEXT,
    regex_rules TEXT,
    tags TEXT,
    unique (relation_type, subject_tag, object_tag, regex_rules, tags)
);

CREATE TABLE [relationships](
    id INTEGER NOT NULL PRIMARY KEY,
    tablet_num TEXT,
    relation_type TEXT,
    subj TEXT,
    obj TEXT,
    providence TEXT,
    time_period TEXT,
    dates_referenced TEXT,
    unique (tablet_num, relation_type)
);