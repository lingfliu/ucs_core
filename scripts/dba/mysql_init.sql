create table if not exists d_node_meta (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    parent_id INTEGER,
    class INTEGER NOT NULL,
    model VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    descrip VARCHAR(255),
);

create table if not exists d_point_meta (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    dnode_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    class INTEGER NOT NULL,
    unit VARCHAR(255),
    dimen INTEGER,
    byte_len INTEGER,
    sample_len INTEGER,
    coder VARCHAR(255),
);

create table if not exists ctl_node_meta (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    parent_id INTEGER,
    class INTEGER NOT NULL,
    model VARCHAR(255) NOT NULL,
);

create table if not exists c_point_meta (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    cnode_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    class INTEGER NOT NULL,
    unit VARCHAR(255),
    dimen INTEGER,
    byte_len INTEGER,
    sample_len INTEGER,
    coder VARCHAR(255),
);

