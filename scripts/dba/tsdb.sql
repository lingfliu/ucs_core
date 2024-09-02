create stable if not exists dp_single values (
    ts timestamp;
    value int;
) tags (dnode_id int, dp_id int, dnode_class int, dp_class int);
